#!/bin/bash

# Qdrant 独立部署脚本（方案A：系统外独立进程）
# 用途：在远端服务器部署并启动 Qdrant，供后端通过 URL 直连
# 前置：在本项目根目录预置 qdrant-binaries/qdrant-<platform> 二进制（例如 qdrant-linux-amd64, qdrant-linux-arm64）
# 说明：
#  - 默认安装目录：/pixelpunk-qdrant
#  - 默认端口：HTTP 6333，gRPC 6334
#  - 数据目录：<INSTALL_DIR>/shared/data/qdrant（绝对路径）
#  - 配置文件：<INSTALL_DIR>/config/qdrant.yaml
#  - 日志文件：<INSTALL_DIR>/logs/qdrant.log
#  - PID 文件：<INSTALL_DIR>/qdrant.pid

set -euo pipefail

# 颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

# 默认参数
REMOTE_HOST=""
REMOTE_USER="root"
SSH_PORT="22"
INSTALL_DIR="/pixelpunk-qdrant"
HTTP_PORT="6333"
GRPC_PORT="6334"
BINARIES_DIR="qdrant-binaries"
QDRANT_VERSION="v1.15.1"
FORCE="no"
REQUIRE_EMPTY_DIR="no"

print_info(){ echo -e "${CYAN}[INFO]${NC} $1"; }
print_ok(){ echo -e "${GREEN}[OK]${NC} $1"; }
print_warn(){ echo -e "${YELLOW}[WARN]${NC} $1"; }
print_err(){ echo -e "${RED}[ERR]${NC} $1"; }
usage(){
  cat <<EOF
用法: $0 --remote-host <ip> [--remote-user root] [--ssh-port 22] \
         [--install-dir /pixelpunk-qdrant] [--http-port 6333] [--grpc-port 6334] \
         [--binaries-dir qdrant-binaries] [--force] [--require-empty-dir]

示例:
  $0 --remote-host 1.2.3.4 --install-dir /pixelpunk-qdrant --http-port 7333 --grpc-port 7334
  # 强制跳过确认（非交互部署）
  $0 --remote-host 1.2.3.4 --install-dir /pixelpunk-qdrant --force
  # 要求安装目录必须为空，否则退出
  $0 --remote-host 1.2.3.4 --install-dir /pixelpunk-qdrant --require-empty-dir
EOF
}

# 参数解析
while [[ $# -gt 0 ]]; do
  case "$1" in
    --remote-host) REMOTE_HOST="$2"; shift 2;;
    --remote-user) REMOTE_USER="$2"; shift 2;;
    --ssh-port) SSH_PORT="$2"; shift 2;;
    --install-dir) INSTALL_DIR="$2"; shift 2;;
    --http-port) HTTP_PORT="$2"; shift 2;;
    --grpc-port) GRPC_PORT="$2"; shift 2;;
    --binaries-dir) BINARIES_DIR="$2"; shift 2;;
    --force) FORCE="yes"; shift 1;;
    --require-empty-dir) REQUIRE_EMPTY_DIR="yes"; shift 1;;
    -h|--help) usage; exit 0;;
    *) print_err "未知参数: $1"; usage; exit 1;;
  esac
done

[[ -z "$REMOTE_HOST" ]] && { print_err "必须指定 --remote-host"; usage; exit 1; }

# 为 ssh 与 scp 分别设置端口参数（ssh 用 -p，scp 用 -P）
SSH_COMMON_OPTS=( -o StrictHostKeyChecking=no -o ConnectTimeout=10 )
SSH_CMD=( ssh -p "$SSH_PORT" "${SSH_COMMON_OPTS[@]}" "${REMOTE_USER}@${REMOTE_HOST}" )
SCP_CMD=( scp -P "$SSH_PORT" "${SSH_COMMON_OPTS[@]}" )

print_info "检测远端可达性..."
if ! "${SSH_CMD[@]}" "echo ok" >/dev/null 2>&1; then
  print_err "无法连接远端 ${REMOTE_USER}@${REMOTE_HOST}:${SSH_PORT}"; exit 1
fi
print_ok "远端连接正常"

print_info "检测远端架构..."
REMOTE_ARCH=$("${SSH_CMD[@]}" "uname -m" | tr -d '\r')
case "$REMOTE_ARCH" in
  x86_64) PLATFORM="linux-amd64";;
  aarch64|arm64) PLATFORM="linux-arm64";;
  *) print_err "不支持的远端架构: $REMOTE_ARCH"; exit 1;;
 esac
print_ok "远端架构: $REMOTE_ARCH (${PLATFORM})"

LOCAL_BIN_PATH="${BINARIES_DIR}/qdrant-${PLATFORM}"
if [[ ! -f "$LOCAL_BIN_PATH" ]]; then
  print_err "未找到本地预下载二进制: $LOCAL_BIN_PATH"
  print_warn "请将对应平台的 Qdrant ${QDRANT_VERSION} 可执行文件保存为 ${LOCAL_BIN_PATH} 并赋予可执行权限"
  exit 1
fi

print_info "安装前检查与提示..."
# 检查目录是否存在与是否为空
DIR_STATE=$("${SSH_CMD[@]}" "bash -lc 'if [ -d \"${INSTALL_DIR}\" ]; then if [ \"\$(ls -A ${INSTALL_DIR} 2>/dev/null | wc -l)\" -gt 0 ]; then echo not_empty; else echo empty; fi; else echo not_exist; fi'" | tr -d '\r')
print_info "远端目录状态: ${DIR_STATE}"

# 检测是否已有运行中的服务（systemd 或 pid 文件或端口）
SERVICE_NAME="qdrant-$(basename "${INSTALL_DIR}")-${HTTP_PORT}"
RUNNING_HINT=$("${SSH_CMD[@]}" "if command -v systemctl >/dev/null 2>&1 && systemctl is-active --quiet '${SERVICE_NAME}'; then echo systemd_active; elif [ -f '${INSTALL_DIR}/qdrant.pid' ]; then P=\$(cat '${INSTALL_DIR}/qdrant.pid' 2>/dev/null || true); if [ -n \"\$P\" ] && kill -0 \$P 2>/dev/null; then echo pid_running; exit 0; fi; elif command -v lsof >/dev/null 2>&1 && lsof -iTCP:${HTTP_PORT} -sTCP:LISTEN -Pn >/dev/null 2>&1; then echo port_listen; else echo none; fi" | tr -d '\r')
print_info "服务运行状态: ${RUNNING_HINT}"

if [[ "$FORCE" != "yes" ]]; then
  # 需要确认的场景：目录非空 或 检测到服务运行
  if [[ "$DIR_STATE" == "not_empty" || "$RUNNING_HINT" != "none" ]]; then
    print_warn "检测到潜在风险：目录非空或已有服务在运行。"
    if [[ "$REQUIRE_EMPTY_DIR" == "yes" && "$DIR_STATE" != "empty" ]]; then
      print_err "已启用 --require-empty-dir，但目录不是空目录。中止。"; exit 1
    fi
    echo -n "确认继续部署并覆盖文件/重启进程? (yes/NO): "
    read -r ANSWER
    if [[ "$ANSWER" != "yes" ]]; then
      print_err "用户取消操作。"; exit 1
    fi
  fi
else
  if [[ "$REQUIRE_EMPTY_DIR" == "yes" && "$DIR_STATE" != "empty" ]]; then
    print_err "已启用 --require-empty-dir，但目录不是空目录。中止。"; exit 1
  fi
fi

print_info "创建远端目录结构..."
REMOTE_BIN_DIR="${INSTALL_DIR}/bin"
REMOTE_CFG_DIR="${INSTALL_DIR}/config"
REMOTE_LOG_DIR="${INSTALL_DIR}/logs"
REMOTE_DATA_DIR="${INSTALL_DIR}/shared/data/qdrant"
REMOTE_BACKUP_DIR="${INSTALL_DIR}/shared/data/backups"
"${SSH_CMD[@]}" "mkdir -p '${REMOTE_BIN_DIR}' '${REMOTE_CFG_DIR}' '${REMOTE_LOG_DIR}' '${REMOTE_DATA_DIR}' '${REMOTE_BACKUP_DIR}'"
print_ok "目录已就绪"

print_info "上传 Qdrant 二进制... (${LOCAL_BIN_PATH} -> ${REMOTE_BIN_DIR}/qdrant)"
REMOTE_TMP_BIN="/tmp/qdrant-upload-$(date +%s)-$$"
"${SCP_CMD[@]}" "$LOCAL_BIN_PATH" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_TMP_BIN}"
"${SSH_CMD[@]}" "chmod +x '${REMOTE_TMP_BIN}' && mv -f '${REMOTE_TMP_BIN}' '${REMOTE_BIN_DIR}/qdrant'"
print_ok "二进制上传完成"

print_info "生成并上传配置文件... (${REMOTE_CFG_DIR}/qdrant.yaml)"
QDRANT_CFG_CONTENT=$(cat <<EOF
storage:
  storage_path: ${REMOTE_DATA_DIR}
  write_consistency_factor: 1
  snapshots_path: ${REMOTE_DATA_DIR}/snapshots
service:
  http_port: ${HTTP_PORT}
  grpc_port: ${GRPC_PORT}
  host: 0.0.0.0
  enable_cors: true
log_level: INFO
cluster:
  enabled: false
optimizer:
  indexing_threshold: 100
  flush_interval_sec: 3
  vacuum_min_vector_number: 50
  max_optimization_threads: 1
wal:
  wal_capacity_mb: 128
  wal_segments_ahead: 4
  wal_sync_interval_ms: 100
hnsw_config:
  m: 16
  ef_construct: 128
  full_scan_threshold: 1000
quantization:
  scalar:
    type: int8
    always_ram: true
telemetry_disabled: true
EOF
)
# 使用 ssh here-doc 写入远端文件
"${SSH_CMD[@]}" "cat > '${REMOTE_CFG_DIR}/qdrant.yaml' <<'CFGEOF'
${QDRANT_CFG_CONTENT}
CFGEOF"
print_ok "配置文件已上传"

# 生成远端服务管理脚本（start/stop/restart/status/logs）
print_info "生成服务管理脚本... (${INSTALL_DIR}/qdrant-service.sh)"
# 通过本地 heredoc 输送到远端，避免本地 $ 变量误展开
"${SSH_CMD[@]}" "cat > '${INSTALL_DIR}/qdrant-service.sh'" <<'SVC'
#!/bin/bash
set -euo pipefail

DIR="__INSTALL_DIR__"
BIN="$DIR/bin/qdrant"
CFG="$DIR/config/qdrant.yaml"
LOG="$DIR/logs/qdrant.log"
PID="$DIR/qdrant.pid"

start() {
  mkdir -p "$(dirname "$LOG")"
  if [ -f "$PID" ]; then
    local p=$(cat "$PID" 2>/dev/null || true)
    if [ -n "$p" ] && kill -0 "$p" 2>/dev/null; then
      echo "Qdrant 已在运行 (PID: $p)"; return 0
    fi
  fi
  echo "启动 Qdrant..."
  nohup "$BIN" --config-path "$CFG" >> "$LOG" 2>&1 & echo $! > "$PID"
  sleep 2
  status
}

stop() {
  if [ -f "$PID" ]; then
    local p=$(cat "$PID" 2>/dev/null || true)
    if [ -n "$p" ] && kill -0 "$p" 2>/dev/null; then
      echo "优雅停止 PID: $p"
      kill -TERM "$p" 2>/dev/null || true
      sleep 3
      kill -0 "$p" 2>/dev/null && { echo "强制终止 PID: $p"; kill -9 "$p" 2>/dev/null || true; }
    fi
    rm -f "$PID"
  else
    echo "未找到 PID 文件，认为未运行"
  fi
}

restart() { stop; sleep 2; start; }

status() {
  if [ -f "$PID" ]; then
    local p=$(cat "$PID" 2>/dev/null || true)
    if [ -n "$p" ] && kill -0 "$p" 2>/dev/null; then
      echo "运行中 (PID: $p)"; return 0
    fi
  fi
  echo "未运行"
}

logs() { tail -n 200 -f "$LOG"; }

usage() {
  echo "用法: $0 {start|stop|restart|status|logs}"
}

case "${1:-}" in
  start) start ;;
  stop) stop ;;
  restart) restart ;;
  status) status ;;
  logs) logs ;;
  *) usage ;;
 esac
SVC
"${SSH_CMD[@]}" "sed -i 's#__INSTALL_DIR__#'\"${INSTALL_DIR}\"'#g' '${INSTALL_DIR}/qdrant-service.sh' && chmod +x '${INSTALL_DIR}/qdrant-service.sh'"
print_ok "服务管理脚本已生成"

# 配置自启（优先使用 systemd，若不可用则使用 @reboot cron 兜底）
print_info "检测 systemd 支持并配置自启动..."
SYSTEMD_OK=$("${SSH_CMD[@]}" "[ -d /run/systemd/system ] && command -v systemctl >/dev/null && echo yes || echo no" | tr -d '\r')
if [[ "$SYSTEMD_OK" == "yes" ]]; then
  SERVICE_NAME="qdrant-$(basename "${INSTALL_DIR}")-${HTTP_PORT}"
  UNIT_CONTENT=$(cat <<EOF
[Unit]
Description=Qdrant (${SERVICE_NAME})
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=${INSTALL_DIR}
ExecStart=${REMOTE_BIN_DIR}/qdrant --config-path ${REMOTE_CFG_DIR}/qdrant.yaml
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF
)
  UNIT_TMP="$(mktemp -t qdrant-unit.XXXXXX)"
  printf "%s\n" "$UNIT_CONTENT" > "$UNIT_TMP"
  REMOTE_UNIT_TMP="/tmp/${SERVICE_NAME}.service"
  "${SCP_CMD[@]}" "$UNIT_TMP" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_UNIT_TMP}"
  rm -f "$UNIT_TMP"
  REMOTE_UNIT_DEST="/etc/systemd/system/${SERVICE_NAME}.service"
  # 在远端一次性执行，使用位置参数避免引用问题
  "${SSH_CMD[@]}" bash -s -- "${REMOTE_UNIT_TMP}" "${REMOTE_UNIT_DEST}" "${SERVICE_NAME}" <<'EOF'
set -euo pipefail
src="$1"
dst="$2"
svc="$3"
mv -f "$src" "$dst"
systemctl daemon-reload
systemctl enable "$svc" || true
EOF
  print_ok "systemd 单元已安装: ${SERVICE_NAME}"
else
  print_warn "未检测到 systemd，使用 crontab @reboot 兜底"
  # 使用 bash -s 避免嵌套引号问题
  "${SSH_CMD[@]}" bash -s -- "${INSTALL_DIR}" "${HTTP_PORT}" <<'EOF'
set -euo pipefail
install_dir="$1"
port="$2"
base_name=$(basename "$install_dir")
mark="# qdrant_autostart $base_name $port"
# 移除旧的自启行
(crontab -l 2>/dev/null | grep -v "qdrant_autostart $base_name $port" || true) | crontab -
# 添加新的自启行
(crontab -l 2>/dev/null; echo "@reboot $install_dir/qdrant-service.sh start $mark") | crontab -
EOF
  print_ok "已添加/更新 crontab @reboot 自启"
fi

# 启动/重启 Qdrant
if [[ "$SYSTEMD_OK" == "yes" ]]; then
  print_info "通过 systemd 启动/重启 Qdrant ..."
  "${SSH_CMD[@]}" "systemctl restart ${SERVICE_NAME}"
else
  print_info "检查并停止已运行的 Qdrant (若存在)..."
  "${SSH_CMD[@]}" "if [ -f '${INSTALL_DIR}/qdrant.pid' ]; then PID=\$(cat '${INSTALL_DIR}/qdrant.pid' 2>/dev/null || true); if [ -n \"\$PID\" ] && kill -0 \$PID 2>/dev/null; then kill -TERM \$PID; sleep 3; fi; fi" || true

  print_info "通过 nohup 启动 Qdrant..."
  START_CMD="nohup '${REMOTE_BIN_DIR}/qdrant' --config-path '${REMOTE_CFG_DIR}/qdrant.yaml' >> '${REMOTE_LOG_DIR}/qdrant.log' 2>&1 & echo \$! > '${INSTALL_DIR}/qdrant.pid'"
  "${SSH_CMD[@]}" "$START_CMD"
fi

print_info "健康检查 http://127.0.0.1:${HTTP_PORT}/ ..."
HEALTH_OK=false
for i in {1..30}; do
  STATUS=$("${SSH_CMD[@]}" "bash -lc 'which curl >/dev/null 2>&1 && curl -s -o /dev/null -w \"%{http_code}\" http://127.0.0.1:${HTTP_PORT}/ || echo 000'" | tr -d '\r')
  if [[ "$STATUS" == "200" ]]; then HEALTH_OK=true; break; fi
  sleep 1
done

if ! $HEALTH_OK; then
  print_err "Qdrant 健康检查失败，请查看远端日志: ${REMOTE_LOG_DIR}/qdrant.log"
  exit 1
fi
print_ok "Qdrant 启动成功 (端口: ${HTTP_PORT}/${GRPC_PORT})"

cat <<OUT

${GREEN}部署完成${NC}
- 安装目录: ${INSTALL_DIR}
- 配置文件: ${INSTALL_DIR}/config/qdrant.yaml
- 数据目录: ${INSTALL_DIR}/shared/data/qdrant
- 日志文件: ${INSTALL_DIR}/logs/qdrant.log
- 健康检查: http://127.0.0.1:${HTTP_PORT}/

${CYAN}后端配置建议:${NC}
- 在项目 config.yaml 中设置 vector.qdrant_url: "http://127.0.0.1:${HTTP_PORT}"
- 多环境/多实例：为每个环境指定不同 --install-dir 和端口 (--http-port/--grpc-port)

${CYAN}管理与自启:${NC}
- 管理脚本: ${INSTALL_DIR}/qdrant-service.sh {start|stop|restart|status|logs}
- 若检测到 systemd: 已安装并启用服务单元 (名称: qdrant-$(basename "${INSTALL_DIR}")-${HTTP_PORT})，可用:
  - systemctl status qdrant-$(basename "${INSTALL_DIR}")-${HTTP_PORT}
  - systemctl restart qdrant-$(basename "${INSTALL_DIR}")-${HTTP_PORT}
- 若无 systemd: 已添加 crontab @reboot 自动启动

OUT

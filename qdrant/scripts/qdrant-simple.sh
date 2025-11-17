#!/bin/bash

# PixelPunk Qdrant 简化管理脚本
# 专注于核心功能：启动、停止、状态检查

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
QDRANT_ROOT="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$QDRANT_ROOT")"
QDRANT_VERSION="v1.15.4"
QDRANT_DIR="$QDRANT_ROOT/bin"
QDRANT_DATA_DIR="$QDRANT_ROOT/data/qdrant"
QDRANT_CONFIG="$QDRANT_ROOT/config/qdrant.yaml"
QDRANT_PORT="6333"
QDRANT_GRPC_PORT="6334"

# 获取系统架构
get_arch() {
    case "$(uname -m)" in
        x86_64) echo "x86_64" ;;
        arm64|aarch64) echo "aarch64" ;;
        *) echo "unsupported" ;;
    esac
}

# 获取系统平台
get_platform() {
    case "$(uname -s)" in
        Darwin) echo "apple-darwin" ;;
        Linux) echo "unknown-linux-gnu" ;;
        *) echo "unsupported" ;;
    esac
}

# 下载Qdrant
download_qdrant() {
    local arch=$(get_arch)
    local platform=$(get_platform)
    
    if [ "$arch" = "unsupported" ] || [ "$platform" = "unsupported" ]; then
        echo -e "${RED}不支持的系统架构或平台${NC}"
        exit 1
    fi
    
    local binary_name="qdrant-${arch}-${platform}"
    local download_url="https://github.com/qdrant/qdrant/releases/download/${QDRANT_VERSION}/${binary_name}.tar.gz"
    
    echo -e "${GREEN}下载 Qdrant ${QDRANT_VERSION}...${NC}"
    echo "架构: ${arch}-${platform}"
    echo "下载地址: ${download_url}"
    
    # 创建bin目录
    mkdir -p "${QDRANT_DIR}"
    
    # 下载并解压
    if curl -L -o "${QDRANT_DIR}/qdrant.tar.gz" "${download_url}"; then
        cd "${QDRANT_DIR}"
        tar -xzf qdrant.tar.gz
        rm qdrant.tar.gz
        chmod +x qdrant
        echo -e "${GREEN}Qdrant 下载并解压成功!${NC}"
    else
        echo -e "${RED}Qdrant 下载失败${NC}"
        exit 1
    fi
}

# 创建配置文件
create_config() {
    echo -e "${GREEN}创建 Qdrant 配置文件...${NC}"
    
    mkdir -p "$(dirname "${QDRANT_CONFIG}")"
    mkdir -p "${QDRANT_DATA_DIR}"
    
    cat > "${QDRANT_CONFIG}" << EOF
storage:
  storage_path: ./data/qdrant
  write_consistency_factor: 1
  snapshots_path: ./data/qdrant/snapshots

service:
  http_port: ${QDRANT_PORT}
  grpc_port: ${QDRANT_GRPC_PORT}
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
    
    echo -e "${GREEN}配置文件创建完成: ${QDRANT_CONFIG}${NC}"
}

# 简化启动函数
start_qdrant() {
    # 检查是否已有进程在运行
    if pgrep -f "qdrant" > /dev/null; then
        echo -e "${RED}检测到 Qdrant 进程已在运行${NC}"
        echo "PID: $(pgrep -f qdrant)"
        echo -e "${YELLOW}请先运行 '$0 stop' 停止现有进程${NC}"
        exit 1
    fi
    
    # 检查端口是否被占用
    if lsof -i ":${QDRANT_PORT}" > /dev/null 2>&1; then
        echo -e "${RED}端口 ${QDRANT_PORT} 已被占用${NC}"
        echo "占用进程:"
        lsof -i ":${QDRANT_PORT}"
        exit 1
    fi
    
    # 检查二进制文件
    if [ ! -f "${QDRANT_DIR}/qdrant" ]; then
        echo -e "${YELLOW}Qdrant 二进制文件不存在，开始下载...${NC}"
        download_qdrant
    fi
    
    # 确保数据目录和快照目录存在
    mkdir -p "${QDRANT_DATA_DIR}"
    mkdir -p "${QDRANT_DATA_DIR}/snapshots"
    
    # 创建配置文件（如果不存在）
    if [ ! -f "${QDRANT_CONFIG}" ]; then
        echo -e "${YELLOW}配置文件不存在，创建默认配置...${NC}"
        create_config
    fi
    
    echo -e "${GREEN}启动 Qdrant 向量数据库...${NC}"
    echo "HTTP 端口: ${QDRANT_PORT}"
    echo "gRPC 端口: ${QDRANT_GRPC_PORT}"
    echo "数据目录: ${QDRANT_DATA_DIR}"
    echo "配置文件: ${QDRANT_CONFIG}"
    echo ""
    echo -e "${YELLOW}使用 Ctrl+C 停止服务${NC}"
    echo ""
    
    # 启动Qdrant (切换到qdrant目录，使相对路径正确工作)
    cd "${QDRANT_ROOT}"
    exec "./bin/qdrant" --config-path "./config/qdrant.yaml"
}

# 停止Qdrant
stop_qdrant() {
    local pids=$(pgrep -f "qdrant")
    if [ -n "$pids" ]; then
        echo -e "${YELLOW}正在停止 Qdrant 进程...${NC}"
        for pid in $pids; do
            echo "停止进程 PID: $pid"
            kill -TERM "$pid" 2>/dev/null || kill -9 "$pid" 2>/dev/null
        done
        
        sleep 2
        
        # 确认进程已停止
        if ! pgrep -f "qdrant" > /dev/null; then
            echo -e "${GREEN}Qdrant 已成功停止${NC}"
        else
            echo -e "${RED}部分进程可能仍在运行，请手动检查${NC}"
        fi
    else
        echo -e "${YELLOW}没有发现运行中的 Qdrant 进程${NC}"
    fi
}

# 检查状态
status_qdrant() {
    echo -e "${GREEN}=== Qdrant 服务状态 ===${NC}"
    
    # 检查进程
    local pids=$(pgrep -f "qdrant")
    if [ -n "$pids" ]; then
        echo -e "${GREEN}✓ Qdrant 进程正在运行${NC}"
        for pid in $pids; do
            echo "  PID: $pid"
        done
    else
        echo -e "${RED}✗ Qdrant 进程未运行${NC}"
    fi
    
    # 检查端口
    if lsof -i ":${QDRANT_PORT}" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ HTTP 端口 ${QDRANT_PORT} 正在监听${NC}"
    else
        echo -e "${RED}✗ HTTP 端口 ${QDRANT_PORT} 未监听${NC}"
    fi
    
    if lsof -i ":${QDRANT_GRPC_PORT}" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ gRPC 端口 ${QDRANT_GRPC_PORT} 正在监听${NC}"
    else
        echo -e "${RED}✗ gRPC 端口 ${QDRANT_GRPC_PORT} 未监听${NC}"
    fi
    
    # 检查数据目录
    if [ -d "${QDRANT_DATA_DIR}" ]; then
        local size=$(du -sh "${QDRANT_DATA_DIR}" 2>/dev/null | cut -f1)
        echo -e "${GREEN}✓ 数据目录存在: ${QDRANT_DATA_DIR} (${size})${NC}"
    else
        echo -e "${RED}✗ 数据目录不存在: ${QDRANT_DATA_DIR}${NC}"
    fi
    
    # 检查配置文件
    if [ -f "${QDRANT_CONFIG}" ]; then
        echo -e "${GREEN}✓ 配置文件存在: ${QDRANT_CONFIG}${NC}"
    else
        echo -e "${RED}✗ 配置文件不存在: ${QDRANT_CONFIG}${NC}"
    fi
    
    echo ""
    echo -e "${GREEN}Web UI: http://localhost:${QDRANT_PORT}/dashboard${NC}"
}

# 重启
restart_qdrant() {
    echo -e "${YELLOW}重启 Qdrant...${NC}"
    stop_qdrant
    sleep 2
    start_qdrant
}

# 帮助信息
show_help() {
    echo "Qdrant 简化管理脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 start       启动 Qdrant 服务"
    echo "  $0 stop        停止 Qdrant 服务"
    echo "  $0 restart     重启 Qdrant 服务"
    echo "  $0 status      检查服务状态"
    echo "  $0 download    下载 Qdrant 二进制文件"
    echo "  $0 config      创建配置文件"
    echo "  $0 help        显示此帮助"
    echo ""
    echo "配置:"
    echo "  版本: ${QDRANT_VERSION}"
    echo "  端口: ${QDRANT_PORT} (HTTP), ${QDRANT_GRPC_PORT} (gRPC)"
    echo "  数据目录: ${QDRANT_DATA_DIR}"
    echo "  配置文件: ${QDRANT_CONFIG}"
}

# 主逻辑
case "${1:-help}" in
    start)
        start_qdrant
        ;;
    stop)
        stop_qdrant
        ;;
    restart)
        restart_qdrant
        ;;
    status)
        status_qdrant
        ;;
    download)
        download_qdrant
        ;;
    config)
        create_config
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo -e "${RED}未知命令: $1${NC}"
        show_help
        exit 1
        ;;
esac


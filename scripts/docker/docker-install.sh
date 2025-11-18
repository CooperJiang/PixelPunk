#!/bin/bash

# PixelPunk Docker 一键安装脚本
# 自动下载配置文件并启动 Docker Compose 服务
#
# 使用方式：
#   curl -fsSL https://download.pixelpunk.cc/shell/docker-install.sh | bash
#
# 或者：
#   wget -qO- https://download.pixelpunk.cc/shell/docker-install.sh | bash

set -e

# 下载地址
DOCKER_COMPOSE_URL="https://download.pixelpunk.cc/docker/docker-compose.yml"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# 默认配置
DEFAULT_PORT=9520
DEFAULT_DATA_DIR="./data"
DEFAULT_LOGS_DIR="./logs"
DEFAULT_UPLOADS_DIR="./uploads"

# 用户配置变量
USER_PORT=""
USER_DATA_DIR=""
USER_LOGS_DIR=""
USER_UPLOADS_DIR=""

# 交互式输入工具
prompt_input() {
    local prompt_text="$1"
    local default_value="$2"
    local input=""

    if [ -t 0 ]; then
        read -r -p "$prompt_text" input || true
    elif [ -e /dev/tty ]; then
        read -r -p "$prompt_text" input < /dev/tty || true
    else
        echo -e "${YELLOW}未检测到交互式终端，使用默认值: ${default_value}${NC}"
        input=""
    fi

    if [ -z "$input" ]; then
        input="$default_value"
    fi

    echo "$input"
}

prompt_confirm() {
    local prompt_text="$1"
    local default_answer="$2"
    local answer
    answer=$(prompt_input "$prompt_text" "$default_answer")
    echo "$answer"
}

# 打印欢迎信息
print_header() {
    clear
    echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}                                                           ${NC}"
    echo -e "${CYAN}    ██████╗ ██╗██╗  ██╗███████╗██╗     ██████╗ ██╗   ██╗███╗   ██╗██╗  ██╗${NC}"
    echo -e "${CYAN}    ██╔══██╗██║╚██╗██╔╝██╔════╝██║     ██╔══██╗██║   ██║████╗  ██║██║ ██╔╝${NC}"
    echo -e "${CYAN}    ██████╔╝██║ ╚███╔╝ █████╗  ██║     ██████╔╝██║   ██║██╔██╗ ██║█████╔╝ ${NC}"
    echo -e "${CYAN}    ██╔═══╝ ██║ ██╔██╗ ██╔══╝  ██║     ██╔═══╝ ██║   ██║██║╚██╗██║██╔═██╗ ${NC}"
    echo -e "${CYAN}    ██║     ██║██╔╝ ██╗███████╗███████╗██║     ╚██████╔╝██║ ╚████║██║  ██╗${NC}"
    echo -e "${CYAN}    ╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝╚═╝      ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝${NC}"
    echo -e "${CYAN}                                                           ${NC}"
    echo -e "${CYAN}               智能资源管理平台 - Docker 一键安装           ${NC}"
    echo -e "${CYAN}                                                           ${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"
    echo ""
}

# 检查依赖
check_dependencies() {
    echo -e "${BLUE}[1/7] 检查系统依赖...${NC}"

    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}✗ 未找到 Docker${NC}"
        echo -e "${YELLOW}请先安装 Docker: https://docs.docker.com/get-docker/${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Docker 已安装: $(docker --version | cut -d ' ' -f3 | cut -d ',' -f1)${NC}"

    # 检查 docker-compose
    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}✗ 未找到 docker-compose${NC}"
        echo -e "${YELLOW}请先安装 docker-compose: https://docs.docker.com/compose/install/${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ docker-compose 已安装: $(docker-compose --version | cut -d ' ' -f4 | cut -d ',' -f1)${NC}"

    # 检查 curl 或 wget
    if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
        echo -e "${RED}✗ 未找到 curl 或 wget${NC}"
        echo -e "${YELLOW}请先安装 curl 或 wget${NC}"
        exit 1
    fi
    if command -v curl &> /dev/null; then
        echo -e "${GREEN}✓ curl 已安装${NC}"
    else
        echo -e "${GREEN}✓ wget 已安装${NC}"
    fi

    echo ""
}

# 下载配置文件
download_configs() {
    echo -e "${BLUE}[2/7] 下载配置文件...${NC}"

    # 下载 docker-compose.yml
    echo -ne "${CYAN}下载 docker-compose.yml${NC}"
    if [ -f "docker-compose.yml" ]; then
        echo -e " ${YELLOW}(已存在，备份为 docker-compose.yml.bak)${NC}"
        mv docker-compose.yml docker-compose.yml.bak
    else
        echo ""
    fi

    if command -v curl &> /dev/null; then
        if curl -fsSL "${DOCKER_COMPOSE_URL}" -o docker-compose.yml; then
            echo -e "${GREEN}✓ docker-compose.yml 下载成功${NC}"
        else
            echo -e "${RED}✗ docker-compose.yml 下载失败${NC}"
            exit 1
        fi
    else
        if wget -q "${DOCKER_COMPOSE_URL}" -O docker-compose.yml; then
            echo -e "${GREEN}✓ docker-compose.yml 下载成功${NC}"
        else
            echo -e "${RED}✗ docker-compose.yml 下载失败${NC}"
            exit 1
        fi
    fi

    # 创建 configs 目录
    mkdir -p configs

    # 检查并删除错误的目录（如果存在）
    if [ -d "configs/config.docker.yaml" ]; then
        echo -e "${YELLOW}检测到错误的配置目录，正在清理...${NC}"
        rm -rf "configs/config.docker.yaml"
    fi

    # 创建 config.docker.yaml（如果不存在）
    if [ ! -f "configs/config.docker.yaml" ]; then
        echo -ne "${CYAN}创建 config.docker.yaml${NC}"
        cat > configs/config.docker.yaml << 'CONFIGEOF'
# PixelPunk Docker 环境配置
# 用于 Docker Compose 部署

app:
  port: 9520
  mode: "release"
  ns: "PixelPunk"

# 数据库配置（使用 Docker 容器名）
database:
  type: "mysql"
  host: "mysql"
  port: 3306
  username: "pixelpunk"
  password: "pixelpunk_pass"
  name: "pixelpunk"
  max_idle_conns: 10
  max_open_conns: 100

# Redis 配置（使用 Docker 容器名）
redis:
  host: "redis"
  port: 6379
  password: ""
  db: 0
  pool_size: 10

# 向量数据库配置（使用 Docker 容器名）
vector:
  enabled: true
  qdrant_url: "http://qdrant:6333"
  timeout: 30
  collection_name: "pixelpunk_images"

# 日志配置
log:
  level: "info"
  output: "stdout"
  format: "json"

# 文件上传配置
upload:
  max_size: 100
  allowed_types:
    - "image/jpeg"
    - "image/png"
    - "image/gif"
    - "image/webp"
    - "image/bmp"
  save_path: "/app/uploads"

# AI 配置
ai:
  enabled: false
  provider: "openai"
  api_key: ""
  base_url: ""
  model: "gpt-4-vision-preview"

# 安全配置
security:
  jwt_secret: "pixelpunk-change-this-secret-key"
  session_timeout: 7200
CONFIGEOF
        echo -e " ${GREEN}✓${NC}"
    else
        echo -e "${GREEN}✓ configs/config.docker.yaml 已存在${NC}"
    fi

    # 验证配置文件是文件而不是目录
    if [ ! -f "configs/config.docker.yaml" ]; then
        echo -e "${RED}✗ 配置文件创建失败或不是有效文件${NC}"
        exit 1
    fi

    echo ""
}

# 交互式配置
interactive_config() {
    echo -e "${BLUE}[3/7] 配置服务参数...${NC}"
    echo -e "${YELLOW}提示：直接按回车使用默认值${NC}"
    echo ""

    # 端口配置
    echo -e "${BOLD}端口配置：${NC}"
    USER_PORT=$(prompt_input "$(echo -e ${CYAN}PixelPunk Web 访问端口 [默认: ${DEFAULT_PORT}]: ${NC})" "$DEFAULT_PORT")
    echo -e "${YELLOW}注意：数据库、Redis、Qdrant 等内部服务仅在容器网络内通信，不对外暴露端口${NC}"

    echo ""
    echo -e "${BOLD}数据目录配置：${NC}"
    USER_DATA_DIR=$(prompt_input "$(echo -e ${CYAN}数据存储目录 [默认: ${DEFAULT_DATA_DIR}]: ${NC})" "$DEFAULT_DATA_DIR")

    USER_LOGS_DIR=$(prompt_input "$(echo -e ${CYAN}日志存储目录 [默认: ${DEFAULT_LOGS_DIR}]: ${NC})" "$DEFAULT_LOGS_DIR")

    USER_UPLOADS_DIR=$(prompt_input "$(echo -e ${CYAN}上传文件目录 [默认: ${DEFAULT_UPLOADS_DIR}]: ${NC})" "$DEFAULT_UPLOADS_DIR")

    echo ""
    echo -e "${GREEN}配置已完成！${NC}"
    echo ""
}

# 显示配置摘要
show_summary() {
    echo -e "${BLUE}[4/7] 配置摘要：${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "  Web 访问地址："
    echo -e "    PixelPunk:  ${GREEN}http://localhost:${USER_PORT}${NC}"
    echo ""
    echo -e "  内部服务（仅容器内通信）："
    echo -e "    MySQL:      容器内 3306 端口"
    echo -e "    Qdrant:     容器内 6333 端口"
    echo -e "    Redis:      容器内 6379 端口"
    echo ""
    echo -e "  数据目录："
    echo -e "    数据:       ${USER_DATA_DIR}"
    echo -e "    日志:       ${USER_LOGS_DIR}"
    echo -e "    上传:       ${USER_UPLOADS_DIR}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    confirm=$(prompt_confirm "$(echo -e ${YELLOW}确认以上配置并继续？[Y/n]: ${NC})" "Y")
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        echo -e "${RED}安装已取消${NC}"
        exit 0
    fi
    echo ""
}

# 创建目录并设置权限
setup_directories() {
    echo -e "${BLUE}[5/7] 创建数据目录并设置权限...${NC}"

    # 需要创建的目录列表
    DIRS=(
        "${USER_DATA_DIR}"
        "${USER_LOGS_DIR}"
        "${USER_UPLOADS_DIR}"
        "${USER_DATA_DIR}/mysql"
        "${USER_DATA_DIR}/qdrant"
        "${USER_DATA_DIR}/redis"
        "./configs"
    )

    # 检查是否需要 sudo
    NEED_SUDO=false
    for dir in "${DIRS[@]}"; do
        parent_dir=$(dirname "${dir}")
        if [ ! -w "${parent_dir}" ] 2>/dev/null && [ ! -d "${dir}" ]; then
            NEED_SUDO=true
            break
        fi
    done

    if [ "$NEED_SUDO" = true ]; then
        echo -e "${YELLOW}部分目录需要管理员权限，可能会提示输入密码${NC}"
        SUDO="sudo"
    else
        SUDO=""
    fi

    # 创建目录
    for dir in "${DIRS[@]}"; do
        if [ -d "${dir}" ]; then
            echo -e "  ${GREEN}✓${NC} ${dir} (已存在)"
        else
            if $SUDO mkdir -p "${dir}" 2>/dev/null; then
                echo -e "  ${GREEN}✓${NC} ${dir} (已创建)"
            else
                echo -e "  ${RED}✗${NC} ${dir} (创建失败)"
                exit 1
            fi
        fi
    done

    # 设置权限
    echo ""
    echo -e "${BLUE}设置目录权限...${NC}"
    if [ -n "$SUDO" ]; then
        $SUDO chown -R $USER:$(id -gn) "${USER_DATA_DIR}" "${USER_LOGS_DIR}" "${USER_UPLOADS_DIR}" 2>/dev/null || true
    fi
    chmod -R 755 "${USER_DATA_DIR}" "${USER_LOGS_DIR}" "${USER_UPLOADS_DIR}" 2>/dev/null || \
        $SUDO chmod -R 755 "${USER_DATA_DIR}" "${USER_LOGS_DIR}" "${USER_UPLOADS_DIR}"
    echo -e "${GREEN}✓ 权限设置完成${NC}"

    echo ""
}

# 启动服务
start_services() {
    echo -e "${BLUE}[6/7] 启动 Docker 服务...${NC}"
    echo ""

    # 启动前再次验证配置文件（防止被误删或变成目录）
    if [ ! -f "configs/config.docker.yaml" ]; then
        echo -e "${RED}✗ 配置文件不存在或不是有效文件${NC}"
        echo -e "${YELLOW}正在重新创建配置文件...${NC}"

        # 如果是目录，先删除
        if [ -d "configs/config.docker.yaml" ]; then
            rm -rf "configs/config.docker.yaml"
        fi

        # 重新生成配置文件
        download_configs
    fi

    # 停止已有服务
    if docker-compose ps 2>/dev/null | grep -q "Up"; then
        echo -e "${YELLOW}检测到已运行的服务，正在停止...${NC}"
        docker-compose down
        echo ""
    fi

    # 检测并更新 PixelPunk 镜像
    echo -e "${CYAN}检查 PixelPunk 镜像更新...${NC}"

    PIXELPUNK_IMAGE="snine98/pixelpunk:latest"

    # 获取本地镜像 ID
    LOCAL_IMAGE_ID=$(docker images -q ${PIXELPUNK_IMAGE} 2>/dev/null)

    if [ -n "$LOCAL_IMAGE_ID" ]; then
        echo -e "${YELLOW}检测到本地镜像，正在检查远程是否有新版本...${NC}"
    else
        echo -e "${YELLOW}本地无镜像，正在下载...${NC}"
    fi

    # 拉取最新镜像（Docker 会自动判断是否需要下载）
    if docker pull ${PIXELPUNK_IMAGE} 2>&1 | tee /tmp/pixelpunk_pull.log | grep -q "Image is up to date"; then
        echo -e "${GREEN}✓ 已是最新版本${NC}"
    elif grep -q "Downloaded newer image" /tmp/pixelpunk_pull.log || grep -q "Status: Downloaded" /tmp/pixelpunk_pull.log; then
        echo -e "${GREEN}✓ 已更新到最新版本${NC}"
    else
        echo -e "${GREEN}✓ 镜像已准备就绪${NC}"
    fi

    rm -f /tmp/pixelpunk_pull.log
    echo ""

    # 启动服务
    echo -e "${CYAN}正在启动服务...${NC}"
    if PORT=$USER_PORT \
       DATA_DIR=$USER_DATA_DIR \
       LOGS_DIR=$USER_LOGS_DIR \
       UPLOADS_DIR=$USER_UPLOADS_DIR \
       docker-compose up -d; then
        echo ""
        echo -e "${GREEN}✓ 服务启动成功${NC}"
        echo ""
    else
        echo ""
        echo -e "${RED}✗ 服务启动失败${NC}"
        echo ""
        echo -e "${YELLOW}诊断信息：${NC}"
        echo -e "1. 检查配置文件："
        echo -e "   ls -la configs/config.docker.yaml"
        ls -la configs/config.docker.yaml 2>&1 | sed 's/^/   /'
        echo ""
        echo -e "2. 检查容器状态："
        docker-compose ps 2>&1 | sed 's/^/   /'
        echo ""
        echo -e "3. 查看错误日志："
        echo -e "   ${CYAN}docker-compose logs pixelpunk${NC}"
        echo ""
        exit 1
    fi
}

# 等待服务就绪
wait_for_services() {
    echo -e "${BLUE}[7/7] 等待服务就绪...${NC}"
    echo ""

    # 等待 MySQL
    echo -ne "${CYAN}MySQL 数据库启动中${NC}"
    for i in {1..30}; do
        if docker-compose exec -T mysql mysqladmin ping -h localhost -u root -proot_password_change_me &>/dev/null; then
            echo -e " ${GREEN}✓${NC}"
            break
        fi
        echo -ne "."
        sleep 2
    done

    # 等待 PixelPunk
    echo -ne "${CYAN}PixelPunk 应用启动中${NC}"
    for i in {1..30}; do
        if curl -s http://localhost:${USER_PORT}/health &>/dev/null; then
            echo -e " ${GREEN}✓${NC}"
            break
        fi
        echo -ne "."
        sleep 2
    done

    echo ""
}

# 显示完成信息
show_completion() {
    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}             ✓ PixelPunk 安装完成！             ${NC}"
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${BOLD}访问地址：${NC}"
    echo -e "  ${CYAN}Web 界面:${NC}  ${GREEN}http://localhost:${USER_PORT}${NC}"
    echo -e "  ${CYAN}Health:${NC}    http://localhost:${USER_PORT}/health"
    echo ""
    echo -e "${BOLD}常用命令：${NC}"
    echo -e "  查看日志:     ${YELLOW}docker-compose logs -f pixelpunk${NC}"
    echo -e "  停止服务:     ${YELLOW}docker-compose stop${NC}"
    echo -e "  重启服务:     ${YELLOW}docker-compose restart${NC}"
    echo -e "  停止并删除:   ${YELLOW}docker-compose down${NC}"
    echo ""
    echo -e "${BOLD}数据目录：${NC}"
    echo -e "  数据文件:     ${USER_DATA_DIR}"
    echo -e "  日志文件:     ${USER_LOGS_DIR}"
    echo -e "  上传文件:     ${USER_UPLOADS_DIR}"
    echo ""
    echo -e "${YELLOW}提示：首次运行需要进行系统初始化，请访问 Web 界面完成设置${NC}"
    echo ""
}

# 主流程
main() {
    print_header
    check_dependencies
    download_configs
    interactive_config
    show_summary
    setup_directories
    start_services
    wait_for_services
    show_completion
}

# 错误处理
trap 'echo -e "\n${RED}安装过程中出现错误${NC}"; exit 1' ERR

# 执行主流程
main

#!/usr/bin/env bash

# PixelPunk 一键安装脚本
# 用途：自动检测平台并下载安装 PixelPunk
# 使用：curl -fsSL http://download.pixelpunk.cc/shell/setup.sh | bash
#      或 wget -qO- http://download.pixelpunk.cc/shell/setup.sh | bash

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# 配置
DOWNLOAD_BASE_URL="http://download.pixelpunk.cc/release"
DEFAULT_VERSION="v1.0.0"
DEFAULT_PORT=9520
# 默认直接安装到当前目录（不创建子目录）
INSTALL_DIR="${PIXELPUNK_INSTALL_DIR:-$(pwd)}"

# 日志函数（输出到 stderr，避免污染函数返回值）
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" >&2
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1" >&2
}

log_error() {
    echo -e "${RED}[✗]${NC} $1" >&2
}

log_warning() {
    echo -e "${YELLOW}[!]${NC} $1" >&2
}

log_step() {
    echo -e "${CYAN}▶${NC} ${BOLD}$1${NC}" >&2
}

# 打印标题
print_banner() {
    clear
    echo -e "${BLUE}" >&2
    cat << 'EOF' >&2
╔════════════════════════════════════════╗
║     PixelPunk 一键安装脚本             ║
║     快速部署您的图床系统               ║
╚════════════════════════════════════════╝
EOF
    echo -e "${NC}" >&2
}

# 检测操作系统
detect_os() {
    local os=""
    case "$(uname -s)" in
        Linux*)
            os="linux"
            ;;
        Darwin*)
            os="darwin"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            log_error "不支持在 Windows 上直接运行此脚本"
            log_info "请下载 Windows 安装包手动安装："
            log_info "${DOWNLOAD_BASE_URL}/pixelpunk-${DEFAULT_VERSION}-windows-amd64.zip"
            exit 1
            ;;
        *)
            log_error "不支持的操作系统: $(uname -s)"
            exit 1
            ;;
    esac
    echo "$os"
}

# 检测架构
detect_arch() {
    local arch=""
    case "$(uname -m)" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        armv7l|armv6l)
            log_error "不支持 32 位 ARM 架构"
            log_info "请使用 64 位系统 (arm64/aarch64)"
            exit 1
            ;;
        i386|i686)
            log_error "不支持 32 位 x86 架构"
            log_info "请使用 64 位系统 (x86_64/amd64)"
            exit 1
            ;;
        *)
            log_error "不支持的架构: $(uname -m)"
            exit 1
            ;;
    esac
    echo "$arch"
}

# 检查依赖
check_dependencies() {
    local missing_deps=()

    # 检查下载工具
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        missing_deps+=("curl 或 wget")
    fi

    # 检查解压工具
    if ! command -v tar >/dev/null 2>&1; then
        missing_deps+=("tar")
    fi

    if [ ${#missing_deps[@]} -gt 0 ]; then
        log_error "缺少必要的依赖工具："
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        echo ""
        log_info "请先安装依赖，例如："
        if [[ "$(detect_os)" == "linux" ]]; then
            echo "  Ubuntu/Debian: sudo apt-get install curl tar"
            echo "  CentOS/RHEL:   sudo yum install curl tar"
        else
            echo "  macOS:         brew install curl"
        fi
        exit 1
    fi
}

# 下载文件
download_file() {
    local url=$1
    local output=$2

    echo "" >&2
    if command -v curl >/dev/null 2>&1; then
        # curl 显示进度条
        curl -L --progress-bar "$url" -o "$output"
    elif command -v wget >/dev/null 2>&1; then
        # wget 显示进度条（移除 -q 以显示进度）
        wget --show-progress --progress=bar:force "$url" -O "$output"
    else
        log_error "未找到下载工具 (curl/wget)"
        exit 1
    fi
    echo "" >&2
}

# 检查版本是否存在
check_version_exists() {
    local url=$1

    if command -v curl >/dev/null 2>&1; then
        if curl --output /dev/null --silent --head --fail "$url"; then
            return 0
        fi
    elif command -v wget >/dev/null 2>&1; then
        if wget --spider "$url" 2>/dev/null; then
            return 0
        fi
    fi
    return 1
}

# 检查端口是否被占用
check_port() {
    local port=$1

    # 使用 netstat 或 ss 检查端口
    if command -v netstat >/dev/null 2>&1; then
        if netstat -tln 2>/dev/null | grep -q ":${port} "; then
            return 1  # 端口被占用
        fi
    elif command -v ss >/dev/null 2>&1; then
        if ss -tln 2>/dev/null | grep -q ":${port} "; then
            return 1  # 端口被占用
        fi
    elif command -v lsof >/dev/null 2>&1; then
        if lsof -i ":${port}" >/dev/null 2>&1; then
            return 1  # 端口被占用
        fi
    fi

    return 0  # 端口可用
}

# 配置端口
configure_port() {
    local selected_port=""
    local port_ok=false

    while [ "$port_ok" = false ]; do
        # 检查默认端口
        if check_port "$DEFAULT_PORT"; then
            log_success "默认端口 $DEFAULT_PORT 可用"
            echo "" >&2

            # 询问用户是否使用默认端口
            local use_default=""
            if [ -t 0 ]; then
                read -p "$(echo -e ${CYAN}使用默认端口 $DEFAULT_PORT? [Y/n]: ${NC})" use_default
            elif [ -e /dev/tty ]; then
                read -p "$(echo -e ${CYAN}使用默认端口 $DEFAULT_PORT? [Y/n]: ${NC})" use_default < /dev/tty
            else
                use_default="Y"
            fi

            if [[ "$use_default" =~ ^[Nn]$ ]]; then
                # 用户选择自定义端口
                local custom_port=""
                if [ -t 0 ]; then
                    read -p "$(echo -e ${CYAN}请输入自定义端口号 [1024-65535]: ${NC})" custom_port
                elif [ -e /dev/tty ]; then
                    read -p "$(echo -e ${CYAN}请输入自定义端口号 [1024-65535]: ${NC})" custom_port < /dev/tty
                else
                    custom_port="$DEFAULT_PORT"
                fi

                # 验证端口号
                if [[ ! "$custom_port" =~ ^[0-9]+$ ]] || [ "$custom_port" -lt 1024 ] || [ "$custom_port" -gt 65535 ]; then
                    log_error "无效的端口号，使用默认端口 $DEFAULT_PORT"
                    selected_port="$DEFAULT_PORT"
                    port_ok=true
                elif check_port "$custom_port"; then
                    log_success "端口 $custom_port 可用"
                    selected_port="$custom_port"
                    port_ok=true
                else
                    log_error "端口 $custom_port 已被占用，请重新选择"
                fi
            else
                selected_port="$DEFAULT_PORT"
                port_ok=true
            fi
        else
            log_warning "默认端口 $DEFAULT_PORT 已被占用"
            echo "" >&2

            # 默认端口被占用，要求用户输入
            local custom_port=""
            if [ -t 0 ]; then
                read -p "$(echo -e ${CYAN}请输入可用端口号 [1024-65535]: ${NC})" custom_port
            elif [ -e /dev/tty ]; then
                read -p "$(echo -e ${CYAN}请输入可用端口号 [1024-65535]: ${NC})" custom_port < /dev/tty
            else
                log_error "默认端口被占用且无法读取用户输入"
                exit 1
            fi

            # 验证端口号
            if [[ ! "$custom_port" =~ ^[0-9]+$ ]] || [ "$custom_port" -lt 1024 ] || [ "$custom_port" -gt 65535 ]; then
                log_error "无效的端口号，请重新输入"
            elif check_port "$custom_port"; then
                log_success "端口 $custom_port 可用"
                selected_port="$custom_port"
                port_ok=true
            else
                log_error "端口 $custom_port 已被占用，请重新选择"
            fi
        fi
    done

    echo "" >&2
    echo "$selected_port"  # 返回值，必须输出到 stdout
}

# 主安装流程
main() {
    print_banner

    echo "" >&2
    log_step "步骤 1/7: 检测系统环境"
    echo "" >&2

    # 检测平台
    local os=$(detect_os)
    local arch=$(detect_arch)
    local platform="${os}-${arch}"

    log_success "操作系统: $os"
    log_success "架构: $arch"
    log_success "目标平台: $platform"

    # 检查依赖
    check_dependencies
    log_success "依赖检查通过"

    echo "" >&2
    log_step "步骤 2/7: 配置服务端口"
    echo "" >&2

    # 配置端口
    local service_port=$(configure_port)
    log_info "选定端口: $service_port"

    echo "" >&2
    log_step "步骤 3/7: 准备安装"
    echo "" >&2

    # 使用默认版本（当前固定为 v1.0.0）
    # TODO: 未来有多版本时，可添加版本选择功能
    version="$DEFAULT_VERSION"
    log_info "安装版本: $version"

    # 构建文件名和 URL
    local archive_name="pixelpunk-${version}-${platform}.tar.gz"
    local download_url="${DOWNLOAD_BASE_URL}/${archive_name}"

    echo "" >&2
    log_step "步骤 4/7: 检查版本可用性"
    echo "" >&2

    log_info "检查: $download_url"
    if ! check_version_exists "$download_url"; then
        log_error "版本不存在或无法访问"
        log_info "请检查版本号是否正确，或访问以下地址查看可用版本："
        log_info "${DOWNLOAD_BASE_URL}/"
        exit 1
    fi
    log_success "版本可用"

    echo "" >&2
    log_step "步骤 5/7: 下载安装包"
    echo "" >&2

    # 创建临时目录
    local tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" EXIT

    local tmp_archive="${tmp_dir}/${archive_name}"

    log_info "下载地址: $download_url"
    log_info "保存位置: $tmp_archive"

    echo "" >&2
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" >&2
    echo -e "${CYAN}📥 正在下载，请稍候...${NC}" >&2
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" >&2

    if download_file "$download_url" "$tmp_archive"; then
        echo ""
        local file_size=$(du -h "$tmp_archive" | cut -f1)
        log_success "下载完成 (大小: $file_size)"
    else
        echo ""
        log_error "下载失败"
        log_info "请检查网络连接或稍后重试"
        exit 1
    fi

    echo "" >&2
    log_step "步骤 6/7: 解压安装包"
    echo "" >&2

    # 检查当前目录
    log_info "安装目录: $INSTALL_DIR"

    # 检查当前目录是否已有 PixelPunk 安装
    local has_pixelpunk=false
    if [ -f "$INSTALL_DIR/pixelpunk" ] || [ -f "$INSTALL_DIR/pixelpunk.sh" ]; then
        has_pixelpunk=true
    fi

    # 检查当前目录是否为空（忽略隐藏文件）
    local file_count=$(ls -A "$INSTALL_DIR" 2>/dev/null | grep -v '^\.' | wc -l)

    if [ "$has_pixelpunk" = true ]; then
        log_warning "检测到当前目录已有 PixelPunk 安装"
        echo ""
        echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" >&2
        echo -e "${YELLOW}⚠️  是否覆盖现有安装?${NC}" >&2
        echo -e "${YELLOW}    现有数据（数据库、上传文件等）将被保留，但程序文件会被替换。${NC}" >&2
        echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" >&2
        echo ""

        local overwrite=""
        if [ -t 0 ]; then
            read -p "$(echo -e ${CYAN}覆盖安装? [y/N]: ${NC})" overwrite
        elif [ -e /dev/tty ]; then
            read -p "$(echo -e ${CYAN}覆盖安装? [y/N]: ${NC})" overwrite < /dev/tty
        else
            log_error "无法读取用户输入，且当前目录已有安装"
            log_info "如需强制覆盖，请设置环境变量："
            echo "  PIXELPUNK_FORCE_INSTALL=1"
            exit 1
        fi

        if [[ ! "$overwrite" =~ ^[Yy]$ ]]; then
            echo ""
            log_info "取消安装"
            exit 0
        fi
        echo ""
    elif [ "$file_count" -gt 0 ]; then
        log_warning "当前目录不为空"
        echo ""
        echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" >&2
        echo -e "${YELLOW}⚠️  当前目录包含其他文件${NC}" >&2
        echo -e "${YELLOW}    建议创建一个空目录进行安装，例如：${NC}" >&2
        echo -e "${CYAN}      mkdir pixelpunk && cd pixelpunk${NC}" >&2
        echo -e "${CYAN}      curl -fsSL https://download.pixelpunk.cc/shell/install.sh | bash${NC}" >&2
        echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" >&2
        echo ""

        local continue_install=""
        if [ -t 0 ]; then
            read -p "$(echo -e ${CYAN}仍要继续安装到当前目录? [y/N]: ${NC})" continue_install
        elif [ -e /dev/tty ]; then
            read -p "$(echo -e ${CYAN}仍要继续安装到当前目录? [y/N]: ${NC})" continue_install < /dev/tty
        else
            log_error "当前目录不为空，无法自动安装"
            exit 1
        fi

        if [[ ! "$continue_install" =~ ^[Yy]$ ]]; then
            echo ""
            log_info "取消安装"
            exit 0
        fi
        echo ""
    fi

    mkdir -p "$INSTALL_DIR"

    # 创建解压目录
    local extract_dir="${tmp_dir}/extract"
    mkdir -p "$extract_dir"

    # 解压到临时目录
    log_info "正在解压..."
    if tar -xzf "$tmp_archive" -C "$extract_dir" 2>&1 | grep -v "Ignoring unknown extended header"; then
        :  # 解压成功，忽略 macOS 扩展属性警告
    fi

    # 验证解压结果
    if [ ! -f "$extract_dir/pixelpunk" ]; then
        log_error "解压失败：未找到 pixelpunk 可执行文件"
        log_info "临时目录内容："
        ls -la "$extract_dir" || true
        exit 1
    fi

    # 复制文件到安装目录
    log_info "复制文件到: $INSTALL_DIR"
    cp -r "$extract_dir"/* "$INSTALL_DIR/"

    log_success "解压完成"

    echo "" >&2
    log_step "步骤 7/7: 运行安装脚本"
    echo "" >&2

    # 切换到安装目录
    cd "$INSTALL_DIR" || {
        log_error "无法切换到安装目录: $INSTALL_DIR"
        exit 1
    }

    # 检查 install.sh 是否存在
    if [ ! -f "install.sh" ]; then
        log_error "安装脚本不存在: $INSTALL_DIR/install.sh"
        log_error "安装包可能不完整"
        exit 1
    fi

    # 执行安装脚本
    chmod +x install.sh
    log_info "执行安装脚本..."
    echo "" >&2
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" >&2
    echo "" >&2

    # 通过环境变量传递端口配置
    PIXELPUNK_PORT="$service_port" ./install.sh

    echo "" >&2
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" >&2
    echo "" >&2
    log_success "安装完成！"

    echo "" >&2
    echo -e "${GREEN}╔════════════════════════════════════════╗${NC}" >&2
    echo -e "${GREEN}║     🎉 PixelPunk 安装成功！            ║${NC}" >&2
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}" >&2
    echo "" >&2
    echo -e "${YELLOW}📍 安装位置:${NC} $INSTALL_DIR"
    echo "" >&2
    echo -e "${YELLOW}📋 常用命令:${NC}" >&2
    echo "  ./pixelpunk.sh status   # 查看服务状态"
    echo "  ./pixelpunk.sh logs     # 查看日志"
    echo "  ./pixelpunk.sh restart  # 重启服务"
    echo "" >&2
}

# 错误处理
handle_error() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        echo ""
        log_error "安装过程中发生错误 (退出码: $exit_code)"
        log_info "如需帮助，请访问: https://github.com/CooperJiang/PixelPunk-v1/issues"
    fi
    exit $exit_code
}

trap 'handle_error' ERR

# 运行主程序
main "$@"

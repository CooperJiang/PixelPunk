#!/bin/bash

# PixelPunk 多架构 Docker 镜像构建和推送脚本
# 支持 linux/amd64, linux/arm64, linux/arm/v7
#
# 使用方式：
#   ./scripts/docker/build-push-multiarch.sh
#
# 环境变量：
#   DOCKER_VERSION           - 版本标签（默认: latest）
#   DOCKER_USERNAME          - Docker Hub 用户名（默认: snine98）
#   PLATFORMS                - 目标平台列表（默认: linux/amd64,linux/arm64,linux/arm/v7）
#   DOCKER_BUILD_CACHE_MODE  - local / registry / none（默认: local）
#   DOCKER_BUILD_CACHE_DIR   - 本地缓存目录（默认: .buildx-cache）
#   DOCKER_BUILD_CACHE_REF   - 远程缓存镜像引用（默认: <image>:buildcache）

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Docker Hub 配置
DOCKER_USERNAME="${DOCKER_USERNAME:-snine98}"
IMAGE_NAME="pixelpunk"
FULL_IMAGE_NAME="${DOCKER_USERNAME}/${IMAGE_NAME}"

# 版本号（可通过环境变量覆盖）
VERSION="${DOCKER_VERSION:-latest}"

# 支持的平台架构（可通过 PLATFORMS 环境变量覆盖）
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64,linux/arm/v7}"

# Buildx 缓存配置
# 注意：多平台构建时 local cache 导出非常慢（可能卡数小时）
# Dockerfile 中的 --mount=type=cache 已经提供足够好的缓存性能
CACHE_MODE="${DOCKER_BUILD_CACHE_MODE:-none}" # local / registry / none (推荐 none)
CACHE_DIR="${DOCKER_BUILD_CACHE_DIR:-.buildx-cache}"
CACHE_REF="${DOCKER_BUILD_CACHE_REF:-${FULL_IMAGE_NAME}:buildcache}"
CACHE_ARGS=()
BASE_BUILD_ARGS=()

echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${CYAN}          PixelPunk 多架构镜像构建和推送${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}镜像信息：${NC}"
echo -e "  镜像名称: ${BOLD}${FULL_IMAGE_NAME}${NC}"
echo -e "  版本标签: ${BOLD}${VERSION}${NC}"
echo -e "  支持平台: ${BOLD}${PLATFORMS}${NC}"
case "${CACHE_MODE}" in
    local)
        echo -e "  缓存模式: ${BOLD}local → ${CACHE_DIR}${NC}"
        ;;
    registry)
        echo -e "  缓存模式: ${BOLD}registry → ${CACHE_REF}${NC}"
        ;;
    *)
        echo -e "  缓存模式: ${BOLD}none${NC}"
        ;;
esac
echo ""

# 检查 Docker
check_docker() {
    echo -e "${BLUE}[1/6] 检查 Docker 环境...${NC}"

    if ! command -v docker &> /dev/null; then
        echo -e "${RED}✗ Docker 未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Docker 已安装: $(docker --version | cut -d ' ' -f3 | cut -d ',' -f1)${NC}"

    # 检查 Docker daemon 是否运行
    if ! docker info &> /dev/null; then
        echo -e "${RED}✗ Docker daemon 未运行${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Docker daemon 运行正常${NC}"

    echo ""
}

# 检查并设置 buildx
setup_buildx() {
    echo -e "${BLUE}[2/6] 设置 Docker Buildx...${NC}"

    # 检查 buildx 是否可用
    if ! docker buildx version &> /dev/null; then
        echo -e "${RED}✗ Docker Buildx 不可用${NC}"
        echo -e "${YELLOW}请更新到 Docker 19.03+ 版本${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Docker Buildx 可用${NC}"

    # 创建或使用 multiarch builder
    BUILDER_NAME="pixelpunk-multiarch"
    if ! docker buildx inspect "${BUILDER_NAME}" &> /dev/null; then
        echo -e "${CYAN}创建多架构 builder: ${BUILDER_NAME}${NC}"
        docker buildx create --name "${BUILDER_NAME}" --use --platform "${PLATFORMS}"
    else
        echo -e "${CYAN}使用现有 builder: ${BUILDER_NAME}${NC}"
        docker buildx use "${BUILDER_NAME}"
    fi

    # 启动 builder
    docker buildx inspect --bootstrap
    echo -e "${GREEN}✓ Buildx builder 已就绪${NC}"

    echo ""
}


# 当前登录用户
get_logged_in_user() {
    docker info --format '{{ .RegistryConfig.IndexConfigs."docker.io".Name }}' 2>/dev/null | tr -d '\n'
}

# 登录 Docker Hub（仅推送时调用）
docker_login() {
    local current_user
    current_user=$(get_logged_in_user)

    if [[ -n "${current_user}" ]]; then
        echo -e "${GREEN}✓ 当前已登录 Docker Hub (${current_user})${NC}"
        if [[ "${current_user}" != "${DOCKER_USERNAME}" ]]; then
            echo -e "${YELLOW}⚠️  目标推送账号为 ${DOCKER_USERNAME}${NC}"
            read -r -p "是否切换账号？[y/N] " relogin_choice
            if [[ ! ${relogin_choice:-N} =~ ^[Yy]$ ]]; then
                echo -e "${YELLOW}继续使用现有账号推送${NC}"
                echo ""
                return 0
            fi
        else
            echo ""
            return 0
        fi
    else
        echo -e "${YELLOW}⚠️  尚未登录到 Docker Hub${NC}"
    fi

    echo -e "${CYAN}正在登录到 Docker Hub (${DOCKER_USERNAME})...${NC}"
    if ! docker login --username "${DOCKER_USERNAME}"; then
        echo -e "${RED}✗ 登录失败，请稍后重试${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ 登录成功${NC}"
    echo ""
}

# Docker 构建会在容器内重新构建前端，无需预检查

# 配置 Buildx 缓存
prepare_cache() {
    echo -e "${BLUE}[3/6] 准备 Buildx 缓存 (${CACHE_MODE})...${NC}"
    case "${CACHE_MODE}" in
        local)
            mkdir -p "${CACHE_DIR}"
            CACHE_ARGS+=(--cache-from type=local,src=${CACHE_DIR})
            CACHE_ARGS+=(--cache-to type=local,dest=${CACHE_DIR},mode=max)
            echo -e "${GREEN}✓ 使用本地缓存目录 ${CACHE_DIR}${NC}"
            ;;
        registry)
            CACHE_ARGS+=(--cache-from type=registry,ref=${CACHE_REF})
            CACHE_ARGS+=(--cache-to type=registry,ref=${CACHE_REF},mode=max)
            echo -e "${GREEN}✓ 使用远程缓存 ${CACHE_REF}${NC}"
            ;;
        none)
            echo -e "${YELLOW}⚠️  已禁用缓存${NC}"
            ;;
        *)
            echo -e "${RED}未知缓存模式: ${CACHE_MODE}${NC}"
            exit 1
            ;;
    esac
    echo ""
}


# 构建多架构镜像（不推送）
build_image() {
    echo -e "${BLUE}[4/6] 构建多架构镜像...${NC}"
    echo -e "${YELLOW}支持的平台: ${PLATFORMS}${NC}"
    echo ""

    BASE_BUILD_ARGS=(
        --platform "${PLATFORMS}"
        --tag "${FULL_IMAGE_NAME}:${VERSION}"
        --tag "${FULL_IMAGE_NAME}:latest"
        --file Dockerfile
    )

    BUILD_START=$(date +%s)

    docker buildx build \
        "${BASE_BUILD_ARGS[@]}" \
        "${CACHE_ARGS[@]}" \
        .

    BUILD_END=$(date +%s)
    BUILD_TIME=$((BUILD_END - BUILD_START))

    echo ""
    echo -e "${GREEN}✓ 构建完成（耗时: ${BUILD_TIME}s）${NC}"
    echo ""
}

push_image() {
    echo -e "${BLUE}[5/6] 推送镜像到 Docker Hub...${NC}"
    docker_login
    docker buildx build \
        "${BASE_BUILD_ARGS[@]}" \
        "${CACHE_ARGS[@]}" \
        --push \
        .
    echo ""
    show_completion
}

maybe_push() {
    read -r -p "镜像已构建，是否推送到 ${FULL_IMAGE_NAME}? [y/N] " push_choice
    if [[ ${push_choice:-N} =~ ^[Yy]$ ]]; then
        push_image
    else
        echo ""
        echo -e "${YELLOW}已跳过推送。如需推送，可重新执行 make docker-build${NC}"
        echo ""
    fi
}

# 显示完成信息
show_completion() {
    echo -e "${BLUE}[6/6] 验证镜像...${NC}"

    # 显示镜像信息
    echo -e "${CYAN}查询镜像清单...${NC}"
    docker buildx imagetools inspect "${FULL_IMAGE_NAME}:${VERSION}" | head -20

    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}          ✓ 多架构镜像构建完成！${NC}"
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${BOLD}镜像信息：${NC}"
    echo -e "  仓库: ${GREEN}${FULL_IMAGE_NAME}${NC}"
    echo -e "  标签: ${GREEN}${VERSION}, latest${NC}"
    echo ""
    echo -e "${BOLD}支持的系统架构：${NC}"
    echo -e "  ${GREEN}✓${NC} linux/amd64    - Intel/AMD 64位系统（服务器、PC）"
    echo -e "  ${GREEN}✓${NC} linux/arm64    - ARM 64位系统（树莓派4+、Mac M1/M2/M3、ARM服务器）"
    echo -e "  ${GREEN}✓${NC} linux/arm/v7   - ARM 32位系统（树莓派3、旧版ARM设备）"
    echo ""
    echo -e "${BOLD}使用方式：${NC}"
    echo -e "  docker pull ${FULL_IMAGE_NAME}:${VERSION}"
    echo -e "  docker run -p 9520:9520 ${FULL_IMAGE_NAME}:${VERSION}"
    echo ""
    echo -e "${BOLD}Docker Compose：${NC}"
    echo -e "  curl -fsSL https://download.pixelpunk.cc/shell/docker-install.sh | bash"
    echo ""
    echo -e "${YELLOW}注意：用户在不同架构系统上拉取镜像时，Docker 会自动选择对应的架构版本${NC}"
    echo ""
}

# 错误处理
trap 'echo -e "\n${RED}构建过程中出现错误${NC}"; exit 1' ERR

# 主流程
main() {
    check_docker
    setup_buildx
    prepare_cache
    build_image
    maybe_push
}

# 执行主流程
main

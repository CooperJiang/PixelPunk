#!/bin/bash
# 本地构建 Docker 镜像（不推送到 Docker Hub）

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

DOCKER_USERNAME="${DOCKER_USERNAME:-snine98}"
IMAGE_NAME="pixelpunk"
FULL_IMAGE_NAME="${DOCKER_USERNAME}/${IMAGE_NAME}"
VERSION="${DOCKER_VERSION:-latest}"
PLATFORM="${PLATFORM:-linux/amd64}"

echo -e "${CYAN}═══════════════════════════════════════${NC}"
echo -e "${CYAN}  PixelPunk 本地镜像构建${NC}"
echo -e "${CYAN}═══════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}镜像信息：${NC}"
echo -e "  镜像名称: ${BOLD}${FULL_IMAGE_NAME}${NC}"
echo -e "  版本标签: ${BOLD}${VERSION}${NC}"
echo -e "  目标平台: ${BOLD}${PLATFORM}${NC}"
echo -e "  构建模式: ${BOLD}本地构建（不推送）${NC}"
echo ""

echo -e "${BLUE}开始构建镜像...${NC}"

docker buildx build \
  --platform "${PLATFORM}" \
  --tag "${FULL_IMAGE_NAME}:${VERSION}" \
  --tag "${FULL_IMAGE_NAME}:latest" \
  --load \
  -f Dockerfile \
  .

echo ""
echo -e "${GREEN}✓ 镜像构建成功！${NC}"
echo ""
echo -e "${BLUE}查看镜像：${NC}"
docker images | grep pixelpunk | head -3

echo ""
echo -e "${BLUE}使用镜像：${NC}"
echo -e "  docker run -d -p 9520:9520 ${FULL_IMAGE_NAME}:${VERSION}"
echo ""

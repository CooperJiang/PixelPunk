#!/bin/bash

# PixelPunk 发布镜像预下载脚本
# 一次性下载所有构建需要的镜像，后续可完全离线构建

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║        📦 PixelPunk 发布构建镜像预下载工具                 ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo

# 需要的基础镜像
REQUIRED_IMAGES=(
    "golang:1.23"
    "ubuntu:22.04"
    "alpine:latest"
)

echo -e "${YELLOW}需要下载的镜像:${NC}"
for img in "${REQUIRED_IMAGES[@]}"; do
    echo "  • $img"
done
echo

# 检查Docker环境
if ! command -v docker >/dev/null 2>&1; then
    echo -e "${RED}错误: Docker未安装${NC}"
    exit 1
fi

if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}错误: Docker服务未运行${NC}"
    exit 1
fi

echo -e "${BLUE}开始下载镜像...${NC}"
echo

success_count=0
total_count=${#REQUIRED_IMAGES[@]}

for img in "${REQUIRED_IMAGES[@]}"; do
    echo -e "${YELLOW}正在下载: $img${NC}"
    
    if docker pull "$img"; then
        echo -e "${GREEN}✓ $img 下载成功${NC}"
        ((success_count++))
    else
        echo -e "${RED}✗ $img 下载失败${NC}"
    fi
    echo
done

# 启用buildx（如果没有）
echo -e "${YELLOW}配置Docker Buildx...${NC}"
if docker buildx version >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Docker Buildx 已可用${NC}"
else
    echo -e "${YELLOW}启用Docker Buildx...${NC}"
    if docker buildx create --use >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Docker Buildx 已启用${NC}"
    else
        echo -e "${YELLOW}⚠️ Docker Buildx 配置可能有问题，但不影响基本功能${NC}"
    fi
fi

echo
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"

if [ $success_count -eq $total_count ]; then
    echo -e "${GREEN}🎉 所有镜像下载成功！($success_count/$total_count)${NC}"
    echo
    echo -e "${GREEN}✅ 现在可以完全离线构建了！${NC}"
    echo
    echo -e "${YELLOW}使用方法:${NC}"
    echo "  make release              # 交互式发布构建"
    echo "  ./scripts/build-webp-multiplatform.sh  # 直接多平台构建"
    echo
    echo -e "${BLUE}支持的平台:${NC}"
    echo "  🐧 Linux AMD64/ARM64 (完整WebP支持)"
    echo "  🍎 macOS Intel/Apple Silicon (智能WebP)"
    echo "  🪟 Windows AMD64 (智能WebP)"
else
    echo -e "${RED}⚠️ 部分镜像下载失败 ($success_count/$total_count)${NC}"
    echo
    echo -e "${YELLOW}你可以手动重试失败的镜像：${NC}"
    for img in "${REQUIRED_IMAGES[@]}"; do
        if ! docker images | grep -q "$img"; then
            echo "  docker pull $img"
        fi
    done
fi

echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
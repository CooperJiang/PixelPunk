#!/bin/bash

# Docker构建资源清理脚本
# 用于清理构建过程中产生的临时容器和镜像

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║           🧹 Docker 构建资源清理工具                       ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo

echo -e "${YELLOW}正在清理 Docker 构建资源...${NC}"
echo

# 1. 清理 buildx builder 容器
echo -e "${BLUE}1. 清理 buildx builder 实例${NC}"
BUILDER_CONTAINERS=$(docker ps -a --format "{{.Names}}" | grep "buildx_buildkit" || true)
if [ -n "$BUILDER_CONTAINERS" ]; then
    echo "$BUILDER_CONTAINERS" | while read -r container; do
        if [ -n "$container" ]; then
            echo "  • 停止并移除: $container"
            docker stop "$container" > /dev/null 2>&1 || true
            docker rm "$container" > /dev/null 2>&1 || true
        fi
    done
    echo -e "${GREEN}  ✓ buildx 容器已清理${NC}"
else
    echo -e "${GREEN}  ✓ 没有需要清理的 buildx 容器${NC}"
fi

echo

# 2. 清理 buildx 构建器
echo -e "${BLUE}2. 清理 buildx 构建器${NC}"
BUILDERS=$(docker buildx ls 2>/dev/null | grep -v "default" | grep -E "linux-builder|pixelpunk" | awk '{print $1}' || true)
if [ -n "$BUILDERS" ]; then
    echo "$BUILDERS" | while read -r builder; do
        if [ -n "$builder" ] && [ "$builder" != "Name/Node" ]; then
            echo "  • 移除构建器: $builder"
            docker buildx rm "$builder" > /dev/null 2>&1 || true
        fi
    done
    echo -e "${GREEN}  ✓ 构建器已清理${NC}"
else
    echo -e "${GREEN}  ✓ 没有需要清理的构建器${NC}"
fi

echo

# 3. 清理构建镜像
echo -e "${BLUE}3. 清理构建相关镜像${NC}"
BUILD_IMAGES=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep -E "pixelpunk.*offline|pixelpunk.*webp" || true)
if [ -n "$BUILD_IMAGES" ]; then
    echo "$BUILD_IMAGES" | while read -r image; do
        if [ -n "$image" ]; then
            echo "  • 移除镜像: $image"
            docker rmi "$image" > /dev/null 2>&1 || true
        fi
    done
    echo -e "${GREEN}  ✓ 构建镜像已清理${NC}"
else
    echo -e "${GREEN}  ✓ 没有需要清理的构建镜像${NC}"
fi

echo

# 4. 清理未使用的资源
echo -e "${BLUE}4. 清理未使用的 Docker 资源${NC}"
echo "  • 清理未使用的容器..."
PRUNED_CONTAINERS=$(docker container prune -f 2>/dev/null | grep "Total reclaimed space" || echo "")
echo "  • 清理未使用的镜像..."
PRUNED_IMAGES=$(docker image prune -f 2>/dev/null | grep "Total reclaimed space" || echo "")

if [ -n "$PRUNED_CONTAINERS" ] || [ -n "$PRUNED_IMAGES" ]; then
    [ -n "$PRUNED_CONTAINERS" ] && echo "    $PRUNED_CONTAINERS"
    [ -n "$PRUNED_IMAGES" ] && echo "    $PRUNED_IMAGES"
else
    echo -e "${GREEN}  ✓ 没有可清理的资源${NC}"
fi

echo

# 5. 显示清理结果
echo -e "${BLUE}5. 清理完成状态${NC}"
REMAINING_BUILDERS=$(docker ps --format "{{.Names}}" | grep "buildx_buildkit" | wc -l || echo "0")
echo "  • 剩余 buildx 容器: $REMAINING_BUILDERS 个"

if [ "$REMAINING_BUILDERS" -eq 0 ]; then
    echo -e "${GREEN}✅ 所有构建资源已清理完成！${NC}"
else
    echo -e "${YELLOW}⚠️  仍有 $REMAINING_BUILDERS 个 buildx 容器运行${NC}"
    echo -e "${YELLOW}   这些可能是其他项目正在使用的构建器${NC}"
fi

echo
echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                     清理完成！                             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
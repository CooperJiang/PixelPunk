#!/bin/bash

# PixelPunk Qdrant 统一管理入口
# 调用 qdrant/ 目录下的管理脚本

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_ROOT="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$SCRIPTS_ROOT")"
QDRANT_SCRIPT="$PROJECT_ROOT/qdrant/scripts/qdrant-simple.sh"

# 检查脚本是否存在
if [ ! -f "$QDRANT_SCRIPT" ]; then
    echo -e "\033[0;31m错误: 找不到Qdrant管理脚本: $QDRANT_SCRIPT\033[0m"
    exit 1
fi

# 转发所有参数到实际脚本
exec "$QDRANT_SCRIPT" "$@"


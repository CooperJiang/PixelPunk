#!/bin/sh
set -e

# PixelPunk Docker 启动脚本
# 自动处理配置文件初始化和验证

CONFIG_FILE="/app/configs/config.yaml"
DOCKER_CONFIG="/app/configs/config.docker.yaml"
EXAMPLE_CONFIG="/app/configs/config.example.yaml"

echo "=== PixelPunk Docker Entrypoint ==="

# 检查配置文件状态
if [ -d "$CONFIG_FILE" ]; then
    echo "⚠️  警告: 配置文件被错误地创建为目录，正在修复..."
    rm -rf "$CONFIG_FILE"
fi

if [ ! -f "$CONFIG_FILE" ]; then
    echo "配置文件不存在，正在初始化..."

    # 优先使用 Docker 专用配置（包含数据库连接等完整配置）
    if [ -f "$DOCKER_CONFIG" ]; then
        cp "$DOCKER_CONFIG" "$CONFIG_FILE"
        echo "✓ 已使用 Docker 默认配置"
    elif [ -f "$EXAMPLE_CONFIG" ]; then
        cp "$EXAMPLE_CONFIG" "$CONFIG_FILE"
        echo "✓ 已使用示例配置"
    else
        echo "❌ 错误: 找不到任何配置模板"
        exit 1
    fi

    echo ""
    echo "📌 重要提示:"
    echo "  • 默认配置适用于 Docker Compose 环境"
    echo "  • 数据库连接: mysql:3306 (用户名: pixelpunk)"
    echo "  • 如需修改，请在宿主机编辑 configs/config.docker.yaml 后重启容器"
    echo ""
    echo "📖 文档: https://docs.pixelpunk.cc"
    echo "🚀 推荐安装方式: curl -fsSL https://download.pixelpunk.cc/shell/docker-install.sh | bash"
    echo ""
fi

# 验证配置文件是否可读
if [ ! -r "$CONFIG_FILE" ]; then
    echo "错误: 配置文件不可读: $CONFIG_FILE"
    exit 1
fi

echo "✓ 配置文件验证通过"
echo "启动 PixelPunk..."
echo ""

# 执行主程序
exec "$@"

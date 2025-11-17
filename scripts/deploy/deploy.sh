#!/bin/bash

# PixelPunk 部署管理系统启动脚本
# 调用真正的部署脚本

# 确保脚本可执行
chmod +x ./deploy/deploy-arrow.sh

# 启动交互式部署系统
./deploy/deploy-arrow.sh "$@"
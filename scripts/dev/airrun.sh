#!/bin/bash

# 检查air是否安装
if ! command -v air &> /dev/null; then
    echo "air 未安装，请先安装:"
    echo "go install github.com/cosmtrek/air@latest"
    exit 1
fi

# 启动air热重载
air
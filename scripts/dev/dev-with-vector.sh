#!/bin/bash

# PixelPunk 开发环境启动脚本（包含向量服务）
# 使用方法: ./scripts/dev-with-vector.sh

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}PixelPunk 开发环境启动（独立向量服务模式）${NC}"
echo

# 检查依赖
if ! command -v pnpm &> /dev/null; then
    echo -e "${RED}错误: pnpm 未安装${NC}"
    echo "请先安装 pnpm: npm install -g pnpm"
    exit 1
fi

# 记录进程ID以便清理
PIDS=()

# 清理函数
cleanup() {
    echo
    echo -e "${YELLOW}正在停止所有服务...${NC}"
    for pid in "${PIDS[@]}"; do
        if kill -0 $pid 2>/dev/null; then
            kill $pid
            echo "停止进程: $pid"
        fi
    done
    # 额外清理
    pkill -f "vector-service" 2>/dev/null || true
    pkill -f "qdrant" 2>/dev/null || true
    pkill -f "air" 2>/dev/null || true
    pkill -f "vite" 2>/dev/null || true
    echo -e "${GREEN}所有服务已停止${NC}"
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 启动向量服务
echo -e "${GREEN}1. 启动向量服务...${NC}"
cd vector-service
make dev &
VECTOR_PID=$!
PIDS+=($VECTOR_PID)
cd ..

echo -e "${YELLOW}等待向量服务启动...${NC}"
sleep 3

# 检查向量服务是否启动成功
if ! curl -s http://localhost:7333/api/v1/health > /dev/null; then
    echo -e "${RED}警告: 向量服务可能未正常启动${NC}"
    echo -e "${YELLOW}请检查 vector-service 目录下的配置和日志${NC}"
fi

# 启动后端服务
echo -e "${GREEN}2. 启动后端服务...${NC}"
./airrun.sh &
BACKEND_PID=$!
PIDS+=($BACKEND_PID)

echo -e "${YELLOW}等待后端服务启动...${NC}"
sleep 2

# 启动前端服务
echo -e "${GREEN}3. 启动前端服务...${NC}"
cd web
pnpm run dev &
FRONTEND_PID=$!
PIDS+=($FRONTEND_PID)
cd ..

echo
echo -e "${GREEN}=== 开发环境启动完成 ===${NC}"
echo -e "${YELLOW}前端地址: http://localhost:5173${NC}"
echo -e "${YELLOW}后端地址: http://localhost:9520${NC}"
echo -e "${YELLOW}向量服务: http://localhost:7333${NC}"
echo -e "${YELLOW}向量服务健康检查: http://localhost:7333/api/v1/health${NC}"
echo
echo -e "${GREEN}使用 Ctrl+C 停止所有服务${NC}"
echo

# 等待所有进程
wait
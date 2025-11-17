#!/bin/bash

# PixelPunk 快速部署脚本
# 使用已编译好的 pixelpunk-linux 二进制文件快速部署

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 配置
REMOTE_USER="root"
REMOTE_HOST="110.42.43.216"
REMOTE_PORT="22"
REMOTE_DIR="/pixelPunk/app"
LOCAL_BINARY="./pixelpunk-linux"

# 检查本地二进制文件
if [ ! -f "$LOCAL_BINARY" ]; then
    echo -e "${RED}✗ 未找到二进制文件: $LOCAL_BINARY${NC}"
    echo -e "${YELLOW}提示: 请先运行以下命令构建:${NC}"
    echo -e "  docker buildx build --platform linux/amd64 --load -t pixelpunk-build:latest ."
    echo -e "  docker create --name temp-pixelpunk pixelpunk-build:latest"
    echo -e "  docker cp temp-pixelpunk:/app/pixelpunk ./pixelpunk-linux"
    echo -e "  docker rm temp-pixelpunk"
    exit 1
fi

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║          🚀 PixelPunk 快速部署                             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# 显示文件信息
echo -e "${CYAN}二进制文件信息:${NC}"
ls -lh "$LOCAL_BINARY"
file "$LOCAL_BINARY"
echo ""

# 显示服务器信息
echo -e "${CYAN}目标服务器:${NC}"
echo -e "  主机: ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PORT}"
echo -e "  目录: ${REMOTE_DIR}"
echo ""

# 确认
read -p "确认部署? (y/N): " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}部署已取消${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}[1/5] 测试 SSH 连接...${NC}"
if ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} "echo '连接成功'" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ SSH 连接正常${NC}"
else
    echo -e "${RED}✗ SSH 连接失败${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}[2/5] 停止远程服务...${NC}"
ssh -o StrictHostKeyChecking=no -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} "
    cd ${REMOTE_DIR} && \

    # 尝试使用 service.sh 停止
    if [ -f ./service.sh ]; then
        ./service.sh stop 2>/dev/null || true
        sleep 1
    fi

    # 强制杀死所有 pixelpunk 进程
    pkill -9 pixelpunk 2>/dev/null || true
    killall -9 pixelpunk 2>/dev/null || true

    # 等待进程完全退出
    sleep 3

    # 验证没有残留进程
    if pgrep -f pixelpunk > /dev/null; then
        echo '警告: 仍有 pixelpunk 进程运行'
        ps aux | grep pixelpunk | grep -v grep
    else
        echo '所有 pixelpunk 进程已清理'
    fi
" && echo -e "${GREEN}✓ 服务已停止${NC}"

echo ""
echo -e "${BLUE}[3/5] 备份旧版本...${NC}"
ssh -o StrictHostKeyChecking=no -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} "
    cd ${REMOTE_DIR} && \
    if [ -f pixelpunk ]; then
        cp pixelpunk pixelpunk.backup.\$(date +%Y%m%d_%H%M%S)
        echo '备份文件: pixelpunk.backup.\$(date +%Y%m%d_%H%M%S)'
    else
        echo '无需备份（首次部署）'
    fi
" && echo -e "${GREEN}✓ 备份完成${NC}"

echo ""
echo -e "${BLUE}[4/5] 上传新版本...${NC}"
echo -e "${CYAN}上传中...${NC}"
if scp -o StrictHostKeyChecking=no -P ${REMOTE_PORT} "$LOCAL_BINARY" ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/pixelpunk; then
    echo -e "${GREEN}✓ 上传成功 ($(ls -lh $LOCAL_BINARY | awk '{print $5}'))${NC}"

    # 验证上传的文件
    echo -e "${CYAN}验证文件...${NC}"
    ssh -o StrictHostKeyChecking=no -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} "
        cd ${REMOTE_DIR} && \
        sync && \
        sleep 1 && \
        if [ -f pixelpunk ]; then
            ls -lh pixelpunk
            echo '✓ 文件存在且可访问'
        else
            echo '✗ 文件不存在！'
            exit 1
        fi
    " || { echo -e "${RED}✗ 文件验证失败${NC}"; exit 1; }
else
    echo -e "${RED}✗ 上传失败${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}[5/5] 设置权限并启动服务...${NC}"
ssh -o StrictHostKeyChecking=no -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} "
    cd ${REMOTE_DIR} && \

    # 设置执行权限
    chmod +x pixelpunk && \
    sync && \
    sleep 1 && \

    # 确保日志目录存在
    mkdir -p logs && \

    # 启动服务
    if [ -f ./service.sh ]; then
        echo '使用 service.sh 启动服务...'
        ./service.sh start
    else
        echo '使用 nohup 启动服务...'
        nohup ./pixelpunk > logs/pixelpunk.log 2>&1 &
        NOHUP_EXIT=\$?
        if [ \$NOHUP_EXIT -ne 0 ]; then
            echo \"启动失败，退出码: \$NOHUP_EXIT\"
            exit \$NOHUP_EXIT
        fi
        echo '启动命令已执行'
    fi
" && echo -e "${GREEN}✓ 启动命令已执行${NC}"

echo ""
echo -e "${CYAN}等待服务启动...${NC}"
sleep 3

# 检查服务状态
echo -e "${BLUE}检查服务状态...${NC}"
ssh -o StrictHostKeyChecking=no -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} "
    cd ${REMOTE_DIR} && \

    if [ -f ./service.sh ]; then
        ./service.sh status
    else
        if pgrep -f 'pixelpunk' > /dev/null; then
            echo -e '✓ 服务运行中'
            ps aux | grep pixelpunk | grep -v grep
        else
            echo -e '✗ 服务未运行'
            echo ''
            echo '最近的日志:'
            if [ -f logs/pixelpunk.log ]; then
                tail -20 logs/pixelpunk.log
            else
                echo '日志文件不存在'
            fi
            exit 1
        fi
    fi
" || {
    echo ""
    echo -e "${RED}✗ 服务启动失败${NC}"
    echo -e "${YELLOW}请手动检查日志: ssh ${REMOTE_USER}@${REMOTE_HOST} -p ${REMOTE_PORT} 'tail -50 ${REMOTE_DIR}/logs/pixelpunk.log'${NC}"
    exit 1
}

echo ""
echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✓ 部署完成！${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${CYAN}查看日志:${NC}"
echo -e "  ssh ${REMOTE_USER}@${REMOTE_HOST} -p ${REMOTE_PORT} 'tail -f ${REMOTE_DIR}/logs/pixelpunk.log'"
echo ""
echo -e "${CYAN}访问地址:${NC}"
echo -e "  http://${REMOTE_HOST}:9520"
echo ""

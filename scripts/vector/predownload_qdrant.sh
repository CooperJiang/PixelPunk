#!/bin/bash

# Qdrant 预下载脚本
# 用于预先下载 Qdrant 二进制文件到本地缓存，避免部署时网络问题

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 配置
QDRANT_VERSION="v1.12.5"
CACHE_DIR="./deploy/cache"

echo -e "${BLUE}=====================================${NC}"
echo -e "${BLUE}   Qdrant 预下载工具${NC}"
echo -e "${BLUE}   版本: $QDRANT_VERSION${NC}"
echo -e "${BLUE}=====================================${NC}"
echo

# 创建缓存目录
mkdir -p "$CACHE_DIR"

# 显示菜单
echo -e "${YELLOW}请选择要预下载的版本:${NC}"
echo -e "  1) Linux x86_64 (最常用)"
echo -e "  2) Linux ARM64"
echo -e "  3) 全部下载"
echo
read -p "请输入选项 [1-3]: " choice

# 下载函数
download_file() {
    local arch="$1"
    local url="$2"
    local filename="$3"
    local filepath="$CACHE_DIR/$filename"
    
    echo
    echo -e "${BLUE}正在下载 $arch 版本...${NC}"
    echo -e "${CYAN}URL: $url${NC}"
    echo -e "${CYAN}保存到: $filepath${NC}"
    
    # 检查是否已存在
    if [ -f "$filepath" ]; then
        echo -e "${YELLOW}文件已存在，是否重新下载? (y/n) [n]: ${NC}"
        read -n 1 redownload
        echo
        if [ "$redownload" != "y" ] && [ "$redownload" != "Y" ]; then
            echo -e "${GREEN}✅ 跳过下载，使用现有文件${NC}"
            return 0
        fi
        rm -f "$filepath"
    fi
    
    # 尝试下载
    echo -e "${BLUE}开始下载...${NC}"
    
    # 使用多种下载方式
    if command -v aria2c > /dev/null 2>&1; then
        # 使用 aria2 多线程下载（最快）
        echo -e "${CYAN}使用 aria2 多线程下载...${NC}"
        aria2c -x 16 -s 16 -k 1M -d "$CACHE_DIR" -o "$filename" "$url" || {
            echo -e "${RED}aria2 下载失败${NC}"
            return 1
        }
    elif command -v wget > /dev/null 2>&1; then
        # 使用 wget
        echo -e "${CYAN}使用 wget 下载...${NC}"
        wget --show-progress -O "$filepath" "$url" || {
            echo -e "${RED}wget 下载失败${NC}"
            return 1
        }
    elif command -v curl > /dev/null 2>&1; then
        # 使用 curl
        echo -e "${CYAN}使用 curl 下载...${NC}"
        curl -L -# -o "$filepath" "$url" || {
            echo -e "${RED}curl 下载失败${NC}"
            return 1
        }
    else
        echo -e "${RED}❌ 未找到下载工具（需要 aria2c、wget 或 curl）${NC}"
        echo -e "${YELLOW}提示: 安装 aria2 可以获得最快的下载速度:${NC}"
        echo -e "${BLUE}  macOS: brew install aria2${NC}"
        echo -e "${BLUE}  Linux: apt/yum install aria2${NC}"
        return 1
    fi
    
    # 验证下载
    if [ -f "$filepath" ]; then
        size=$(ls -lh "$filepath" | awk '{print $5}')
        echo -e "${GREEN}✅ 下载成功！文件大小: $size${NC}"
        return 0
    else
        echo -e "${RED}❌ 下载失败${NC}"
        return 1
    fi
}

# 执行下载
success_count=0
fail_count=0

case $choice in
    1)
        # Linux x86_64
        url="https://github.com/qdrant/qdrant/releases/download/${QDRANT_VERSION}/qdrant-x86_64-unknown-linux-musl.tar.gz"
        filename="qdrant-x86_64-linux.tar.gz"
        if download_file "Linux x86_64" "$url" "$filename"; then
            ((success_count++))
        else
            ((fail_count++))
            echo -e "${YELLOW}提示: 您可以手动下载文件:${NC}"
            echo -e "${BLUE}  下载地址: $url${NC}"
            echo -e "${BLUE}  保存位置: $CACHE_DIR/$filename${NC}"
        fi
        ;;
    2)
        # Linux ARM64
        url="https://github.com/qdrant/qdrant/releases/download/${QDRANT_VERSION}/qdrant-aarch64-unknown-linux-musl.tar.gz"
        filename="qdrant-aarch64-linux.tar.gz"
        if download_file "Linux ARM64" "$url" "$filename"; then
            ((success_count++))
        else
            ((fail_count++))
            echo -e "${YELLOW}提示: 您可以手动下载文件:${NC}"
            echo -e "${BLUE}  下载地址: $url${NC}"
            echo -e "${BLUE}  保存位置: $CACHE_DIR/$filename${NC}"
        fi
        ;;
    3)
        # 全部下载
        # Linux x86_64
        url="https://github.com/qdrant/qdrant/releases/download/${QDRANT_VERSION}/qdrant-x86_64-unknown-linux-musl.tar.gz"
        filename="qdrant-x86_64-linux.tar.gz"
        if download_file "Linux x86_64" "$url" "$filename"; then
            ((success_count++))
        else
            ((fail_count++))
        fi
        
        # Linux ARM64
        url="https://github.com/qdrant/qdrant/releases/download/${QDRANT_VERSION}/qdrant-aarch64-unknown-linux-musl.tar.gz"
        filename="qdrant-aarch64-linux.tar.gz"
        if download_file "Linux ARM64" "$url" "$filename"; then
            ((success_count++))
        else
            ((fail_count++))
        fi
        ;;
    *)
        echo -e "${RED}无效选项${NC}"
        exit 1
        ;;
esac

# 显示结果
echo
echo -e "${BLUE}=====================================${NC}"
if [ $fail_count -eq 0 ]; then
    echo -e "${GREEN}✅ 所有文件下载成功！${NC}"
    echo -e "${GREEN}文件已缓存到: $CACHE_DIR${NC}"
    echo
    echo -e "${YELLOW}现在可以运行部署命令:${NC}"
    echo -e "${BLUE}  make deploy${NC}"
    echo
    echo -e "${CYAN}即使服务器网络不好，也能快速完成 Qdrant 安装${NC}"
else
    echo -e "${YELLOW}⚠️  部分文件下载失败${NC}"
    echo -e "${GREEN}成功: $success_count 个${NC}"
    echo -e "${RED}失败: $fail_count 个${NC}"
    echo
    echo -e "${YELLOW}对于下载失败的文件，您可以:${NC}"
    echo -e "  1. 重新运行此脚本"
    echo -e "  2. 使用代理后重试"
    echo -e "  3. 手动下载到指定位置"
fi
echo -e "${BLUE}=====================================${NC}"

# 显示缓存目录内容
echo
echo -e "${CYAN}缓存目录内容:${NC}"
ls -lh "$CACHE_DIR" 2>/dev/null | grep -E "\.tar\.gz$" || echo "  (空)"
echo

exit $fail_count
# PixelPunk Scripts 脚本说明

本目录包含 PixelPunk 项目的所有自动化脚本，分为构建、部署、开发、Docker 和工具等类别。

## 📁 目录结构

```
scripts/
├── build/          # 构建相关脚本
├── deploy/         # 部署相关脚本
├── dev/            # 开发环境脚本
├── docker/         # Docker 相关脚本
├── tools/          # 通用工具脚本
├── vector/         # 向量数据库脚本
└── makefiles/      # Makefile 模块
```

---

## 🔨 build/ - 构建脚本

### `build-release.sh`
**用途**: 跨平台发布包构建（主要构建脚本）
**Make 命令**: `make release`
**功能**:
- 支持 5 个平台: Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)
- 自动下载并打包 Qdrant v1.11.3
- 构建前端和后端
- 生成 tar.gz 安装包
- 三层缓存机制加速构建

**示例**:
```bash
# 交互式构建所有平台
make release

# 指定版本号
RELEASE_VERSION=v1.2.0 ./scripts/build/build-release.sh

# 只构建特定平台（2 = Linux ARM64）
make release PLATFORMS=2
```

---

### `build-macos-local.sh`
**用途**: macOS 本地快速构建
**Make 命令**: `make build-macos` (如果配置)
**功能**:
- 仅在 macOS 上构建当前架构
- 跳过交叉编译
- 用于本地开发测试

---

### `sync-frontend.sh`
**用途**: 同步前端构建产物到后端 static 目录
**调用时机**: 在 `build-release.sh` 中自动调用
**功能**:
- 复制 `web/dist/` → `internal/static/dist/`
- 确保前端资源嵌入到 Go 二进制

---

### `build-base-images.sh`
**用途**: 构建 Docker 基础镜像
**功能**:
- 创建包含编译依赖的基础镜像
- 加速后续 Docker 构建

---

### `build-webp-offline.sh` & `build-webp-turbo.sh`
**用途**: WebP 图像库编译脚本
**功能**:
- 离线构建 libwebp
- 用于 Windows 交叉编译

---

## 🚀 deploy/ - 部署脚本

### `deploy.sh`
**用途**: 主部署脚本（服务器部署）
**Make 命令**: `make deploy`
**功能**:
- 上传二进制到服务器
- 重启服务
- 支持配置文件部署

---

### `quick-deploy.sh`
**用途**: 快速部署脚本
**Make 命令**: `make quick-deploy`
**功能**:
- 快速重新部署（无需完整构建）
- 适用于小改动快速上线

---

### `deploy-qdrant.sh`
**用途**: 独立部署 Qdrant 向量数据库
**功能**:
- 下载并启动 Qdrant
- 配置向量存储服务

---

## 💻 dev/ - 开发环境脚本

### `airrun.sh`
**用途**: Go 代码热重载（使用 Air）
**Make 命令**: `make dev`
**功能**:
- 监听代码变化
- 自动重新编译并重启
- 开发时实时预览

**示例**:
```bash
make dev
# 或直接运行
./scripts/dev/airrun.sh
```

---

### `web-pm.sh`
**用途**: 前端项目管理脚本
**Make 命令**: `make web-dev`
**功能**:
- 启动前端开发服务器
- 管理前端依赖
- 前端构建

**示例**:
```bash
# 启动前端开发服务器
./scripts/dev/web-pm.sh run dev

# 构建前端生产版本
./scripts/dev/web-pm.sh build
```

---

### `dev-with-vector.sh`
**用途**: 启动完整开发环境（含向量数据库）
**功能**:
- 同时启动 Qdrant
- 启动后端热重载
- 一站式开发环境

---

## 🐳 docker/ - Docker 脚本

### `build-push-multiarch.sh`
**用途**: 构建并推送多架构 Docker 镜像
**Make 命令**: `make docker-build`
**功能**:
- 支持 linux/amd64, linux/arm64, linux/arm/v7
- 使用 Docker Buildx
- 自动推送到 Docker Hub
- BuildKit 缓存加速

**示例**:
```bash
make docker-build
# 或指定版本
DOCKER_VERSION=v1.2.0 make docker-build
```

---

### `build-local-only.sh`
**用途**: 仅构建本地架构 Docker 镜像
**功能**:
- 快速本地测试
- 不推送到 Registry

---

### `docker-install.sh`
**用途**: Docker 环境安装脚本
**功能**:
- 安装 Docker 和 Docker Compose
- 配置 Docker 环境

---

## 🛠️ tools/ - 工具脚本

### `install.sh`
**用途**: PixelPunk 一键安装脚本
**功能**:
- 下载最新版本
- 自动解压到当前目录
- 配置启动脚本
- 用户最常用的安装方式

**使用**:
```bash
curl -fsSL https://your-domain/install.sh | bash
# 或
wget -qO- https://your-domain/install.sh | bash
```

---

### `cleanup-docker.sh`
**用途**: Docker 环境清理
**Make 命令**: `make docker-clean`
**功能**:
- 清理悬空镜像
- 清理构建缓存
- 释放磁盘空间

---

### `prepare-release-images.sh`
**用途**: 准备发布镜像
**功能**:
- 预处理镜像资源
- 准备发布文件

---

## 🔍 vector/ - 向量数据库脚本

### `qdrant.sh`
**用途**: Qdrant 管理脚本
**Make 命令**:
- `make qdrant-start` - 启动
- `make qdrant-stop` - 停止
- `make qdrant-restart` - 重启
- `make qdrant-status` - 查看状态

**功能**:
- 下载 Qdrant 二进制
- 启动/停止/重启 Qdrant
- 健康检查

**示例**:
```bash
# 启动 Qdrant
./scripts/vector/qdrant.sh start

# 查看状态
./scripts/vector/qdrant.sh status

# 停止 Qdrant
./scripts/vector/qdrant.sh stop
```

---

### `predownload_qdrant.sh`
**用途**: 预下载 Qdrant 所有平台版本
**功能**:
- 批量下载所有平台的 Qdrant
- 缓存到 `.cache/qdrant/`
- 加速后续构建

---

## 📋 makefiles/ - Makefile 模块

### `Makefile.release`
**用途**: 定义发布构建相关命令
**包含命令**:
- `make release` - 发布构建
- `make release-clean` - 清理发布文件
- `make release-list` - 列出已构建版本

---

### `Makefile.docker`
**用途**: 定义 Docker 相关命令
**包含命令**:
- `make docker-build` - 构建推送镜像
- `make docker-up` - 启动 Compose
- `make docker-down` - 停止 Compose
- `make docker-logs` - 查看日志

---

## 🎯 常用工作流程

### 开发流程
```bash
# 1. 启动开发环境
make dev              # 后端热重载
make web-dev          # 前端开发服务器（另一个终端）

# 2. 启动 Qdrant（如需测试向量功能）
make qdrant-start
```

---

### 发布流程
```bash
# 1. 构建发布包
make release
# 输出: build/release/pixelpunk-v1.0.0-{platform}.tar.gz

# 2. 构建 Docker 镜像
make docker-build
# 输出: snine98/pixelpunk:latest (多架构)
```

---

### 部署流程
```bash
# 安装包部署
curl -fsSL https://your-domain/install.sh | bash

# Docker 部署
docker run -d -p 9520:9520 snine98/pixelpunk:latest

# Docker Compose 部署
make docker-up
```

---

## 🔧 环境变量

### 构建相关
```bash
RELEASE_VERSION=v1.2.0    # 发布版本号
SKIP_FRONTEND=true        # 跳过前端构建
SKIP_QDRANT=true          # 跳过 Qdrant 下载
FORCE_REBUILD=true        # 强制重新构建
```

### Docker 相关
```bash
DOCKER_VERSION=v1.2.0     # Docker 镜像版本
DOCKER_USERNAME=snine98   # Docker Hub 用户名
```

---

## 📞 脚本维护

### 添加新脚本时的注意事项:
1. 使用 `#!/bin/bash` 开头
2. 添加脚本说明注释
3. 添加 `set -e` 防止错误传播
4. 使用相对路径或 `$PROJECT_ROOT`
5. 更新本 README

### 调试脚本
```bash
# 开启调试模式
bash -x ./scripts/build/build-release.sh

# 检查语法
bash -n ./scripts/build/build-release.sh
```

---

## 📄 许可证

这些脚本是 PixelPunk 项目的一部分，遵循项目许可证。

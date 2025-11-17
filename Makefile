# PixelPunk 智能资源管理项目 Makefile

.PHONY: help dev deploy deploy-qdrant dev-backend dev-frontend dev-qdrant qdrant-start qdrant-stop qdrant-restart qdrant-status

# 默认目标
.DEFAULT_GOAL := help

# 项目配置
APP_NAME := pixelpunk
TMP_DIR := tmp
WEB_DIR := web

# 环境变量
PKG_CONFIG_PATH := /opt/homebrew/lib/pkgconfig
export PKG_CONFIG_PATH

# 颜色输出
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
NC     := \033[0m

## 显示帮助信息
help:
	@echo "$(GREEN)═══════════════════════════════════════$(NC)"
	@echo "$(GREEN)  PixelPunk 开发与部署命令$(NC)"
	@echo "$(GREEN)═══════════════════════════════════════$(NC)"
	@echo ""
	@echo "$(YELLOW)🚀 开发命令:$(NC)"
	@echo "  make dev          - 启动后端热更新 (推荐)"
	@echo "  make dev-backend  - 仅启动后端热更新"
	@echo "  make dev-frontend - 仅启动前端热更新"
	@echo "  make dev-qdrant   - 启动 Qdrant 向量数据库"
	@echo ""
	@echo "$(YELLOW)📦 打包发布 (生成用户安装包):$(NC)"
	@echo "  make release                - 🎯 交互式打包 (推荐)"
	@echo "  make release-help           - 查看详细打包选项"
	@echo "  make release-list           - 查看已生成的安装包"
	@echo "  make release-clean          - 清理打包目录"
	@echo ""
	@echo "  $(GREEN)单平台打包:$(NC)"
	@echo "  make release-linux-amd64    - Linux x86_64 安装包"
	@echo "  make release-linux-arm64    - Linux ARM64 安装包"
	@echo "  make release-darwin-amd64   - macOS Intel 安装包"
	@echo "  make release-darwin-arm64   - macOS Apple Silicon 安装包"
	@echo "  make release-windows-amd64  - Windows x86_64 安装包"
	@echo ""
	@echo "$(YELLOW)🐳 Docker 部署:$(NC)"
	@echo "  make docker-build           - 构建并推送多架构镜像"
	@echo "  make docker-push            - 构建并推送多架构镜像"
	@echo "  make docker-up              - 启动服务（docker-compose）"
	@echo "  make docker-down            - 停止服务"
	@echo "  make docker-help            - Docker 完整命令列表"
	@echo ""
	@echo "$(YELLOW)🚀 服务器部署 (开发者使用):$(NC)"
	@echo "  make deploy                 - 统一部署管理系统"
	@echo "  make deploy-quick           - 快速部署 (使用已构建二进制)"
	@echo "  make deploy-qdrant          - 单独部署 Qdrant 数据库"
	@echo "  make clean-cache            - 清理 Docker 构建缓存"
	@echo ""
	@echo "$(YELLOW)🗄️  向量数据库管理:$(NC)"
	@echo "  make qdrant-start           - 启动 Qdrant"
	@echo "  make qdrant-stop            - 停止 Qdrant"
	@echo "  make qdrant-status          - 查看状态"
	@echo "  make qdrant-restart         - 重启 Qdrant"
	@echo ""
	@echo "$(YELLOW)💡 快速开始:$(NC)"
	@echo "  • 本地开发: make dev (后端) + make dev-frontend (前端，可选)"
	@echo "  • 打包发布: make release"
	@echo "  • 服务器部署: make deploy"
	@echo ""
	@echo "$(GREEN)📁 安装包输出目录: build/release/$(NC)"
	@echo ""

## 启动后端开发环境（热更新）
dev:
	@echo "$(GREEN)启动后端开发环境...$(NC)"
	@echo "$(YELLOW)仅启动后端热重载，前端请另开终端运行: make dev-frontend$(NC)"
	@echo "$(YELLOW)后端地址: http://localhost:9520$(NC)"
	@echo "$(YELLOW)前端地址: http://localhost:5173 (需要另外启动)$(NC)"
	@echo "$(YELLOW)使用 Ctrl+C 停止后端服务$(NC)"
	@echo ""
	@# 检测并清理后端端口
	@if lsof -ti:9520 >/dev/null 2>&1; then \
		echo "$(YELLOW)⚠️  检测到端口 9520 已被占用，正在清理旧进程...$(NC)"; \
		pkill -f "air" 2>/dev/null || true; \
		pkill -f "tmp/main" 2>/dev/null || true; \
		pkill -f "scripts/dev/airrun.sh" 2>/dev/null || true; \
		lsof -ti:9520 | xargs kill -9 2>/dev/null || true; \
		sleep 1; \
		echo "$(GREEN)✅ 旧进程已清理$(NC)"; \
	fi
	@$(MAKE) dev-backend


## 仅启动后端热更新
dev-backend:
	@echo "$(GREEN)启动后端热重载...$(NC)"
	@# 检测并清理后端端口
	@if lsof -ti:9520 >/dev/null 2>&1; then \
		echo "$(YELLOW)⚠️  检测到端口 9520 已被占用，正在清理旧进程...$(NC)"; \
		pkill -f "air" 2>/dev/null || true; \
		pkill -f "tmp/main" 2>/dev/null || true; \
		pkill -f "scripts/dev/airrun.sh" 2>/dev/null || true; \
		lsof -ti:9520 | xargs kill -9 2>/dev/null || true; \
		sleep 1; \
		echo "$(GREEN)✅ 旧进程已清理$(NC)"; \
	fi
	@./scripts/dev/airrun.sh

## 仅启动前端热更新
dev-frontend:
	@echo "$(GREEN)启动前端开发服务器...$(NC)"
	@# 检测并清理前端端口
	@if lsof -ti:5173 >/dev/null 2>&1; then \
		echo "$(YELLOW)⚠️  检测到端口 5173 已被占用，正在清理旧进程...$(NC)"; \
		pkill -f "vite" 2>/dev/null || true; \
		lsof -ti:5173 | xargs kill -9 2>/dev/null || true; \
		sleep 1; \
		echo "$(GREEN)✅ 旧进程已清理$(NC)"; \
	fi
	@./scripts/dev/web-pm.sh run dev

## 仅启动Qdrant数据库
dev-qdrant:
	@echo "$(GREEN)启动Qdrant数据库...$(NC)"
	@echo "$(YELLOW)Qdrant地址: http://localhost:6333$(NC)"
	./scripts/vector/qdrant.sh start

## Qdrant数据库管理
qdrant-start:
	@echo "$(GREEN)启动Qdrant数据库...$(NC)"
	@./scripts/vector/qdrant.sh start

qdrant-stop:
	@echo "$(GREEN)停止Qdrant数据库...$(NC)"
	@./scripts/vector/qdrant.sh stop

qdrant-restart:
	@echo "$(GREEN)重启Qdrant数据库...$(NC)"
	@./scripts/vector/qdrant.sh restart

qdrant-status:
	@echo "$(GREEN)检查Qdrant数据库状态...$(NC)"
	@./scripts/vector/qdrant.sh status

## 统一交互式部署管理系统
deploy:
	@echo "$(GREEN)启动 PixelPunk 统一部署管理系统...$(NC)"
	@echo "$(YELLOW)使用方向键选择，Enter确认，q退出$(NC)"
	@./scripts/deploy/deploy.sh

## 快速部署（使用已构建的 pixelpunk-linux 二进制）
deploy-quick:
	@echo "$(GREEN)快速部署 PixelPunk...$(NC)"
	@chmod +x ./scripts/deploy/quick-deploy.sh
	@./scripts/deploy/quick-deploy.sh

## 单独部署向量数据库（交互式，完整输出）
# 从 deploy/deploy-arrow.sh 读取远端默认配置，如有需要可通过环境变量 REMOTE_HOST/REMOTE_USER/REMOTE_PORT 覆盖
DEPLOY_SCRIPT := deploy/deploy-arrow.sh
REMOTE_HOST_DEF := $(shell sed -n 's/^REMOTE_HOST="\(.*\)"/\1/p' $(DEPLOY_SCRIPT) | head -1)
REMOTE_USER_DEF := $(shell sed -n 's/^REMOTE_USER="\(.*\)"/\1/p' $(DEPLOY_SCRIPT) | head -1)
REMOTE_PORT_DEF := $(shell sed -n 's/^REMOTE_PORT="\(.*\)"/\1/p' $(DEPLOY_SCRIPT) | head -1)

REMOTE_HOST_RESOLVED := $(if $(REMOTE_HOST),$(REMOTE_HOST),$(REMOTE_HOST_DEF))
REMOTE_USER_RESOLVED := $(if $(REMOTE_USER),$(REMOTE_USER),$(REMOTE_USER_DEF))
REMOTE_PORT_RESOLVED := $(if $(REMOTE_PORT),$(REMOTE_PORT),$(REMOTE_PORT_DEF))

deploy-qdrant:
	@echo "$(GREEN)部署向量数据库(Qdrant) - 单独流程$(NC)"
	@echo "$(YELLOW)目标服务器$(NC): $(REMOTE_USER_RESOLVED)@$(REMOTE_HOST_RESOLVED):$(REMOTE_PORT_RESOLVED)"
	@echo "$(YELLOW)说明$(NC): 回车可使用默认值，部署过程将完整输出"
	@read -p "安装目录 [默认 /pixelpunk-qdrant] : " INSTALL_DIR; \
	  if [ -z "$$INSTALL_DIR" ]; then INSTALL_DIR="/pixelpunk-qdrant"; fi; \
	  read -p "HTTP 端口 [默认 6333] : " HTTP_PORT; \
	  if [ -z "$$HTTP_PORT" ]; then HTTP_PORT=6333; fi; \
	  read -p "gRPC 端口 [默认 6334] : " GRPC_PORT; \
	  if [ -z "$$GRPC_PORT" ]; then GRPC_PORT=6334; fi; \
	  echo ""; \
	  echo "$(GREEN)开始部署 Qdrant 到 $(REMOTE_HOST_RESOLVED)$(NC)"; \
	  ./scripts/deploy/deploy-qdrant.sh \
	    --remote-host "$(REMOTE_HOST_RESOLVED)" \
	    --remote-user "$(REMOTE_USER_RESOLVED)" \
	    --ssh-port "$(REMOTE_PORT_RESOLVED)" \
	    --install-dir "$$INSTALL_DIR" \
	    --http-port "$$HTTP_PORT" \
	    --grpc-port "$$GRPC_PORT"

## 清理 Docker 构建缓存
clean-cache:
	@echo "$(YELLOW)清理 Docker 构建缓存...$(NC)"
	@docker builder prune -f
	@echo "$(GREEN)✓ 缓存清理完成$(NC)"

# 包含发布构建系统
-include scripts/makefiles/Makefile.release

# 包含 Docker 构建系统
-include scripts/makefiles/Makefile.docker

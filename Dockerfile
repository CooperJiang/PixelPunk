# PixelPunk Docker 镜像
# 多阶段构建，优化镜像大小

# ============================================
# 阶段 1: 构建前端
# ============================================
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web

# 复制前端依赖文件
COPY web/package*.json ./
COPY web/pnpm-lock.yaml ./

# 启用 corepack 并安装依赖（lockfileVersion 5.4 需要 pnpm v7）
# 使用缓存挂载加速 pnpm 下载
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    corepack enable && \
    corepack prepare pnpm@7.33.7 --activate && \
    pnpm install --frozen-lockfile

# 复制前端源代码
COPY web/ ./

# 构建前端
RUN pnpm run build

# ============================================
# 阶段 2: 构建后端
# ============================================
FROM golang:1.23-alpine AS backend-builder

# 安装构建依赖（包括 WebP 支持）
RUN apk add --no-cache \
    gcc \
    g++ \
    make \
    pkgconfig \
    libwebp-dev \
    musl-dev

WORKDIR /app

# 复制 Go 依赖文件
COPY go.mod go.sum ./

# 配置 Go 代理并下载依赖（使用缓存挂载加速）
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 复制后端源代码
COPY . .

# 复制前端构建产物到后端静态目录
RUN rm -rf internal/static/dist
COPY --from=frontend-builder /app/web/dist ./internal/static/dist

# 构建后端（启用 CGO 支持 WebP，使用缓存加速编译）
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 \
    GOOS=linux \
    go build \
    -ldflags="-w -s -X main.Version=docker" \
    -o pixelpunk \
    ./cmd/main.go

# ============================================
# 阶段 3: 运行时镜像
# ============================================
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    libwebp \
    libstdc++ \
    libgcc \
    && rm -rf /var/cache/apk/*

# 设置时区（默认上海）
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=backend-builder /app/pixelpunk .

# 创建必要的目录
RUN mkdir -p configs data logs uploads

# 复制配置文件模板
COPY configs/config.example.yaml ./configs/config.example.yaml

# 设置权限
RUN chmod +x pixelpunk

# 暴露端口
EXPOSE 9520

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9520/health || exit 1

# 启动应用
CMD ["./pixelpunk"]

# 开发文档

本文档帮助开发者快速搭建 PixelPunk 开发环境。

## 环境要求

- **Go**: 1.23.0 或更高版本
- **Node.js**: 18.0 或更高版本
- **pnpm**: 最新版本
- **Qdrant**: 向量数据库（项目自带）

## 快速开始

### 1. 配置文件

首先需要创建配置文件：

```bash
# 复制配置模板
cp configs/config.example.yaml configs/config.yaml
```

根据需要修改 `configs/config.yaml` 中的配置：

```yaml
app:
  port: 9520 # 后端服务端口
  mode: "debug" # 开发模式使用 debug，生产使用 release

database:
  type: "sqlite" # 数据库类型：sqlite 或 mysql
  path: "./data/app.db" # SQLite 数据库文件路径

redis:
  host: "127.0.0.1" # Redis 地址
  port: 6379 # Redis 端口
  password: "" # Redis 密码（如有）

vector:
  enabled: true # 是否启用向量搜索
  qdrant_url: "http://localhost:6333" # Qdrant 地址
```

### 2. 启动向量数据库

在启动项目前，需要先启动 Qdrant 向量数据库：

```bash
make qdrant-start
```

> 其他 Qdrant 管理命令：
>
> - `make qdrant-stop` - 停止数据库
> - `make qdrant-restart` - 重启数据库
> - `make qdrant-status` - 查看状态

### 3. 启动后端（方式一：热更新）

**推荐在开发时使用，代码修改自动重启：**

```bash
make dev
```

### 4. 启动后端（方式二：直接运行）

```bash
go run cmd/main.go
```

后端启动后访问地址：`http://localhost:9520`

### 5. 启动前端

前端项目在 `web/` 目录下，使用 pnpm 管理依赖。

```bash
# 进入前端目录
cd web

# 安装依赖（首次运行）
pnpm install

# 启动开发服务器
pnpm dev
```

或者在项目根目录使用 Make 命令：

```bash
make dev-frontend
```

前端启动后访问地址：`http://localhost:3800`

## 开发流程

### 推荐的启动顺序

1. 启动 Qdrant：`make qdrant-start`
2. 启动后端：`make dev`（另开一个终端）
3. 启动前端：`make dev-frontend`（再开一个终端）

### 前端开发

前端使用 Vue 3 + Vite 开发，主要目录结构：

```
web/
├── src/
│   ├── pages/          # 页面组件
│   ├── components/     # 通用组件
│   ├── composables/    # 组合式函数
│   ├── stores/         # 状态管理
│   └── utils/          # 工具函数
└── public/             # 静态资源
```

常用命令：

```bash
pnpm dev              # 启动开发服务器
pnpm build            # 打包生产版本
pnpm lint:check       # 检查代码规范
pnpm lint:fix         # 自动修复代码规范
```

### 后端开发

后端使用 Go + Gin 框架，主要目录结构：

```
.
├── cmd/                # 程序入口
├── internal/           # 内部代码
│   ├── controllers/    # 控制器
│   ├── services/       # 业务逻辑
│   ├── models/         # 数据模型
│   └── middleware/     # 中间件
└── pkg/                # 公共包
```

## 常见问题

### 端口被占用

如果端口被占用，Make 命令会自动清理旧进程。手动清理方式：

```bash
# 清理后端端口 9520
lsof -ti:9520 | xargs kill -9

# 清理前端端口 5173
lsof -ti:5173 | xargs kill -9

# 清理 Qdrant 端口 6333
lsof -ti:6333 | xargs kill -9
```

### 数据库初始化

首次运行时，后端会自动创建数据库表结构，无需手动初始化。

### 前端 API 代理

开发时前端会自动代理 API 请求到后端，配置在 `web/vite.config.ts` 中。

## 更多资源

- [主题开发文档](../web/docs/THEME-DEVELOPMENT-GUIDE.md) - 如何开发自定义前端主题
- [国际化文档](../web/docs/I18N-GUIDE.md) - 如何添加多语言支持
- [Makefile 命令](../Makefile) - 查看所有可用的 Make 命令

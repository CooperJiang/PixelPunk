# 部署文档

本文档介绍如何将 PixelPunk 部署到生产环境。

---

## 🌟 支持的平台

| 操作系统 | 架构 | 安装包名称 | 一键安装 | 状态 |
|---------|------|-----------|---------|------|
| Linux | x86_64 | `pixelpunk-v*-linux-amd64.tar.gz` | ✅ 支持 | ✅ 稳定 |
| Linux | ARM64 | `pixelpunk-v*-linux-arm64.tar.gz` | ✅ 支持 | ✅ 稳定 |
| macOS | Intel | `pixelpunk-v*-darwin-amd64.tar.gz` | ✅ 支持 | ✅ 稳定 |
| macOS | Apple Silicon | `pixelpunk-v*-darwin-arm64.tar.gz` | ✅ 支持 | ✅ 稳定 |
| Windows | x86_64 | `pixelpunk-v*-windows-amd64.zip` | ❌ 手动安装 | ✅ 支持 |

> **💡 说明**：Linux 和 macOS 支持一键安装脚本，Windows 需要手动下载安装包。

---

## 📦 部署方式

### 方式一：一键安装脚本（最简单，推荐）

适用于 **Linux** 和 **macOS** 系统，一行命令完成安装。

#### 快速安装

```bash
curl -fsSL http://download.pixelpunk.cc/shell/setup.sh | bash
```

或使用 wget：

```bash
wget -qO- http://download.pixelpunk.cc/shell/setup.sh | bash
```

#### 安装过程

脚本会自动执行以下操作：

1. **检测系统环境**
   - 自动识别操作系统（Linux/macOS）
   - 自动识别架构（x86_64/ARM64）
   - 检查必要依赖（curl/wget、tar）

2. **下载安装包**
   - 自动选择对应平台的安装包
   - 从 `http://download.pixelpunk.cc/release/` 下载
   - 当前版本：v1.0.0

3. **安装配置**
   - 解压到 `~/pixelpunk`（可自定义）
   - 创建必要目录结构
   - 生成初始配置文件

4. **启动服务**
   - 配置主程序端口（默认 9800）
   - 启动 PixelPunk 服务
   - 可选择设置开机自启动

#### 自定义安装目录

```bash
PIXELPUNK_INSTALL_DIR=/opt/pixelpunk curl -fsSL http://download.pixelpunk.cc/shell/setup.sh | bash
```

#### 安装完成

安装成功后，你会看到：

```
════════════════════════════════════════
  ✨ 安装完成！
════════════════════════════════════════

📌 下一步操作：

1. 打开浏览器访问系统配置页面：
   http://localhost:9800/setup

2. 在配置页面完成以下设置：
   • 数据库配置 (MySQL/SQLite)
   • Redis 配置 (可选)
   • 向量数据库配置 (可选，支持以图搜图)
   • 管理员账号创建
```

#### 服务管理

```bash
cd ~/pixelpunk  # 或你自定义的安装目录

./pixelpunk.sh start    # 启动服务
./pixelpunk.sh stop     # 停止服务
./pixelpunk.sh restart  # 重启服务
./pixelpunk.sh status   # 查看状态
./pixelpunk.sh logs     # 查看日志
```

---

### 方式二：手动下载安装（适用于所有平台）

适用于所有平台（Linux、macOS、Windows），适合无法使用一键脚本的情况。

#### 1. 下载安装包

访问下载地址，选择对应平台的安装包：

```
http://download.pixelpunk.cc/release/
```

或从 GitHub Releases 下载：

```
https://github.com/CooperJiang/PixelPunk-v1/releases
```

**平台选择**：
- Linux x86_64: `pixelpunk-v1.0.0-linux-amd64.tar.gz`
- Linux ARM64: `pixelpunk-v1.0.0-linux-arm64.tar.gz`
- macOS Intel: `pixelpunk-v1.0.0-darwin-amd64.tar.gz`
- macOS Apple Silicon: `pixelpunk-v1.0.0-darwin-arm64.tar.gz`
- Windows x86_64: `pixelpunk-v1.0.0-windows-amd64.zip`

#### 2. 解压安装包

**Linux/macOS**：
```bash
# 解压
tar -xzf pixelpunk-v1.0.0-linux-amd64.tar.gz

# 进入目录
cd pixelpunk-v1.0.0-linux-amd64
```

**Windows**：
```powershell
# 右键解压 zip 文件，或使用命令：
Expand-Archive pixelpunk-v1.0.0-windows-amd64.zip

# 进入目录
cd pixelpunk-v1.0.0-windows-amd64
```

#### 3. 安装包内容

```
pixelpunk-v1.0.0-linux-amd64/
├── pixelpunk               # 主程序（Linux/macOS）
├── pixelpunk.exe           # 主程序（Windows）
├── install.sh              # 安装脚本（Linux/macOS）
├── pixelpunk.sh            # 服务管理脚本
├── configs/
│   └── config.example.yaml # 配置模板
├── qdrant/                 # 内置向量数据库
│   ├── bin/qdrant
│   └── config/
└── README.txt              # 使用说明
```

#### 4. 运行安装

**Linux/macOS**：

```bash
# 运行安装脚本
./install.sh
```

安装脚本会自动：
- 配置主程序端口（默认 9800）
- 创建必要的目录结构
- 生成配置文件
- 启动 PixelPunk 服务

**Windows**：

```powershell
# 创建配置文件
copy configs\config.example.yaml configs\config.yaml

# 编辑配置文件（设置端口等）
notepad configs\config.yaml

# 启动主程序
.\pixelpunk.exe
```

#### 5. 服务管理

**Linux/macOS**：

```bash
./pixelpunk.sh start    # 启动服务
./pixelpunk.sh stop     # 停止服务
./pixelpunk.sh restart  # 重启服务
./pixelpunk.sh status   # 查看状态
./pixelpunk.sh logs     # 查看日志
```

**Windows**：

直接运行 `pixelpunk.exe` 或通过任务管理器管理进程。

#### 6. 访问系统

安装完成后，打开浏览器访问：

```
http://localhost:9800/setup
```

在配置页面完成：
- 数据库配置（MySQL/SQLite）
- Redis 配置（可选）
- 向量数据库配置（可选）
- 管理员账号创建

---

### 方式三：Docker 部署（容器化，推荐生产环境）

适用于容器化环境，快速部署，支持一键启动。

#### 3.1 使用 Docker Hub 镜像（最简单）

直接从 Docker Hub 拉取并运行：

```bash
# 拉取镜像
docker pull snine98/pixelpunk:latest

# 下载 docker-compose.yml
curl -O https://raw.githubusercontent.com/snine98/pixelpunk/main/docker-compose.yml

# 启动所有服务（PixelPunk + MySQL + Qdrant + Redis）
docker-compose up -d
```

服务地址：
- **PixelPunk 主应用**: `http://localhost:9520`
- **Qdrant 向量数据库**: `http://localhost:6333`
- **MySQL 数据库**: `localhost:3306`
- **Redis 缓存**: `localhost:6379`

#### 3.2 查看服务状态

```bash
# 查看运行状态
docker-compose ps

# 查看日志
docker-compose logs -f pixelpunk

# 停止服务
docker-compose down

# 重启服务
docker-compose restart
```

#### 3.3 自定义配置

如需自定义配置，编辑 `configs/config.docker.yaml`，然后重启服务：

```bash
# 创建配置目录
mkdir -p configs

# 复制配置模板
cp configs/config.example.yaml configs/config.docker.yaml

# 编辑配置（修改数据库密码、端口等）
vim configs/config.docker.yaml

# 重启服务以应用配置
docker-compose restart
```

#### 3.4 本地构建 Docker 镜像（开发者）

如果需要自行构建镜像：

```bash
# 方式1: 使用 Make 命令（推荐）
make docker-build

# 方式2: 使用脚本
./scripts/docker/build-docker.sh

# 方式3: 直接使用 Docker 命令
docker build -f Dockerfile -t snine98/pixelpunk:latest .
```

**推送到 Docker Hub**：

```bash
# 登录 Docker Hub
docker login

# 推送镜像
make docker-push

# 或使用脚本
./scripts/docker/push-docker.sh
```

#### 3.5 Docker 完整命令列表

```bash
# 构建相关
make docker-build           # 构建镜像（完整构建，包含前端）
make docker-build-quick     # 快速构建（使用已有前端）
make docker-push            # 推送到 Docker Hub

# 部署相关
make docker-up              # 启动所有服务
make docker-down            # 停止所有服务
make docker-restart         # 重启服务
make docker-ps              # 查看服务状态

# 日志相关
make docker-logs            # 查看所有服务日志
make docker-logs-app        # 仅查看应用日志

# 清理相关
make docker-clean           # 清理容器和缓存
make docker-clean-all       # 完全清理（包括镜像和数据）

# 帮助
make docker-help            # 查看完整命令列表
```

#### 3.6 Docker Compose 服务说明

项目提供的 `docker-compose.yml` 包含以下服务：

| 服务名 | 镜像 | 端口 | 说明 |
|--------|------|------|------|
| **pixelpunk** | `snine98/pixelpunk:latest` | 9520 | 主应用服务 |
| **mysql** | `mysql:8.0` | 3306 | MySQL 数据库 |
| **qdrant** | `qdrant/qdrant:latest` | 6333, 6334 | Qdrant 向量数据库 |
| **redis** | `redis:7-alpine` | 6379 | Redis 缓存（可选） |

**数据持久化**：
- MySQL 数据: `mysql-data` volume
- Qdrant 数据: `qdrant-data` volume
- Redis 数据: `redis-data` volume
- 上传文件: `./uploads` 目录映射

**生产环境建议**：
1. 修改 MySQL root 密码（`docker-compose.yml` 中的 `MYSQL_ROOT_PASSWORD`）
2. 修改应用数据库密码（`configs/config.docker.yaml`）
3. 配置持久化存储（使用外部卷或 NFS）
4. 配置反向代理（Nginx）并启用 HTTPS

---

### 方式四：开发者服务器部署

适用于直接部署到远程服务器（需要构建环境），仅供开发者使用。

```bash
# 使用 Make 命令交互式部署
make deploy
```

这个命令会：
1. 同步前端代码
2. 构建后端
3. 上传到服务器
4. 配置并启动服务

---

## 🔧 生产环境优化

### 1. 数据库优化

**使用 MySQL 而非 SQLite**：

```yaml
database:
  type: "mysql"
  host: "your-mysql-host"
  port: 3306
  username: "pixelpunk"
  password: "strong_password"
  name: "pixelpunk"
```

**创建数据库**：

```sql
CREATE DATABASE pixelpunk CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'pixelpunk'@'%' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON pixelpunk.* TO 'pixelpunk'@'%';
FLUSH PRIVILEGES;
```

### 2. 反向代理配置

**Nginx 配置示例**：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # 上传文件大小限制
    client_max_body_size 100M;

    location / {
        proxy_pass http://localhost:9520;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket 支持
    location /ws {
        proxy_pass http://localhost:9520;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

**Caddy 配置示例**：

```caddy
your-domain.com {
    reverse_proxy localhost:9520

    # 自动 HTTPS
    encode gzip

    # 上传大小限制
    request_body {
        max_size 100MB
    }
}
```

### 3. 系统服务配置

**Systemd 服务文件**（`/etc/systemd/system/pixelpunk.service`）：

```ini
[Unit]
Description=PixelPunk Image Hosting Service
After=network.target

[Service]
Type=simple
User=pixelpunk
Group=pixelpunk
WorkingDirectory=/opt/pixelpunk
ExecStart=/opt/pixelpunk/pixelpunk
Restart=on-failure
RestartSec=5s

# 环境变量
Environment="APP_MODE=release"

# 日志
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 重载配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start pixelpunk

# 开机自启
sudo systemctl enable pixelpunk

# 查看状态
sudo systemctl status pixelpunk

# 查看日志
sudo journalctl -u pixelpunk -f
```

### 4. 安全建议

- ✅ 使用强密码
- ✅ 配置防火墙，只开放必要端口
- ✅ 启用 HTTPS
- ✅ 定期备份数据库和上传文件
- ✅ 限制上传文件大小和类型
- ✅ 配置 Redis 密码
- ✅ 使用非 root 用户运行服务

---

## 🔍 健康检查

### 检查服务状态

```bash
# 检查进程
ps aux | grep pixelpunk

# 检查端口
netstat -tlnp | grep 9520

# 检查日志
tail -f /var/log/pixelpunk/app.log
```

### 数据库连接测试

```bash
# MySQL
mysql -h localhost -u pixelpunk -p pixelpunk

# Redis
redis-cli ping
```

### Qdrant 连接测试

```bash
curl http://localhost:6333/health
```

---

## 📊 监控与日志

### 日志位置

- **应用日志**: `/var/log/pixelpunk/app.log`（或配置的路径）
- **访问日志**: Nginx/Caddy 日志
- **系统日志**: `journalctl -u pixelpunk`

### 监控建议

- 使用 Prometheus + Grafana 监控系统资源
- 配置告警通知
- 定期检查磁盘空间
- 监控数据库性能

---

## 🔄 更新与维护

### 更新版本

```bash
# 1. 备份数据
cp -r /opt/pixelpunk/data /opt/pixelpunk/data.backup
mysqldump -u pixelpunk -p pixelpunk > backup.sql

# 2. 停止服务
./pixelpunk.sh stop

# 3. 替换二进制文件
mv pixelpunk pixelpunk.old
cp pixelpunk-new pixelpunk
chmod +x pixelpunk

# 4. 启动服务
./pixelpunk.sh start

# 5. 检查状态
./pixelpunk.sh status
```

### 数据备份

建议定期备份：

```bash
# 备份脚本示例
#!/bin/bash
BACKUP_DIR="/backup/pixelpunk/$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR

# 备份数据库
mysqldump -u pixelpunk -p pixelpunk > $BACKUP_DIR/database.sql

# 备份上传文件
tar -czf $BACKUP_DIR/uploads.tar.gz /opt/pixelpunk/data/uploads

# 备份配置
cp /opt/pixelpunk/configs/config.yaml $BACKUP_DIR/

# 删除 7 天前的备份
find /backup/pixelpunk -type d -mtime +7 -exec rm -rf {} \;
```

---

## ❓ 常见问题

### 端口被占用

```bash
# 查看占用端口的进程
lsof -i:9520

# 修改配置文件中的端口
vim configs/config.yaml
```

### 数据库连接失败

- 检查数据库服务是否运行
- 确认配置文件中的连接信息
- 测试网络连接

### 向量搜索不可用

- 检查 Qdrant 服务状态
- 确认 `vector.enabled` 配置为 `true`
- 查看 Qdrant 日志

---

## 📞 获取帮助

- 📖 查看 [开发文档](./DEVELOPMENT.md)
- 🐛 提交 [Issue](https://github.com/CooperJiang/PixelPunk-v1/issues)
- 💬 参与 [Discussions](https://github.com/CooperJiang/PixelPunk-v1/discussions)

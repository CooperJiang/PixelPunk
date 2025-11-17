# Docker 部署指南

本文档详细说明如何使用 Docker 部署 PixelPunk。

---

## 📦 镜像信息

- **Docker Hub**: `snine98/pixelpunk`
- **最新版本**: `snine98/pixelpunk:latest`
- **指定版本**: `snine98/pixelpunk:v1.0.0`

---

## 🚀 快速开始

### 使用 Docker Compose（推荐）

**1. 拉取镜像和配置文件**

```bash
# 拉取 Docker 镜像
docker pull snine98/pixelpunk:latest

# 下载 docker-compose.yml
curl -O https://raw.githubusercontent.com/snine98/pixelpunk/main/docker-compose.yml

# 下载配置文件
mkdir -p configs
curl -o configs/config.docker.yaml https://raw.githubusercontent.com/snine98/pixelpunk/main/configs/config.docker.yaml
```

**2. 启动服务**

```bash
docker-compose up -d
```

**3. 访问服务**

- PixelPunk 主应用: `http://localhost:9520`
- Qdrant 向量数据库: `http://localhost:6333`
- MySQL 数据库: `localhost:3306`

---

## 🔧 配置管理

### 环境变量

通过 `docker-compose.yml` 配置环境变量：

```yaml
services:
  pixelpunk:
    environment:
      - TZ=Asia/Shanghai           # 时区
      - APP_MODE=release            # 运行模式
```

### 配置文件

编辑 `configs/config.docker.yaml` 修改数据库连接、端口等配置：

```yaml
app:
  port: 9520
  mode: "release"

database:
  type: "mysql"
  host: "mysql"
  username: "pixelpunk"
  password: "pixelpunk_pass"  # ⚠️ 生产环境请修改
```

修改后重启服务：

```bash
docker-compose restart
```

---

## 📊 服务管理

### 查看状态

```bash
# 查看运行状态
docker-compose ps

# 查看详细信息
docker ps
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 仅查看 PixelPunk 日志
docker-compose logs -f pixelpunk

# 查看最近 100 行
docker-compose logs --tail=100 pixelpunk
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启单个服务
docker-compose restart pixelpunk

# 重新加载配置
docker-compose up -d --force-recreate
```

### 停止服务

```bash
# 停止服务（保留数据）
docker-compose stop

# 停止并删除容器（保留数据卷）
docker-compose down

# 完全清除（包括数据卷）
docker-compose down -v
```

---

## 💾 数据持久化

### 数据卷

Docker Compose 使用以下数据卷：

| 卷名 | 用途 | 说明 |
|------|------|------|
| `mysql-data` | MySQL 数据库 | 持久化数据库数据 |
| `qdrant-data` | Qdrant 向量库 | 持久化向量数据 |
| `redis-data` | Redis 缓存 | 持久化缓存数据 |
| `./uploads` | 上传文件 | 用户上传的图片 |

### 备份数据

```bash
# 备份 MySQL
docker exec pixelpunk-mysql mysqldump -u pixelpunk -ppixelpunk_pass pixelpunk > backup.sql

# 备份上传文件
tar -czf uploads-backup.tar.gz uploads/

# 备份数据卷
docker run --rm -v mysql-data:/data -v $(pwd):/backup alpine tar -czf /backup/mysql-backup.tar.gz /data
```

### 恢复数据

```bash
# 恢复 MySQL
docker exec -i pixelpunk-mysql mysql -u pixelpunk -ppixelpunk_pass pixelpunk < backup.sql

# 恢复上传文件
tar -xzf uploads-backup.tar.gz

# 恢复数据卷
docker run --rm -v mysql-data:/data -v $(pwd):/backup alpine tar -xzf /backup/mysql-backup.tar.gz -C /
```

---

## 🛠️ Make 命令（开发者）

如果在项目源码目录，可以使用 Make 命令简化操作：

### 构建命令

```bash
make docker-build           # 构建镜像（完整构建）
make docker-build-quick     # 快速构建（使用已有前端）
make docker-push            # 推送到 Docker Hub
make docker-login           # 登录 Docker Hub
```

### 部署命令

```bash
make docker-up              # 启动所有服务
make docker-down            # 停止所有服务
make docker-restart         # 重启服务
make docker-ps              # 查看服务状态
```

### 日志命令

```bash
make docker-logs            # 查看所有服务日志
make docker-logs-app        # 仅查看应用日志
```

### 清理命令

```bash
make docker-clean           # 清理容器和缓存
make docker-clean-all       # 完全清理（包括镜像和数据）
```

### 帮助命令

```bash
make docker-help            # 查看完整命令列表
```

---

## 🔒 安全建议

### 1. 修改默认密码

编辑 `docker-compose.yml`：

```yaml
services:
  mysql:
    environment:
      MYSQL_ROOT_PASSWORD: your_strong_root_password
      MYSQL_PASSWORD: your_strong_user_password
```

编辑 `configs/config.docker.yaml`：

```yaml
database:
  password: "your_strong_user_password"
```

### 2. 使用独立网络

```yaml
networks:
  pixelpunk-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### 3. 限制资源使用

```yaml
services:
  pixelpunk:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 512M
```

### 4. 使用只读文件系统

```yaml
services:
  pixelpunk:
    read_only: true
    tmpfs:
      - /tmp
      - /app/logs
```

---

## 🌐 反向代理配置

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name pixelpunk.example.com;

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

### Caddy 配置示例

```caddy
pixelpunk.example.com {
    reverse_proxy localhost:9520
}
```

---

## 🐛 常见问题

### 1. 端口已被占用

```bash
# 查看端口占用
lsof -i :9520

# 修改端口（编辑 docker-compose.yml）
ports:
  - "8080:9520"  # 使用 8080 端口
```

### 2. 容器无法启动

```bash
# 查看详细日志
docker-compose logs pixelpunk

# 查看容器状态
docker inspect pixelpunk
```

### 3. 数据库连接失败

```bash
# 检查 MySQL 是否就绪
docker-compose ps mysql

# 检查数据库日志
docker-compose logs mysql

# 测试连接
docker exec -it pixelpunk-mysql mysql -u pixelpunk -ppixelpunk_pass -e "SELECT 1"
```

### 4. 镜像拉取失败

```bash
# 使用镜像加速器
# 编辑 /etc/docker/daemon.json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com"
  ]
}

# 重启 Docker
sudo systemctl restart docker

# 重新拉取
docker pull snine98/pixelpunk:latest
```

---

## 📚 相关文档

- [完整部署文档](DEPLOYMENT.md)
- [开发文档](DEVELOPMENT.md)
- [配置说明](../configs/README.md)

---

## 🆘 获取帮助

- **问题反馈**: [GitHub Issues](https://github.com/CooperJiang/PixelPunk-v1/issues)
- **社区讨论**: [GitHub Discussions](https://github.com/CooperJiang/PixelPunk-v1/discussions)

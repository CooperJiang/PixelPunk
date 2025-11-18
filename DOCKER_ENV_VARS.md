# Docker 环境变量配置指南

PixelPunk 支持通过环境变量配置所有参数，适用于单容器部署场景。

## 配置优先级

```
默认配置 < 配置文件 < 环境变量 (最高优先级)
```

环境变量会覆盖配置文件中的相同配置项。

## 应用配置

| 环境变量 | 说明 | 默认值 | 示例 |
|---------|------|--------|------|
| `APP_APP_PORT` | 应用端口 | 9520 | 9520 |
| `APP_APP_MODE` | 运行模式 | release | release/debug |
| `APP_APP_NS` | 命名空间 | pixelpunk | pixelpunk |

## 数据库配置

| 环境变量 | 说明 | 必填 | 示例 |
|---------|------|------|------|
| `APP_DB_TYPE` | 数据库类型 | ✅ | mysql/sqlite |
| `APP_DB_HOST` | 数据库地址 | ✅ (MySQL) | 192.168.1.100 |
| `APP_DB_PORT` | 数据库端口 | ✅ (MySQL) | 3306 |
| `APP_DB_USERNAME` | 数据库用户名 | ✅ (MySQL) | pixelpunk |
| `APP_DB_PASSWORD` | 数据库密码 | ✅ (MySQL) | your_password |
| `APP_DB_NAME` | 数据库名称 | ✅ (MySQL) | pixelpunk |
| `APP_DB_PATH` | 数据库文件路径 | ✅ (SQLite) | /app/data/pixelpunk.db |

## Redis 配置

| 环境变量 | 说明 | 默认值 | 示例 |
|---------|------|--------|------|
| `APP_REDIS_HOST` | Redis 地址 | localhost | 192.168.1.101 |
| `APP_REDIS_PORT` | Redis 端口 | 6379 | 6379 |
| `APP_REDIS_PASSWORD` | Redis 密码 | (空) | your_redis_pass |
| `APP_REDIS_DB` | Redis 数据库 | 0 | 0 |

## 向量数据库配置

| 环境变量 | 说明 | 默认值 | 示例 |
|---------|------|--------|------|
| `APP_VECTOR_ENABLED` | 是否启用向量搜索 | true | true/false |
| `APP_VECTOR_QDRANT_URL` | Qdrant 地址 | - | http://192.168.1.102:6333 |
| `APP_VECTOR_TIMEOUT` | 请求超时时间(秒) | 30 | 30 |
| `APP_VECTOR_OPENAI_API_KEY` | OpenAI API Key | - | sk-xxx |
| `APP_VECTOR_OPENAI_BASE_URL` | OpenAI API 地址 | - | https://api.openai.com/v1 |
| `APP_VECTOR_OPENAI_MODEL` | 向量化模型 | - | text-embedding-ada-002 |

## 上传配置

| 环境变量 | 说明 | 默认值 | 示例 |
|---------|------|--------|------|
| `APP_UPLOAD_MAX_FILE_SIZE` | 最大文件大小(字节) | - | 104857600 |
| `APP_UPLOAD_ALLOWED_TYPES` | 允许的文件类型 | - | image/jpeg,image/png |

---

## 部署示例

### 示例 1: 连接外部 MySQL + Redis

```bash
docker run -d \
  --name pixelpunk \
  -p 9520:9520 \
  -v $(pwd)/uploads:/app/uploads \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  -e APP_DB_TYPE=mysql \
  -e APP_DB_HOST=192.168.1.100 \
  -e APP_DB_PORT=3306 \
  -e APP_DB_USERNAME=pixelpunk \
  -e APP_DB_PASSWORD=your_password \
  -e APP_DB_NAME=pixelpunk \
  -e APP_REDIS_HOST=192.168.1.101 \
  -e APP_REDIS_PORT=6379 \
  -e APP_VECTOR_ENABLED=false \
  snine98/pixelpunk:latest
```

### 示例 2: 使用 SQLite（无需外部数据库）

```bash
docker run -d \
  --name pixelpunk \
  -p 9520:9520 \
  -v $(pwd)/uploads:/app/uploads \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  -e APP_DB_TYPE=sqlite \
  -e APP_DB_PATH=/app/data/pixelpunk.db \
  -e APP_REDIS_HOST=192.168.1.101 \
  -e APP_REDIS_PORT=6379 \
  -e APP_VECTOR_ENABLED=false \
  snine98/pixelpunk:latest
```

### 示例 3: 完整配置（含向量搜索）

```bash
docker run -d \
  --name pixelpunk \
  -p 9520:9520 \
  -v $(pwd)/uploads:/app/uploads \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  -e APP_DB_TYPE=mysql \
  -e APP_DB_HOST=192.168.1.100 \
  -e APP_DB_PORT=3306 \
  -e APP_DB_USERNAME=pixelpunk \
  -e APP_DB_PASSWORD=your_password \
  -e APP_DB_NAME=pixelpunk \
  -e APP_REDIS_HOST=192.168.1.101 \
  -e APP_REDIS_PORT=6379 \
  -e APP_VECTOR_ENABLED=true \
  -e APP_VECTOR_QDRANT_URL=http://192.168.1.102:6333 \
  -e APP_VECTOR_OPENAI_API_KEY=sk-xxx \
  -e APP_VECTOR_OPENAI_BASE_URL=https://api.openai.com/v1 \
  snine98/pixelpunk:latest
```

---

## 与 Docker Compose 的区别

| 部署方式 | 配置方式 | 数据库/Redis | 适用场景 |
|---------|---------|-------------|----------|
| **Docker Compose** | 默认配置 | 自动创建容器 | 快速部署、测试环境 |
| **单容器 + 环境变量** | 环境变量 | 使用已有服务 | 生产环境、已有数据库 |
| **单容器 + 配置文件** | 挂载配置文件 | 使用已有服务 | 复杂配置场景 |

---

## 注意事项

1. **必需环境变量**：
   - MySQL: `APP_DB_TYPE`, `APP_DB_HOST`, `APP_DB_USERNAME`, `APP_DB_PASSWORD`, `APP_DB_NAME`
   - SQLite: `APP_DB_TYPE`, `APP_DB_PATH`

2. **数据持久化**：
   - 务必挂载 `/app/uploads`（用户上传的文件）
   - 务必挂载 `/app/data`（SQLite 数据库或其他数据）
   - 可选挂载 `/app/logs`（日志文件）

3. **网络访问**：
   - 容器内需要能访问外部数据库和 Redis
   - 如果使用 Docker 网络，可以使用容器名或服务名

4. **首次启动**：
   - 首次启动会自动创建数据库表
   - 建议先准备好数据库和 Redis 服务

---

## 故障排查

### 配置未生效？

检查环境变量名称是否正确：
```bash
docker exec pixelpunk env | grep APP_
```

### 数据库连接失败？

1. 检查数据库地址是否可达：
   ```bash
   docker exec pixelpunk ping -c 3 your_db_host
   ```

2. 检查数据库用户名密码是否正确

3. 查看日志：
   ```bash
   docker logs pixelpunk
   ```

### 如何查看当前配置？

```bash
# 查看环境变量
docker exec pixelpunk env | grep APP_

# 查看生成的配置文件
docker exec pixelpunk cat /app/configs/config.yaml
```

---

## 推荐部署方式

- **开发/测试环境**：使用 Docker Compose 一键部署
- **生产环境（已有数据库）**：使用单容器 + 环境变量
- **生产环境（复杂配置）**：使用单容器 + 挂载配置文件

更多信息请访问：https://docs.pixelpunk.cc

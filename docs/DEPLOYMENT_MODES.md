# PixelPunk 部署模式设计文档

## 📋 概述

PixelPunk 支持三种部署模式，以适应不同的使用场景和用户需求。系统通过环境变量自动识别部署模式，并智能调整Setup流程。

---

## 🎯 三种部署模式

### 1️⃣ Standalone 模式（安装包部署）

**适用场景**: 用户手动安装，自行配置所有服务

**特点**:
- 完整的Setup流程
- 用户需要配置数据库、向量数据库等所有服务
- 最灵活，支持各种自定义配置

**启动方式**:
```bash
# 下载安装包解压后
./pixelpunk

# 或从源码编译运行
go run ./cmd/main.go
```

**环境变量**:
```bash
DEPLOY_MODE=standalone  # 默认值，可省略
```

**Setup流程**: **Level 2（完整配置）**
- ✅ 站点基本信息
- ✅ 管理员账号
- ✅ 数据库配置（MySQL/SQLite）
- ✅ 向量数据库配置（Qdrant）
- ✅ Redis配置（可选）
- ✅ AI服务配置（可选）
- ✅ 邮件服务配置（可选）

---

### 2️⃣ Docker 模式（单容器部署）

**适用场景**: 快速体验、小规模部署、无Docker Compose环境

**特点**:
- 单个Docker容器即可运行
- 内嵌SQLite数据库（无需外部MySQL）
- 可选外部Qdrant（通过环境变量配置）
- 简化Setup流程

**构建镜像**:
```bash
# 构建All-in-One镜像
docker build -f Dockerfile.all-in-one -t pixelpunk:aio .
```

**启动方式**:
```bash
# 最简单启动（仅SQLite）
docker run -d \
  -p 9520:9520 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/uploads:/app/uploads \
  --name pixelpunk \
  pixelpunk:aio

# 带Qdrant配置启动
docker run -d \
  -p 9520:9520 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/uploads:/app/uploads \
  -e QDRANT_URL=http://your-qdrant:6333 \
  --name pixelpunk \
  pixelpunk:aio
```

**环境变量**:
```bash
DEPLOY_MODE=docker                          # 自动设置
QDRANT_URL=http://your-qdrant:6333         # 可选，向量数据库URL
REDIS_HOST=your-redis                       # 可选，Redis主机
REDIS_PORT=6379                             # 可选，Redis端口
```

**Setup流程**: **Level 1（基础配置）**
- ✅ 站点基本信息
- ✅ 管理员账号
- ❌ 数据库配置（已内嵌SQLite）
- ❌ 其他服务配置（通过环境变量）

---

### 3️⃣ Compose 模式（微服务部署）

**适用场景**: 生产环境、中大规模部署、需要服务隔离

**特点**:
- 服务分离：PixelPunk、MySQL、Qdrant、Redis各自独立容器
- 配置预设：config.docker.yaml预先配置好所有服务连接
- 最简Setup流程：仅需创建管理员账号
- 易于扩展和维护

**启动方式**:
```bash
# 一键启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

**环境变量**:
```bash
DEPLOY_MODE=compose      # 在docker-compose.yml中设置
CONFIG_PRESET=true       # 标识配置已预设
```

**Setup流程**: **Level 1（最简配置）**
- ✅ 站点基本信息
- ✅ 管理员账号
- ❌ 所有服务配置（已在config.docker.yaml预设）

---

## 🔧 技术实现

### 后端API

**GET /api/v1/setup/status** 返回扩展的安装状态：

```json
{
  "installed": false,
  "database_ok": false,
  "redis_ok": true,
  "message": "系统未安装，请先完成安装配置",
  "deploy_mode": "compose",     // standalone/docker/compose
  "setup_level": 1,             // 0=已安装 1=简化配置 2=完整配置
  "config_preset": true         // 配置是否预设
}
```

### 前端Setup页面适配

**建议实现逻辑** (Vue示例):

```typescript
// composables/useSetup.ts
export const useSetup = () => {
  const setupStatus = ref(null)

  const fetchStatus = async () => {
    const res = await axios.get('/api/v1/setup/status')
    setupStatus.value = res.data
  }

  const needsDatabaseConfig = computed(() => {
    return setupStatus.value?.setup_level === 2
  })

  const canSkipSetup = computed(() => {
    const status = setupStatus.value
    return status?.deploy_mode === 'compose' &&
           status?.config_preset === true &&
           status?.installed === true
  })

  return {
    setupStatus,
    fetchStatus,
    needsDatabaseConfig,
    canSkipSetup
  }
}
```

**Setup页面组件** (Vue):

```vue
<template>
  <div class="setup-page">
    <!-- Level 1: 基础配置（所有模式） -->
    <section class="basic-config">
      <h2>站点配置</h2>
      <el-input v-model="siteName" placeholder="站点名称" />

      <h2>管理员账号</h2>
      <el-input v-model="adminUsername" placeholder="用户名" />
      <el-input v-model="adminPassword" type="password" placeholder="密码" />
      <el-input v-model="adminEmail" placeholder="邮箱" />
    </section>

    <!-- Level 2: 完整配置（仅standalone模式） -->
    <section v-if="needsDatabaseConfig" class="advanced-config">
      <h2>数据库配置</h2>
      <el-select v-model="dbType">
        <el-option label="MySQL" value="mysql" />
        <el-option label="SQLite" value="sqlite" />
      </el-select>

      <template v-if="dbType === 'mysql'">
        <el-input v-model="dbHost" placeholder="主机" />
        <el-input v-model="dbPort" placeholder="端口" />
        <el-input v-model="dbUsername" placeholder="用户名" />
        <el-input v-model="dbPassword" type="password" placeholder="密码" />
        <el-input v-model="dbName" placeholder="数据库名" />
      </template>

      <h2>向量数据库配置</h2>
      <el-input v-model="qdrantUrl" placeholder="Qdrant URL" />

      <!-- 其他高级配置... -->
    </section>

    <!-- 部署模式提示 -->
    <div class="deploy-info">
      <el-tag v-if="deployMode === 'compose'" type="success">
        Docker Compose 部署 - 数据库已预配置
      </el-tag>
      <el-tag v-else-if="deployMode === 'docker'" type="info">
        Docker 单容器部署 - 使用内嵌数据库
      </el-tag>
      <el-tag v-else type="warning">
        安装包部署 - 需要完整配置
      </el-tag>
    </div>

    <el-button type="primary" @click="handleInstall">
      开始安装
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { useSetup } from '@/composables/useSetup'

const { setupStatus, fetchStatus, needsDatabaseConfig } = useSetup()

const deployMode = computed(() => setupStatus.value?.deploy_mode)

onMounted(() => {
  fetchStatus()
})

// ... 其他逻辑
</script>
```

---

## 📦 构建和使用

### 构建不同模式的镜像

```bash
# 1. 微服务模式镜像（用于docker-compose）
docker build -t snine98/pixelpunk:latest .

# 2. All-in-One镜像（单容器部署）
docker build -f Dockerfile.all-in-one -t snine98/pixelpunk:aio .
```

### 使用Docker Compose部署

```bash
# 1. 确保配置文件存在
ls configs/config.docker.yaml

# 2. 启动所有服务
docker-compose up -d

# 3. 访问 http://localhost:9520
# 首次访问进入Setup页面，仅需创建管理员账号
```

### 使用单容器Docker部署

```bash
# 1. 运行容器
docker run -d \
  -p 9520:9520 \
  -v pixelpunk-data:/app/data \
  -v pixelpunk-uploads:/app/uploads \
  --name pixelpunk \
  snine98/pixelpunk:aio

# 2. 访问 http://localhost:9520
# 首次访问进入Setup页面，仅需创建管理员账号
```

### 使用安装包部署

```bash
# 1. 解压安装包
tar -xzf pixelpunk-linux-amd64.tar.gz
cd pixelpunk

# 2. 启动应用
./pixelpunk

# 3. 访问 http://localhost:9520
# 首次访问进入Setup页面，需要完整配置数据库等服务
```

---

## 🔒 安全建议

### 生产环境Docker Compose部署

1. **修改默认密码**:
```yaml
# docker-compose.yml
mysql:
  environment:
    MYSQL_ROOT_PASSWORD: your-strong-password  # 修改此处
    MYSQL_PASSWORD: your-strong-password       # 修改此处
```

2. **限制CORS**:
```yaml
# configs/config.docker.yaml
cors:
  allow_origins:
    - "https://your-domain.com"  # 替换为实际域名
```

3. **使用环境变量文件**:
```bash
# 创建 .env 文件
echo "MYSQL_ROOT_PASSWORD=your-strong-password" > .env
echo "MYSQL_PASSWORD=your-strong-password" >> .env

# docker-compose.yml 中使用
mysql:
  environment:
    MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
```

4. **持久化数据**:
```yaml
volumes:
  - ./data:/app/data           # 数据库文件
  - ./uploads:/app/uploads     # 上传文件
  - ./logs:/app/logs           # 日志文件
```

---

## ❓ FAQ

**Q: 如何在Docker模式下连接外部Qdrant？**

A: 通过环境变量配置：
```bash
docker run -d \
  -e QDRANT_URL=http://your-qdrant-host:6333 \
  ...
```

**Q: Compose模式下如何自定义配置？**

A: 编辑 `configs/config.docker.yaml` 文件，重启服务生效：
```bash
docker-compose restart pixelpunk
```

**Q: 如何从Docker模式迁移到Compose模式？**

A:
1. 导出SQLite数据（如有）
2. 启动Compose服务
3. 在Setup时选择导入数据或全新安装

**Q: Setup完成后还能修改配置吗？**

A:
- Standalone: 直接修改 `configs/config.yaml`
- Docker/Compose: 通过管理员后台配置页面修改

---

## 📝 总结

| 模式 | 复杂度 | 适用场景 | Setup级别 |
|------|--------|----------|-----------|
| **Standalone** | ⭐⭐⭐ | 自定义部署 | Level 2（完整） |
| **Docker** | ⭐ | 快速体验 | Level 1（简化） |
| **Compose** | ⭐⭐ | 生产环境 | Level 1（最简） |

选择适合您的部署模式，开始使用PixelPunk！

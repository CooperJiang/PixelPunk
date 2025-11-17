# 发布与安装指南

本文档说明如何打包发布 PixelPunk 以及用户如何安装。

---

## 📦 打包发布流程

### 1. 打包新版本

运行打包命令，会提示输入版本号：

```bash
make release
```

**交互示例**：
```
═══════════════════════════════════════
  PixelPunk 发布构建
═══════════════════════════════════════

请输入版本号（默认: v1.0.0）:
版本号 [v1.0.0]: 1.2.0    # 输入版本号，支持 1.2.0 或 v1.2.0 格式

构建版本: v1.2.0

[开始构建...]
```

**支持的版本号格式**：
- `1.0.0` → 自动转换为 `v1.0.0`
- `v1.0.0` → 使用原格式
- `1.0.0-beta` → 转换为 `v1.0.0-beta`（支持预发布版本）

### 2. 打包结果

打包完成后，文件位于 `build/release/` 目录：

```
build/release/
├── pixelpunk-v1.2.0-linux-amd64.tar.gz
├── pixelpunk-v1.2.0-linux-arm64.tar.gz
├── pixelpunk-v1.2.0-darwin-amd64.tar.gz
├── pixelpunk-v1.2.0-darwin-arm64.tar.gz
└── pixelpunk-v1.2.0-windows-amd64.zip
```

### 3. 单独打包特定平台

如果只需要打包某个平台：

```bash
make release-linux-amd64     # Linux x86_64
make release-linux-arm64     # Linux ARM64
make release-darwin-amd64    # macOS Intel
make release-darwin-arm64    # macOS Apple Silicon
make release-windows-amd64   # Windows x86_64
```

每个命令都会提示输入版本号。

### 4. 上传到服务器

将打包好的文件上传到：

```
http://download.pixelpunk.cc/release/
```

上传所有生成的 `.tar.gz` 和 `.zip` 文件。

**示例**：
```bash
# 使用 scp 上传
scp build/release/pixelpunk-v1.2.0-*.tar.gz user@download.pixelpunk.cc:/path/to/release/
scp build/release/pixelpunk-v1.2.0-*.zip user@download.pixelpunk.cc:/path/to/release/

# 或使用 rsync
rsync -avz build/release/ user@download.pixelpunk.cc:/path/to/release/
```

### 5. 上传安装脚本

将 `scripts/tools/setup.sh` 上传到：

```
http://download.pixelpunk.cc/shell/setup.sh
```

**示例**：
```bash
scp scripts/tools/setup.sh user@download.pixelpunk.cc:/path/to/shell/setup.sh
```

---

## 🚀 用户安装方式

### 方式一：一键安装（推荐）

用户可以通过一行命令直接安装：

```bash
curl -fsSL http://download.pixelpunk.cc/shell/setup.sh | bash
```

或使用 wget：

```bash
wget -qO- http://download.pixelpunk.cc/shell/setup.sh | bash
```

**安装流程**：
1. 自动检测操作系统和架构
2. 下载 v1.0.0 版本安装包（当前默认版本）
3. 解压到 `~/pixelpunk`
4. 运行安装脚本配置服务
5. 启动 PixelPunk

> 💡 **注意**：当前只安装 v1.0.0 版本，后续有多版本支持后会添加版本选择功能。

**自定义安装目录**：
```bash
PIXELPUNK_INSTALL_DIR=/opt/pixelpunk curl -fsSL http://download.pixelpunk.cc/shell/setup.sh | bash
```

### 方式二：手动安装

#### 1. 下载安装包

访问下载地址，选择对应平台的安装包：

```
http://download.pixelpunk.cc/release/
```

**平台对应关系**：
- Linux x86_64: `pixelpunk-v{版本}-linux-amd64.tar.gz`
- Linux ARM64: `pixelpunk-v{版本}-linux-arm64.tar.gz`
- macOS Intel: `pixelpunk-v{版本}-darwin-amd64.tar.gz`
- macOS Apple Silicon: `pixelpunk-v{版本}-darwin-arm64.tar.gz`
- Windows: `pixelpunk-v{版本}-windows-amd64.zip`

#### 2. 解压安装包

**Linux/macOS**：
```bash
tar -xzf pixelpunk-v1.0.0-linux-amd64.tar.gz
cd pixelpunk-v1.0.0-linux-amd64
```

**Windows**：
```powershell
# 解压 zip 文件
cd pixelpunk-v1.0.0-windows-amd64
```

#### 3. 运行安装脚本

**Linux/macOS**：
```bash
./install.sh
```

**Windows**：
```powershell
# 手动启动 pixelpunk.exe
.\pixelpunk.exe
```

---

## 🔧 安装包内容

每个安装包包含以下文件：

```
pixelpunk-v1.0.0-linux-amd64/
├── pixelpunk              # 主程序（Linux/macOS）
├── pixelpunk.exe          # 主程序（Windows）
├── install.sh             # 安装脚本（自动配置）
├── pixelpunk.sh           # 服务管理脚本
├── configs/               # 配置文件目录
│   └── config.example.yaml
├── qdrant/                # 内置向量数据库
│   ├── bin/qdrant
│   └── config/
└── README.txt             # 使用说明
```

---

## 📋 服务管理

安装完成后，可以使用以下命令管理服务：

```bash
./pixelpunk.sh start    # 启动服务
./pixelpunk.sh stop     # 停止服务
./pixelpunk.sh restart  # 重启服务
./pixelpunk.sh status   # 查看状态
./pixelpunk.sh logs     # 查看日志
```

---

## 🌍 支持的平台

| 操作系统 | 架构 | 安装包名称 | 状态 |
|---------|------|-----------|------|
| Linux | x86_64 | `linux-amd64.tar.gz` | ✅ 支持 |
| Linux | ARM64 | `linux-arm64.tar.gz` | ✅ 支持 |
| macOS | Intel | `darwin-amd64.tar.gz` | ✅ 支持 |
| macOS | Apple Silicon | `darwin-arm64.tar.gz` | ✅ 支持 |
| Windows | x86_64 | `windows-amd64.zip` | ✅ 支持 |

---

## ⚙️ 环境变量

### 打包时

```bash
# 指定版本号（跳过交互输入）
RELEASE_VERSION=v1.2.0 make release

# 跳过前端构建（调试用）
SKIP_FRONTEND=true make release

# 跳过 Qdrant 下载（调试用）
SKIP_QDRANT=true make release
```

### 安装时

```bash
# 自定义安装目录
export PIXELPUNK_INSTALL_DIR=/opt/pixelpunk
curl -fsSL http://download.pixelpunk.cc/shell/setup.sh | bash
```

---

## 🐛 常见问题

### Q1: 打包时版本号格式错误

**错误**：
```
✗ 版本号格式错误！请使用格式：v1.0.0 或 1.0.0
```

**解决**：
- 使用正确格式：`1.0.0`、`v1.0.0`、`1.0.0-beta` 等
- 不要包含空格或特殊字符

### Q2: 下载安装包失败

**错误**：
```
✗ 版本不存在或无法访问
```

**解决**：
1. 检查版本号是否正确
2. 确认文件已上传到服务器
3. 访问 `http://download.pixelpunk.cc/release/` 查看可用版本

### Q3: 安装脚本权限不足

**错误**：
```
Permission denied: ./install.sh
```

**解决**：
```bash
chmod +x install.sh
./install.sh
```

### Q4: 端口被占用

**错误**：
```
✗ 端口 9800 已被占用
```

**解决**：
- 安装时选择其他端口
- 或停止占用该端口的程序：
  ```bash
  lsof -i :9800
  kill -9 <PID>
  ```

---

## 📞 获取帮助

- 📖 [完整文档](https://github.com/CooperJiang/PixelPunk-v1/tree/main/docs)
- 🐛 [报告问题](https://github.com/CooperJiang/PixelPunk-v1/issues)
- 💬 [讨论交流](https://github.com/CooperJiang/PixelPunk-v1/discussions)

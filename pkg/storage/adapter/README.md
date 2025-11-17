# PixelPunk 存储适配器架构文档（新版）

## 📋 概述

PixelPunk 图床系统采用插件化的存储适配器架构，支持多种存储后端，包括本地存储和各种云存储服务。每个存储适配器都实现统一的 `StorageAdapter` 接口，确保不同存储后端的无缝切换和一致性体验。

## 🏗️ 架构设计

### 设计原则
- **统一接口**：所有存储适配器实现相同的接口规范
- **插件化**：支持动态加载和切换不同的存储后端
- **可扩展性**：新存储后端可以轻松接入
- **配置驱动**：通过配置文件管理不同存储的参数
- **错误统一**：标准化的错误处理和类型

### 核心组件
```
pkg/storage/adapter/
├── adapter.go          # 接口定义和公共类型
├── local.go           # 本地存储适配器
├── cos.go             # 腾讯云COS适配器  
├── oss.go             # 阿里云OSS适配器
├── rainyun.go         # 雨云S3适配器
└── README.md          # 本文档
```

更多相关目录：
```
pkg/storage/
├── manager/           # 统一存储管理器（Adapter 缓存、调度、Base64 统一编码）
├── factory/           # 适配器工厂与注册
├── path/              # 路径工具：CleanObjectPath / EnsureObjectKey / ExtractObjectPathFromURL
├── url/               # URL 策略（代理/直链/CDN/HTTPS）
├── compat.go          # 兼容层（历史接口的过渡，逐步收敛）
└── storage.go         # 对外统一门面
```

## 📚 业务场景分析

### 1. 多样化存储需求

#### 场景描述
不同用户和部署环境对存储有不同需求：
- **个人部署**：需要简单的本地存储，无额外成本
- **中小企业**：需要云存储的扩展性，但成本敏感
- **大型企业**：需要企业级云存储，重视稳定性和性能
- **多地部署**：需要就近存储，减少访问延迟

#### 业务价值
- **降低成本**：根据规模选择合适的存储方案
- **提升性能**：就近存储减少网络延迟
- **增强稳定性**：多存储后端提供容灾备份
- **简化运维**：统一接口降低维护复杂度

### 2. 图片处理流水线

#### 上传流程
```
用户上传 → 图片处理 → 生成缩略图 → 存储原图和缩略图 → 返回访问URL
```

#### 关键特性
- **自动缩略图**：支持多种尺寸和格式的缩略图生成
- **格式转换**：支持WebP等现代格式，减少带宽消耗
- **质量压缩**：可配置的图片质量，平衡文件大小和视觉效果
- **水印处理**：支持在上传时自动添加水印

### 3. 访问控制和URL管理

#### URL类型
1. **代理URL（ID-based）**：`/f/{fileID}/[displayName]`（原图），`/t/{fileID}/[displayName]`（缩略图）
   - 通过服务器代理访问
   - 支持权限控制和访问日志
   - 适用于本地存储、私有内容或需要严格控制的场景

2. **直链URL**：`https://bucket.endpoint.com/images/<object-key>`
   - 直接访问云存储
   - 减少服务器负载
   - 适用于公开内容和高并发场景

#### 配置驱动的访问策略
```go
// 配置示例
{
    "hide_remote_url": true,     // 隐藏远程URL，强制使用代理
    "access_control": "private", // 访问控制类型
    "cdn_domain": "cdn.example.com", // CDN域名
    "force_https": true          // 强制HTTPS
}
```

### 4. 多租户隔离与游客策略

#### 对象键规范（云端统一）
- 登录用户：别名分片目录（每个用户有稳定别名 alias，并带两位 shard）
  - 原图：`images/<shard>/<alias>/(folder)/file`
  - 缩略图：`thumbnails/<shard>/<alias>/(folder)/file`
- 游客（userID=0）：固定目录（简化）
  - 原图：`images/guest/file`
  - 缩略图：`thumbnails/guest/file`

Local 适配器会将 `images/` 映射到 `uploads/images/`，将 `thumbnails/` 映射到 `uploads/thumbnails/`。

对象键构造不要手写拼接，优先使用：
- `tenant.BuildObjectKey(userID, folder, file)`
- `tenant.BuildThumbObjectKey(userID, folder, file)`
- `path.EnsureObjectKey(userID, input, isThumb)`（当输入为逻辑路径或完整URL时）

## 🔧 StorageAdapter 接口详解（新版）

### 核心方法

#### 1. 基本操作
```go
// 上传文件（核心功能）
Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error)

// 删除文件
Delete(ctx context.Context, path string) error

// 检查文件是否存在
Exists(ctx context.Context, path string) (bool, error)
```

#### 2. URL生成（新版）
```go
// 获取基础访问URL（直链场景下用于构建 RemoteURL，完整URL由上层策略构建）
GetURL(path string, options *URLOptions) (string, error)
```
> 新版不再在适配器层提供 `GetFullURL/GetCDNURL`；完整 URL 由上层 URL 策略或 ID-based 路由（/i,/t）构建。

#### 3. 文件读取
```go
// 读取文件内容（用于代理模式）
ReadFile(ctx context.Context, path string) (io.ReadCloser, error)
```
> 说明：Base64 编码已在 Manager 层统一实现（通过 ReadFile 读取后编码）。适配器无需再实现 Base64 方法。

#### 4. 权限控制
```go
// 设置对象访问权限
SetObjectACL(ctx context.Context, path string, acl string) error
```

#### 5. 配置和管理
```go
// 初始化适配器
Initialize(config map[string]interface{}) error

// 健康检查
HealthCheck(ctx context.Context) error

// 获取适配器类型
GetType() string

// 获取适配器能力
GetCapabilities() Capabilities
```

### 数据结构

#### UploadRequest 上传请求
```go
type UploadRequest struct {
    File         *multipart.FileHeader // HTTP上传的文件
    ProcessedData []byte               // 预处理数据（优先级更高）
    UserID       uint                  // 用户ID（用于路径隔离）
    FolderPath   string                // 文件夹路径
    FileName     string                // 文件名
    ContentType  string                // MIME类型
    Options      *UploadOptions        // 上传选项
}
```

#### UploadResult 上传结果
```go
type UploadResult struct {
    OriginalPath   string // 原始文件存储路径
    ThumbnailPath  string // 缩略图存储路径
    URL            string // 逻辑访问路径（存储到数据库）
    ThumbnailURL   string // 缩略图逻辑路径
    FullURL        string // 完整访问URL
    FullThumbURL   string // 完整缩略图URL
    RemoteURL      string // 云存储直链URL
    RemoteThumbURL string // 云存储缩略图直链
    Size           int64  // 文件大小
    Width          int    // 图片宽度
    Height         int    // 图片高度
    Hash           string // 文件MD5哈希
    ContentType    string // 内容类型
    Format         string // 图片格式
}
```

#### Capabilities 适配器能力
```go
type Capabilities struct {
    SupportsSignedURL bool     // 是否支持签名URL
    SupportsCDN       bool     // 是否支持CDN
    SupportsResize    bool     // 是否支持在线缩放
    SupportsWebP      bool     // 是否支持WebP
    MaxFileSize       int64    // 最大文件大小限制
    SupportedFormats  []string // 支持的文件格式
}
```

## 📊 现有适配器对比

| 适配器 | 类型 | 特点 | 适用场景 | CDN支持 | 签名URL |
|--------|------|------|----------|---------|---------|
| **Local** | 本地存储 | 简单、无成本 | 个人部署、开发测试 | ❌ | ❌ |
| **COS** | 腾讯云 | 企业级、稳定 | 中大型企业 | ✅ | ✅ |
| **OSS** | 阿里云 | 功能丰富 | 国内企业 | ✅ | ✅ |
| **RainyUN** | 雨云S3 | 成本优化 | 中小企业 | ✅ | ✅ |

### S3 渠道使用指引（统一）

- 一个“通用 S3”渠道可对接绝大多数对象存储：
  - 海外/通用：AWS S3、Cloudflare R2、Wasabi、Backblaze B2、DigitalOcean Spaces、MinIO、Ceph RGW、Linode/Akamai、Vultr、Scaleway、OVH、Oracle、IBM 等
  - 国内/私有化：华为 OBS、百度 BOS、金山 KS3、火山 TOS、UCloud US3、青云 QingStor、网易 NOS、京东云等
- 路径样式（use_path_style）选择：
  - 虚拟主机样式（否）：`bucket.host/key` —— 适用于 AWS 官方域名/支持 bucket 子域名的厂商
  - 路径样式（是）：`host/bucket/key` —— 适用于 MinIO/Ceph RGW/Cloudflare R2/自建或反代 Endpoint、IP/含端口 Endpoint、无泛域名证书
  - 建议：连接 AWS 官方 → 选“否”；填写了自定义 Endpoint/反代/IP/端口 → 选“是”
- 常见故障提示：
  - SignatureDoesNotMatch：校验 AccessKey/SecretKey、path-style 选择、Endpoint 与 Region 是否匹配
  - PermanentRedirect/301/Region mismatch：Region 与存储桶实际区域不一致，或使用了错误的 Endpoint
  - InvalidBucketName：存储桶命名不符合规范（小写字母/数字/短横线）
  - RequestTimeTooSkewed/Clock Skew：本机时间与服务端偏差过大，请同步时间

### Upyun/又拍云（补充）

- REST 签名：`UPYUN operator:Base64(HMAC-SHA1(md5(password hex), Method&URI&Date[&Content-MD5]))`
- 中文/空格路径：使用逐段编码 + RawPath，签名与请求使用一致的“已编码路径”
- 测试上传后删除：避免“立刻 DELETE”引发 429（concurrent put or delete），我们已在连接测试中做延迟/异步清理

### WebDAV（补充）

- 直链与代理：大部分 WebDAV 需要认证，建议默认走系统代理（allow_direct=false）；如需直链，需要服务端对 GET 公开或配置自定义域
- 可通过 Alist/aliyundrive-webdav 将各网盘（如 115/阿里云盘 等）映射为 WebDAV，再接入我们的 WebDAV 渠道

### 路径处理对比（更新）

| 适配器 | 本地物理路径 | 对象键（云端统一） | 说明 |
|--------|--------------|--------------------|------|
| Local | `uploads/images/(...)` | `images/(...)` | `images/` → `uploads/images/`；`thumbnails/` 同理 |
| COS/OSS/S3 | N/A | `images/<shard>/<alias>/(folder)/file` | 游客固定 `images/guest/file` |

## 🚀 新存储适配器接入指南

### 步骤1：创建适配器文件

创建 `pkg/storage/adapter/your_storage.go`：

```go
package adapter

import (
    "context"
    "io"
    // 导入所需的SDK包
)

// YourStorageAdapter 你的存储适配器
type YourStorageAdapter struct {
    config     Config
    client     *YourSDKClient // SDK客户端
    bucket     string
    region     string
    endpoint   string
    // 其他必要字段
}

// 确保实现了StorageAdapter接口
var _ StorageAdapter = (*YourStorageAdapter)(nil)
```

### 开发检查清单（推荐标准）

- 配置与初始化
  - 在适配器内部定义强类型 `Config`，从 `config.NewMapConfig` 转换并校验必要字段。
  - `Initialize` 中仅做解析和轻量校验，不做真实 I/O。
- 上传流程
  - 统一从 `UploadRequest` 读取源数据：若 `ProcessedData` 不为空则优先使用，否则从 `File.Open()` 读取。
  - 压缩与缩略图：优先使用 `pkg/storage/pipeline.GenerateOrFallback` 实现一致化缩略图生成；参数（质量/尺寸）用请求值或安全默认值。
  - 对象路径：不要手写。使用 `tenant.BuildObjectKey/BuildThumbObjectKey` 或 `path.EnsureObjectKey` 统一生成。
- URL 生成
  - 适配器实现 `GetURL/GetFullURL/GetCDNURL` 仅做基础拼接；是否隐藏远程URL、是否强制HTTPS、是否使用CDN由上层 `pkg/storage/url/strategy.go` 决策。
- 访问控制
  - `SetObjectACL` 映射到云厂商 ACL，支持 `public-read/private` 至少两类；本地可直接返回成功。
- 健康检查
  - `HealthCheck` 做最小可用性检查（如列目录/Head 对象）。
- 日志
  - 使用 `pkg/logger` 的 `logger.Debug/Info/Warn/Error`，避免使用 `fmt.Printf`。
- 错误
  - 统一返回 `adapter.NewStorageError`，按 `not_found/permission/quota/network/internal` 分类。

以上规范可确保新渠道与现有渠道在行为与指标上保持一致，同时降低维护复杂度。

### 步骤2：实现必需方法

#### 2.1 基础配置方法
```go
// GetType 返回存储类型标识符
func (a *YourStorageAdapter) GetType() string {
    return "your_storage" // 唯一标识符
}

// Initialize 初始化适配器
func (a *YourStorageAdapter) Initialize(configData map[string]interface{}) error {
    cfg := config.NewMapConfig(configData)
    
    // 读取配置
    a.bucket = cfg.GetString("bucket")
    a.region = cfg.GetString("region")
    a.endpoint = cfg.GetString("endpoint")
    accessKey := cfg.GetString("access_key")
    secretKey := cfg.GetString("secret_key")
    
    // 验证必需配置
    if a.bucket == "" || accessKey == "" || secretKey == "" {
        return errors.New("missing required configuration")
    }
    
    // 初始化SDK客户端
    client, err := your_sdk.NewClient(&your_sdk.Config{
        AccessKey: accessKey,
        SecretKey: secretKey,
        Region:    a.region,
        Endpoint:  a.endpoint,
    })
    if err != nil {
        return fmt.Errorf("failed to create client: %w", err)
    }
    
    a.client = client
    a.config = cfg
    return nil
}

// GetCapabilities 返回适配器能力
func (a *YourStorageAdapter) GetCapabilities() Capabilities {
    return Capabilities{
        SupportsSignedURL: true,  // 根据实际情况
        SupportsCDN:       true,  // 根据实际情况
        SupportsResize:    false, // 根据实际情况
        SupportsWebP:      false, // 根据实际情况
        MaxFileSize:       100 * 1024 * 1024, // 100MB
        SupportedFormats:  []string{"jpg", "jpeg", "png", "gif", "webp"},
    }
}

// HealthCheck 健康检查
func (a *YourStorageAdapter) HealthCheck(ctx context.Context) error {
    // 执行简单的API调用验证连接
    _, err := a.client.HeadBucket(ctx, a.bucket)
    return err
}
```

#### 2.2 核心上传功能
```go
// Upload 上传文件
func (a *YourStorageAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
    // 1. 验证请求
    if err := a.validateFile(req); err != nil {
        return nil, err
    }
    
    // 2. 读取和预处理图片数据
    var src io.Reader
    if len(req.ProcessedData) > 0 {
        src = bytes.NewReader(req.ProcessedData)
    } else {
        file, err := req.File.Open()
        if err != nil {
            return nil, err
        }
        defer file.Close()
        src = file
    }
    
    // 3. 图片处理（压缩、格式转换等）
    processedSrc, finalContentType, width, height, err := a.processImage(src, req)
    if err != nil {
        return nil, err
    }
    
    // 4. 生成文件路径
    originalFileName := req.FileName
    objectPath := a.buildObjectPath(req.UserID, req.FolderPath, originalFileName)
    logicalPath := utils.BuildLogicalPath(req.FolderPath, originalFileName)
    
    // 5. 上传原图
    originalResult, err := a.uploadToStorage(processedSrc, objectPath, finalContentType)
    if err != nil {
        return nil, err
    }
    
    // 6. 生成缩略图（如果需要）
    var thumbResult *UploadResult
    if req.Options != nil && req.Options.GenerateThumb {
        thumbResult, err = a.generateThumbnail(processedSrc, req, objectPath)
        if err != nil {
            // 缩略图失败不影响主流程，记录日志即可
            logger.Error("Failed to generate thumbnail: %v", err)
        }
    }
    
    // 7. 构建返回结果
    result := &UploadResult{
        OriginalPath:  objectPath,
        URL:           logicalPath,
        RemoteURL:     originalResult.RemoteURL,
        Size:          originalResult.Size,
        Width:         width,
        Height:        height,
        ContentType:   finalContentType,
        Format:        strings.TrimPrefix(filepath.Ext(originalFileName), "."),
        Hash:          originalResult.Hash,
    }
    
    // 添加缩略图信息
    if thumbResult != nil {
        result.ThumbnailPath = thumbResult.OriginalPath
        result.ThumbnailURL = thumbResult.URL
        result.RemoteThumbURL = thumbResult.RemoteURL
    }
    
    return result, nil
}

// buildObjectPath 构建对象路径
func (a *YourStorageAdapter) buildObjectPath(userID uint, folderPath, fileName string) string {
    return utils.BuildObjectPath(userID, folderPath, fileName)
}
```

#### 2.3 文件操作方法
```go
// ReadFile 读取文件内容
func (a *YourStorageAdapter) ReadFile(ctx context.Context, path string) (io.ReadCloser, error) {
    obj, err := a.client.GetObject(ctx, a.bucket, path)
    if err != nil {
        return nil, a.handleError(err)
    }
    return obj.Body, nil
}

// Delete 删除文件
func (a *YourStorageAdapter) Delete(ctx context.Context, path string) error {
    err := a.client.DeleteObject(ctx, a.bucket, path)
    return a.handleError(err)
}

// Exists 检查文件是否存在
func (a *YourStorageAdapter) Exists(ctx context.Context, path string) (bool, error) {
    _, err := a.client.HeadObject(ctx, a.bucket, path)
    if err != nil {
        if a.isNotFoundError(err) {
            return false, nil
        }
        return false, a.handleError(err)
    }
    return true, nil
}
```

#### 2.4 URL生成方法
```go
// GetURL 获取访问URL
func (a *YourStorageAdapter) GetURL(path string, options *URLOptions) (string, error) {
    if options != nil && options.Expires > 0 {
        // 生成签名URL
        return a.generateSignedURL(path, options.Expires)
    }
    
    // 生成公开URL
    baseURL := fmt.Sprintf("https://%s.%s", a.bucket, a.endpoint)
    return fmt.Sprintf("%s/%s", baseURL, path), nil
}

// 完整URL/ CDN URL 统一由上层 URL 策略生成，此处不再提供
```

#### 2.5 Base64方法（已由 Manager 统一实现）
```go
// 适配器无需实现 Base64 接口。使用 Manager：
// storage.GetBase64(ctx, channelID, objectKey)
// storage.GetThumbnailBase64(ctx, channelID, objectKey)
```

#### 2.6 权限控制
```go
// SetObjectACL 设置对象访问权限
func (a *YourStorageAdapter) SetObjectACL(ctx context.Context, path string, acl string) error {
    // 根据存储服务的ACL支持情况实现
    return a.client.SetObjectACL(ctx, a.bucket, path, acl)
}
```

### 步骤3：实现辅助方法

```go
// validateFile 验证文件
func (a *YourStorageAdapter) validateFile(req *UploadRequest) error {
    // 检查文件大小
    if req.File != nil && req.File.Size > a.GetCapabilities().MaxFileSize {
        return NewStorageError(ErrorTypeQuotaExceeded, "file too large", nil)
    }
    
    // 检查文件格式
    ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(req.FileName)), ".")
    supportedFormats := a.GetCapabilities().SupportedFormats
    if !contains(supportedFormats, ext) {
        return NewStorageError(ErrorTypeInvalidFormat, "unsupported format", nil)
    }
    
    return nil
}

// processImage 处理图片
func (a *YourStorageAdapter) processImage(src io.Reader, req *UploadRequest) (io.Reader, string, int, int, error) {
    // 实现图片压缩、格式转换等处理逻辑
    // 返回处理后的数据流、内容类型、宽度、高度
}

// generateThumbnail 生成缩略图
func (a *YourStorageAdapter) generateThumbnail(src io.Reader, req *UploadRequest, originalPath string) (*UploadResult, error) {
    // 实现缩略图生成逻辑
}

// handleError 统一错误处理
func (a *YourStorageAdapter) handleError(err error) error {
    if err == nil {
        return nil
    }
    
    // 根据具体的SDK错误类型，转换为统一的StorageError
    if a.isNotFoundError(err) {
        return NewStorageError(ErrorTypeNotFound, "file not found", err)
    }
    if a.isPermissionError(err) {
        return NewStorageError(ErrorTypePermission, "permission denied", err)
    }
    
    return NewStorageError(ErrorTypeInternal, "storage operation failed", err)
}

// isNotFoundError 检查是否为文件不存在错误
func (a *YourStorageAdapter) isNotFoundError(err error) bool {
    // 根据具体SDK的错误类型判断
    return strings.Contains(err.Error(), "NoSuchKey") || 
           strings.Contains(err.Error(), "not found")
}

// isPermissionError 检查是否为权限错误
func (a *YourStorageAdapter) isPermissionError(err error) bool {
    // 根据具体SDK的错误类型判断
    return strings.Contains(err.Error(), "Access Denied") ||
           strings.Contains(err.Error(), "permission")
}
```

### 步骤4：注册适配器

在 `pkg/storage/factory/factory.go` 中注册新适配器：

```go
func init() {
    RegisterAdapter("your_storage", func() adapter.StorageAdapter {
        return &adapter.YourStorageAdapter{}
    })
}
```

### 步骤5：配置示例

```yaml
# 在配置文件中添加新存储配置
storage_channels:
  your_storage_channel:
    name: "Your Storage"
    type: "your_storage"
    config:
      bucket: "your-bucket-name"
      region: "your-region"
      endpoint: "your-endpoint.com"
      access_key: "your-access-key"
      secret_key: "your-secret-key"
      cdn_domain: "cdn.yourdomain.com"  # 可选
    options:
      hide_remote_url: false
      access_control: "public"
      force_https: true
```

## 🧪 测试指南

### 单元测试示例

```go
func TestYourStorageAdapter_Upload(t *testing.T) {
    adapter := &YourStorageAdapter{}
    
    // 初始化适配器
    config := map[string]interface{}{
        "bucket":     "test-bucket",
        "access_key": "test-key",
        "secret_key": "test-secret",
    }
    
    err := adapter.Initialize(config)
    assert.NoError(t, err)
    
    // 创建测试文件
    file := createTestImageFile(t)
    defer file.Close()
    
    req := &UploadRequest{
        File:       file,
        UserID:     1,
        FolderPath: "test",
        FileName:   "test.jpg",
        ContentType: "image/jpeg",
    }
    
    // 执行上传
    result, err := adapter.Upload(context.Background(), req)
    assert.NoError(t, err)
    assert.NotEmpty(t, result.URL)
    assert.NotEmpty(t, result.RemoteURL)
}
```

### 集成测试

```go
func TestYourStorageAdapter_Integration(t *testing.T) {
    // 测试完整的上传-读取-删除流程
    adapter := setupTestAdapter(t)
    
    // 上传文件
    result := uploadTestFile(t, adapter)
    
    // 验证文件存在
    exists, err := adapter.Exists(context.Background(), result.OriginalPath)
    assert.NoError(t, err)
    assert.True(t, exists)
    
    // 读取文件内容
    reader, err := adapter.ReadFile(context.Background(), result.OriginalPath)
    assert.NoError(t, err)
    defer reader.Close()
    
    // 删除文件
    err = adapter.Delete(context.Background(), result.OriginalPath)
    assert.NoError(t, err)
    
    // 验证文件已删除
    exists, err = adapter.Exists(context.Background(), result.OriginalPath)
    assert.NoError(t, err)
    assert.False(t, exists)
}
```

## ⚠️ 注意事项

### 1. 错误处理
- 统一使用 `StorageError` 类型
- 正确分类错误类型（NotFound、Permission、Network等）
- 提供有意义的错误消息

### 2. 性能优化
- 使用连接池管理SDK客户端
- 实现合理的超时设置
- 考虑并发安全性

### 3. 安全考虑
- 敏感信息（密钥）不要记录到日志
- 验证所有输入参数
- 正确处理权限控制

### 4. 兼容性
- 遵循现有的路径规范
- 保持接口向后兼容
- 考虑不同版本SDK的兼容性

## 📖 参考资源

- [接口定义](./adapter.go)
- [本地存储实现](./local.go)
- [腾讯云COS实现](./cos.go)
- [阿里云OSS实现](./oss.go)
- [雨云S3实现](./rainyun.go)
- 路径工具：`pkg/storage/path`（CleanObjectPath / EnsureObjectKey / ExtractObjectPathFromURL）

---

*如有疑问或需要协助，请联系开发团队或查看相关文档。*

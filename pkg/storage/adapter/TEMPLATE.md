# 新渠道适配器接入模板（示例）

> 本文提供最小可用的适配器实现骨架，便于拷贝后按需补充。请按注释替换 `your_storage` 相关标识。

## 1) 适配器骨架

```go
package adapter

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "path/filepath"
    "strings"

    "pixelpunk/pkg/storage/config"
    "pixelpunk/pkg/storage/tenant"
    "pixelpunk/pkg/storage/utils"
)

type YourStorageAdapter struct {
    // SDK 客户端与配置
    bucket       string
    region       string
    endpoint     string
    accessKey    string
    secretKey    string
    customDomain string
    useHTTPS     bool

    initialized bool
}

var _ StorageAdapter = (*YourStorageAdapter)(nil)

func NewYourStorageAdapter() StorageAdapter { return &YourStorageAdapter{} }

func (a *YourStorageAdapter) GetType() string { return "your_storage" }

func (a *YourStorageAdapter) Initialize(configData map[string]interface{}) error {
    cfg := config.NewMapConfig(configData)
    a.bucket       = cfg.GetString("bucket")
    a.region       = cfg.GetString("region")
    a.endpoint     = cfg.GetString("endpoint")
    a.accessKey    = cfg.GetString("access_key")
    a.secretKey    = cfg.GetString("secret_key")
    a.customDomain = cfg.GetString("custom_domain")
    a.useHTTPS     = cfg.GetBoolWithDefault("use_https", true)

    if a.bucket == "" || a.accessKey == "" || a.secretKey == "" {
        return NewStorageError(ErrorTypeInternal, "missing required configs", nil)
    }
    // TODO: 构造 SDK 客户端（略）
    a.initialized = true
    return nil
}

func (a *YourStorageAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
    if !a.initialized { return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil) }

    // 1) 读取最终入库数据（优先使用 ProcessedData）
    var data []byte
    if len(req.ProcessedData) > 0 {
        data = req.ProcessedData
    } else {
        f, err := req.File.Open(); if err != nil { return nil, err }
        defer f.Close()
        data, err = io.ReadAll(f); if err != nil { return nil, err }
    }

    // 2) 生成对象键与逻辑路径
    objectKey, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, req.FileName)
    if err != nil { return nil, NewStorageError(ErrorTypeInternal, "build object key", err) }
    logicalPath := utils.BuildLogicalPath(req.FolderPath, req.FileName)

    // 3) 调用 SDK 上传（示意）
    // _, err = sdk.PutObject(ctx, a.bucket, objectKey, bytes.NewReader(data), ...)
    // if err != nil { return nil, NewStorageError(ErrorTypeNetwork, "upload failed", err) }

    // 4) 结果
    return &UploadResult{
        OriginalPath:   objectKey,
        URL:            logicalPath,
        RemoteURL:      objectKey,
        Size:           int64(len(data)),
        ContentType:    a.getContentType(filepath.Ext(req.FileName)),
        Format:         strings.TrimPrefix(strings.ToLower(filepath.Ext(req.FileName)), "."),
    }, nil
}

func (a *YourStorageAdapter) ReadFile(ctx context.Context, path string) (io.ReadCloser, error) {
    if !a.initialized { return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil) }
    // 示例：读取对象
    // resp, err := sdk.GetObject(ctx, a.bucket, path)
    // if err != nil { return nil, NewStorageError(ErrorTypeNotFound, "not found", err) }
    // return resp.Body, nil
    return nil, NewStorageError(ErrorTypeInternal, "not implemented", nil)
}

func (a *YourStorageAdapter) Delete(ctx context.Context, path string) error {
    if !a.initialized { return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil) }
    // 示例：删除对象
    // _, err := sdk.DeleteObject(ctx, a.bucket, path)
    // return err
    return nil
}

func (a *YourStorageAdapter) Exists(ctx context.Context, path string) (bool, error) {
    if !a.initialized { return false, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil) }
    // 示例：Head 对象
    // _, err := sdk.HeadObject(ctx, a.bucket, path)
    // if isNotFound(err) { return false, nil }
    // return err == nil, err
    return false, nil
}

func (a *YourStorageAdapter) GetURL(path string, options *URLOptions) (string, error) {
    if !a.initialized { return "", NewStorageError(ErrorTypeInternal, "adapter not initialized", nil) }
    scheme := "https"; if options != nil && !options.ForceHTTPS && !a.useHTTPS { scheme = "http" }
    if a.customDomain != "" { return fmt.Sprintf("%s://%s/%s", scheme, a.customDomain, path), nil }
    if a.endpoint != "" { return fmt.Sprintf("%s://%s/%s", scheme, a.endpoint, path), nil }
    return "/" + path, nil
}

// 完整URL/ CDN URL 统一由上层 URL 策略生成，此处不再提供

func (a *YourStorageAdapter) SetObjectACL(ctx context.Context, path string, acl string) error {
    // 如 SDK 支持：设置对象 ACL（public-read/private）
    return nil
}

func (a *YourStorageAdapter) HealthCheck(ctx context.Context) error {
    if !a.initialized { return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil) }
    // 示例：列目录/Head bucket 验证连通性
    return nil
}

func (a *YourStorageAdapter) GetCapabilities() Capabilities {
    return Capabilities{
        SupportsSignedURL: true,
        SupportsCDN:       true,
        SupportsResize:    false,
        SupportsWebP:      true,
        MaxFileSize:       5 * 1024 * 1024 * 1024,
        SupportedFormats:  []string{"jpg","jpeg","png","gif","webp"},
    }
}

func (a *YourStorageAdapter) getContentType(ext string) string {
    switch strings.ToLower(strings.TrimPrefix(ext, ".")) {
    case "jpg","jpeg": return "image/jpeg"
    case "png": return "image/png"
    case "gif": return "image/gif"
    case "webp": return "image/webp"
    default: return "application/octet-stream"
    }
}
```

## 2) 在工厂中注册

```go
// pkg/storage/factory/factory.go
func init() {
    RegisterAdapter("your_storage", func() adapter.StorageAdapter {
        return adapter.NewYourStorageAdapter()
    })
}
```

## 3) 配置样例

```yaml
storage_channels:
  my-your-storage:
    name: "Your Storage"
    type: "your_storage"
    config:
      bucket: "your-bucket"
      region: "your-region"
      endpoint: "your.endpoint.com"
      access_key: "***"
      secret_key: "***"
      custom_domain: "cdn.example.com" # 可选
```

> 注意：对象键构造不要手写，统一使用 `tenant.BuildObjectKey/BuildThumbObjectKey` 或 `path.EnsureObjectKey`。

# pkg/storage/path 工具说明

提供对象键与 URL 的标准化工具，避免在各处手工拼接路径，降低错误与维护成本。

## 函数清单

- `CleanObjectPath(p string) string`
  - 统一分隔符为 `/`、去重连续斜杠、去首尾斜杠
  - 示例：`CleanObjectPath("//images//a/b/") == "images/a/b"`

- `IsHTTPURL(s string) bool`
  - 判断字符串是否为 `http://` 或 `https://` URL

- `ExtractObjectPathFromURL(fullURL string) string`
  - 从完整直链 URL 中提取对象键（`images/...` 或 `thumbnails/...`）
  - 兼容 OSS/COS/S3 兼容制式（域名差异）

- `EnsureObjectKey(userID uint, input string, isThumb bool) string`
  - 将任意输入（对象键/逻辑路径/完整URL）规范为对象键
  - 自动识别 `thumb/` 前缀或文件名包含 `_thumb.` 判断缩略图
  - 登录用户按别名分片；游客（userID=0）固定 `images/guest/file` 与 `thumbnails/guest/file`

## 使用建议

- 所有“路径或 URL → 对象键”的转换，统一调用 `EnsureObjectKey`
- 构造对象键时优先用 `tenant.BuildObjectKey/BuildThumbObjectKey`（`EnsureObjectKey` 内部已使用）
- 避免在业务代码手工 `fmt.Sprintf("images/...", ...)` 拼接


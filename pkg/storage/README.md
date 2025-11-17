# pkg/storage 模块说明（新版）

本模块提供统一的存储抽象，覆盖适配器注册、选择与调度、路径规则、URL 生成以及常用工具。目标是让“新增一个第三方渠道”变得简单、可维护、且行为一致。

## 目录结构

```
pkg/storage/
├── adapter/   # 具体渠道适配器与接口定义
├── factory/   # 适配器工厂与注册
├── manager/   # 存储管理器（Adapter 缓存、调度、Base64 统一编码）
├── path/      # 路径和 URL 工具（CleanObjectPath / EnsureObjectKey / ExtractObjectPathFromURL）
├── url/       # URL 策略（代理/直链/CDN/HTTPS）
├── compat.go  # 兼容层（历史接口的过渡，逐步收敛）
└── storage.go # 对外统一门面
```

## 核心流程

- 上传：`storage.Storage.Upload` → `manager.StorageManager.Upload` → 具体 `adapter.Upload`
- 删除：`storage.Storage.Delete` → `manager.Delete` → `adapter.Delete`
- URL：`storage.Storage.GetURL/GetFullURL/GetCDNURL` → `manager` → `adapter.GetURL/...`
- Base64：统一在 `manager` 层使用 `adapter.ReadFile` 读取后编码（适配器无需实现 Base64）

## 对象键规范（Object Key）

- 登录用户：按用户稳定别名 + 两位分片构造
  - 原图：`images/<shard>/<alias>/(folder)/file`
  - 缩略图：`thumbnails/<shard>/<alias>/(folder)/file`
- 游客（userID=0）：固定目录（简化）
  - 原图：`images/guest/file`
  - 缩略图：`thumbnails/guest/file`

不要手写对象键！统一使用：
- `tenant.BuildObjectKey(userID, folder, file)`
- `tenant.BuildThumbObjectKey(userID, folder, file)`
- 当输入为逻辑路径或URL时，使用 `path.EnsureObjectKey(userID, input, isThumb)`

辅助工具：
- `path.CleanObjectPath(p)` 清理对象键（斜杠与冗余）
- `path.IsHTTPURL(s)` 判断是否是 http/https URL
- `path.ExtractObjectPathFromURL(s)` 从完整 URL 提取对象键

## Base64 统一（适配器无需实现）

- Manager 层提供统一的 Base64 编码：
  - `storage.Storage.GetBase64(ctx, channelID, objectKey)`
  - `storage.Storage.GetThumbnailBase64(ctx, channelID, objectKey)`
- 适配器只需正确实现 `ReadFile`，上层统一编码，减少重复与差异。

## 新增渠道（Adapter）接入步骤

详见 `pkg/storage/adapter/README.md`，关键要点：
- 实现 `StorageAdapter` 接口（无需实现 Base64 方法）。
- `Initialize` 校验配置并创建 SDK 客户端；`GetCapabilities` 如实声明能力。
- `Upload`：优先使用 `req.ProcessedData`，保持“只读一次”的处理与上传；对象键用 `tenant.BuildObjectKey/BuildThumbObjectKey`。
- `ReadFile`：读取对象并返回 `io.ReadCloser`；其他层会复用它实现 Base64 与代理。
- `GetURL/GetFullURL/GetCDNURL`：只做基础拼接（域名/CDN/HTTPS 等策略由 `url` 模块或上层控制）。
- 在 `factory` 中注册适配器类型。

## 删除与本地映射

- 本地适配器将对象键 `images/`、`thumbnails/` 映射到 `uploads/images/`、`uploads/thumbnails/`。
- 本地 `Delete` 兼容对象键、物理绝对路径以及带有 `uploads/` 前缀的相对路径，便于历史兼容与测试。

## 兼容层（compat）

- `compat.go` 提供历史接口到新架构的过渡。新代码应优先使用 `storage.Storage` 与 `manager.StorageManager`。
- 对象键构造与 URL 提取统一迁移至 `path` 包，避免在 compat 内重复实现。

## 日志与错误

- 使用 `pkg/logger` 记录关键字段：`channelID`、`objectKey`、`isThumb`、`userID`。
- 返回 `adapter.NewStorageError` 分类错误，便于上层统一处理。

## 小贴士

- 直链 vs 代理：是否暴露直链、是否强制 HTTPS、是否使用 CDN 由上层策略控制。
- MD5 计算：建议在适配器内部对“最终上传的数据”计算 MD5，以保证入库哈希稳定。
- 游客时效：业务层已限制游客最大 7 天，清理任务由上层调度实现。


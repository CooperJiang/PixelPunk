# 兼容层说明（compat）

本说明文档解释 `pkg/storage/compat.go` 的定位与迁移建议。

## 定位

- 为旧版本（直接依赖“存储提供者”接口）的调用方提供过渡入口。
- 将历史调用适配到新架构（Storage/Manager/Adapter/Path）。
- 在不破坏旧接口的前提下，逐步引导所有调用迁移到新 API。

## 迁移建议（优先使用新接口）

- 上传/删除/URL：改用 `storage.Storage`（或直接使用 `manager.StorageManager`）。
- Base64：改用 `storage.Storage.GetBase64 / GetThumbnailBase64`（Manager 层统一通过 ReadFile 编码）。
- 对象键：改用 `tenant.BuildObjectKey/BuildThumbObjectKey` 或 `path.EnsureObjectKey`，不要在 compat 内重复实现解析或拼接。

## 已经完成的收敛

- URL→对象键提取：统一使用 `path.ExtractObjectPathFromURL`。
- 路径清理：统一使用 `path.CleanObjectPath`。
- 逻辑路径/URL→对象键：统一使用 `path.EnsureObjectKey`。

## 计划与注意

- 当所有上游调用完成迁移时，删除 compat 层或将其移动到 `internal/`，避免新代码依赖。
- 兼容层返回的对象键/URL 行为与新实现保持一致，以减少“相同输入不同输出”的情况。


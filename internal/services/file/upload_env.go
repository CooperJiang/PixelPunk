package file

/* Environment preparation split from upload_service.go (no behavior change). */

import (
	"pixelpunk/internal/services/storage"
	"pixelpunk/pkg/errors"
	pkgStorage "pixelpunk/pkg/storage"
)

func prepareUploadEnvironment(ctx *UploadContext) error {
	if err := pkgStorage.EnsureUserDirectories(ctx.UserID); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "初始化用户目录失败")
	}

	channel, err := storage.GetDefaultChannel()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "获取存储渠道失败")
	}
	ctx.StorageChannel = channel
	ctx.FileID = generateFileID()
	return nil
}

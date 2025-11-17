package storage

/* Connection testing utilities split from storage_channel_service.go (no behavior change). */

import (
	"context"
	"fmt"
	"strings"

	"mime/multipart"

	"pixelpunk/pkg/assets"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/storage/adapter"
	"pixelpunk/pkg/storage/factory"
	"pixelpunk/pkg/storage/manager"
)

func TestConnection(channelID string) error {
	if _, err := GetChannelConfigMap(channelID); err != nil {
		return err
	}

	channel, err := GetChannelByID(channelID)
	if err != nil {
		return err
	}

	return testChannelByUpload(channelID, channel.Type)
}

func testChannelByUpload(channelID, channelType string) error {
	mgr, err := createStorageManager()
	if err != nil {
		return fmt.Errorf("创建存储管理器失败: %w", err)
	}

	testFile, err := createTestFileHeader()
	if err != nil {
		return fmt.Errorf("获取测试文件失败: %w", err)
	}

	contentType := "image/png"
	if testFile.Header != nil {
		if v := testFile.Header.Get("Content-Type"); v != "" {
			contentType = v
		}
	}
	req := &adapter.UploadRequest{
		File:        testFile,
		UserID:      0,
		FolderPath:  "_cyber_test",
		FileName:    testFile.Filename,
		ContentType: contentType,
		Options:     nil,
	}
	result, putErr := mgr.Upload(context.Background(), channelID, req)
	if putErr == nil {
		if result != nil {
			if result.OriginalPath != "" {
				_ = mgr.Delete(context.Background(), channelID, result.OriginalPath)
			}
			if result.ThumbnailPath != "" {
				_ = mgr.Delete(context.Background(), channelID, result.ThumbnailPath)
			}
		}
		return nil
	}

	if mapped := mapConnectivityError(channelType, putErr); mapped != nil {
		return mapped
	}
	return errors.Wrap(putErr, errors.CodeThirdPartyService, "上传测试失败")
}

func mapConnectivityError(channelType string, err error) error {
	if err == nil {
		return nil
	}
	s := err.Error()
	if strings.Contains(s, "AccessDenied") || strings.Contains(s, "Forbidden") || strings.Contains(s, "SignatureDoesNotMatch") {
		return errors.New(errors.CodeUnauthorized, "认证失败：请检查凭据/签名/权限设置")
	}
	lower := strings.ToLower(s)
	if strings.Contains(s, "SignatureDoesNotMatch") {
		return errors.New(errors.CodeUnauthorized, "S3签名不匹配：请检查 AccessKey/SecretKey 是否正确、path-style 选择是否正确、以及自定义域/endpoint 与 region 是否匹配")
	}
	if strings.Contains(s, "PermanentRedirect") || strings.Contains(s, "301") || strings.Contains(lower, "moved permanently") || strings.Contains(lower, "wrong region") {
		return errors.New(errors.CodeInvalidParameter, "S3区域不匹配：请将配置的 Region 与存储桶实际区域保持一致，或使用对应区域的 Endpoint")
	}
	if strings.Contains(s, "InvalidBucketName") {
		return errors.New(errors.CodeInvalidParameter, "S3存储桶名称无效：请检查命名是否符合规范（小写字母、数字、短横线，3-63字符）")
	}
	if strings.Contains(s, "RequestTimeTooSkewed") || strings.Contains(lower, "clock skew") {
		return errors.New(errors.CodeTimeout, "本机时间与服务器存在偏差：请同步系统时间后重试")
	}
	if channelType == "azureblob" {
		if strings.Contains(s, "AuthenticationFailed") || strings.Contains(s, "AuthorizationFailure") {
			return errors.New(errors.CodeUnauthorized, "Azure认证失败：请检查 account_name/account_key 是否正确，若使用 SAS 请确认有效期与权限")
		}
		if strings.Contains(s, "AuthorizationPermissionMismatch") || strings.Contains(s, "InsufficientAccountPermissions") {
			return errors.New(errors.CodeForbidden, "Azure权限不足：请为容器授予写入权限，或使用具备写权限的密钥/SAS")
		}
		if strings.Contains(lower, "container") && strings.Contains(lower, "not exist") {
			return errors.New(errors.CodeNotFound, "容器不存在：请检查 container 名称是否正确，或先在 Azure 门户中创建")
		}
		if strings.Contains(lower, "request body is too large") || strings.Contains(lower, "entity too large") {
			return errors.New(errors.CodeUploadLimitExceeded, "上传体积超限：请调整大小或升级存储限制")
		}
		if strings.Contains(lower, "x-ms-version") {
			return errors.New(errors.CodeInvalidParameter, "x-ms-version 不兼容：请联系管理员升级服务端或使用兼容的 API 版本")
		}
	}
	if strings.Contains(strings.ToLower(s), "timeout") || strings.Contains(strings.ToLower(s), "dial tcp") || strings.Contains(strings.ToLower(s), "connection refused") {
		return errors.New(errors.CodeTimeout, "网络连接失败或超时：请检查网络连通、Endpoint/自定义域名与端口配置")
	}
	if strings.Contains(s, "429") || strings.Contains(strings.ToLower(s), "rate") {
		return errors.New(errors.CodeRateLimited, "请求过于频繁：服务端限流，请稍后再试")
	}
	return nil
}

func createStorageManager() (*manager.StorageManager, error) {
	channelRepo := &testChannelRepository{}

	globalFactory := factory.GetGlobalFactory()

	mgr := manager.NewStorageManagerWithFactory(channelRepo, globalFactory)

	return mgr, nil
}

func createTestFileHeader() (*multipart.FileHeader, error) {
	return assets.GetTestFileHeader()
}

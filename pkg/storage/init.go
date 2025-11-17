package storage

import (
	"pixelpunk/pkg/storage/adapter"
	"pixelpunk/pkg/storage/factory"
)

// init 包初始化函数，自动注册内置适配器
func init() {
	factory.RegisterGlobalAdapter("local", adapter.NewLocalAdapter)

	factory.RegisterGlobalAdapter("oss", adapter.NewOSSAdapter)

	// 注册腾讯云COS存储适配器
	factory.RegisterGlobalAdapter("cos", adapter.NewCOSAdapter)

	// 注册雨云RainyUN存储适配器
	factory.RegisterGlobalAdapter("rainyun", adapter.NewRainyunAdapter)

	// 注册 Cloudflare R2 存储适配器
	factory.RegisterGlobalAdapter("r2", adapter.NewR2Adapter)

	// 注册 七牛云 Kodo (S3 兼容) 存储适配器
	factory.RegisterGlobalAdapter("qiniu", adapter.NewQiniuAdapter)

	// 注册 又拍云 Upyun (REST API) 存储适配器
	factory.RegisterGlobalAdapter("upyun", adapter.NewUpyunAdapter)

	// 注册 通用 S3 适配器（AWS S3 / 兼容 S3）
	factory.RegisterGlobalAdapter("s3", adapter.NewS3Adapter)

	// 注册 Azure Blob 存储适配器（原生 REST）
	factory.RegisterGlobalAdapter("azureblob", adapter.NewAzureBlobAdapter)

	// 注册 WebDAV 存储适配器
	factory.RegisterGlobalAdapter("webdav", adapter.NewWebDAVAdapter)

	// 注册 SFTP 与 FTP 存储适配器
	factory.RegisterGlobalAdapter("sftp", adapter.NewSFTPAdapter)
	factory.RegisterGlobalAdapter("ftp", adapter.NewFTPAdapter)

	// 注册常见 S3 兼容厂商（作为 S3 适配器别名）
	factory.RegisterGlobalAdapter("obs", adapter.NewS3Adapter)      // 华为云 OBS
	factory.RegisterGlobalAdapter("bos", adapter.NewS3Adapter)      // 百度云 BOS
	factory.RegisterGlobalAdapter("ks3", adapter.NewS3Adapter)      // 金山云 KS3
	factory.RegisterGlobalAdapter("tos", adapter.NewS3Adapter)      // 火山引擎 TOS
	factory.RegisterGlobalAdapter("us3", adapter.NewS3Adapter)      // UCloud US3
	factory.RegisterGlobalAdapter("wasabi", adapter.NewS3Adapter)   // Wasabi
	factory.RegisterGlobalAdapter("spaces", adapter.NewS3Adapter)   // DigitalOcean Spaces
	factory.RegisterGlobalAdapter("b2", adapter.NewS3Adapter)       // Backblaze B2 S3
	factory.RegisterGlobalAdapter("linode", adapter.NewS3Adapter)   // Linode/Akamai
	factory.RegisterGlobalAdapter("vultr", adapter.NewS3Adapter)    // Vultr Object Storage
	factory.RegisterGlobalAdapter("scaleway", adapter.NewS3Adapter) // Scaleway Object Storage
	factory.RegisterGlobalAdapter("ovh", adapter.NewS3Adapter)      // OVHcloud Object Storage
	factory.RegisterGlobalAdapter("oci", adapter.NewS3Adapter)      // Oracle Cloud Object Storage
	factory.RegisterGlobalAdapter("ibmcos", adapter.NewS3Adapter)   // IBM Cloud Object Storage
	factory.RegisterGlobalAdapter("qingstor", adapter.NewS3Adapter) // 青云 QingStor
	factory.RegisterGlobalAdapter("nos", adapter.NewS3Adapter)      // 网易 NOS
	factory.RegisterGlobalAdapter("jdcloud", adapter.NewS3Adapter)  // 京东云对象存储
}

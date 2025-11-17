package common

import "net/http"

const (
	CodeTypeRegister      = "register"
	CodeTypeResetPassword = "reset_password"
	CodeTypeChangeEmail   = "change_email"
)

const (
	StatusSuccess             = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusForbidden           = http.StatusForbidden
	StatusNotFound            = http.StatusNotFound
	StatusServerError         = http.StatusInternalServerError
	StatusInternalServerError = http.StatusInternalServerError

	DefaultPageSize = 100
	MaxPageSize     = 1000

	UserStatusNormal   = 1
	UserStatusDisabled = 2
	UserStatusDeleted  = 3

	UserRoleSuperAdmin = 1
	UserRoleAdmin      = 2
	UserRoleUser       = 3
)

const (
	LogLevelSuccess = 1
	LogLevelInfo    = 2
	LogLevelWarning = 3
	LogLevelError   = 4
)

const (
	LogTypeChatDialog = 1
	LogTypeOther      = 2
)

const (
	AITaggingStatusNone    = "none"
	AITaggingStatusPending = "pending"
	AITaggingStatusDone    = "done"
	AITaggingStatusFailed  = "failed"
	AITaggingStatusSkipped = "skipped"
	AITaggingStatusIgnored = "ignored"
)

const (
	TaggingActionAuto   = "auto"
	TaggingActionManual = "manual"
)

const (
	VectorStatusPending    = "pending"
	VectorStatusProcessing = "processing"
	VectorStatusCompleted  = "completed"
	VectorStatusFailed     = "failed"
	VectorStatusReset      = "reset"
	VectorStatusStale      = "stale"
)

const (
	TaggingStatusTrigger   = "trigger"
	TaggingStatusReset     = "reset"
	TaggingStatusRetry     = "retry"
	TaggingStatusScheduled = "scheduled"
	TaggingStatusIgnore    = "ignore"
	TaggingStatusUnignore  = "unignore"
)

const (
	ShareStatusNormal   = 1
	ShareStatusExpired  = 2
	ShareStatusDeleted  = 3
	ShareStatusDisabled = 4
)

const (
	ShareItemTypeFolder = "folder"
	ShareItemTypeFile   = "file"
)

const (
	MessageStatusUnread  = 1
	MessageStatusRead    = 2
	MessageStatusDeleted = 3
)

const (
	MessagePriorityHigh   = 1
	MessagePriorityNormal = 2
	MessagePriorityLow    = 3
)

const (
	MessageTypeSystemMaintenance = "system.maintenance"
	MessageTypeSystemUpdate      = "system.update"

	MessageTypeAccountRegister         = "account.register"
	MessageTypeAccountStorageGranted   = "account.storage_granted"
	MessageTypeAccountBandwidthGranted = "account.bandwidth_granted"

	MessageTypeContentReviewPending  = "content.review_pending"
	MessageTypeContentReviewApproved = "content.review_approved"
	MessageTypeContentReviewRejected = "content.review_rejected"

	MessageTypeFileDeletedByAdmin      = "file.deleted_by_admin"
	MessageTypeFileBatchDeletedByAdmin = "file.batch_deleted_by_admin"
	MessageTypeFileHardDeletedByAdmin  = "file.hard_deleted_by_admin"
	MessageTypeFileExpiryWarning       = "file.expiry_warning"
	MessageTypeFileThumbnailFailed     = "file.thumbnail_failed"

	MessageTypeStorageQuotaWarning   = "storage.quota_warning"
	MessageTypeStorageQuotaExceeded  = "storage.quota_exceeded"
	MessageTypeStorageQuotaIncreased = "storage.quota_increased"
	MessageTypeStorageQuotaDecreased = "storage.quota_decreased"

	MessageTypeSecurityLoginAlert      = "security.login_alert"
	MessageTypeSecurityPasswordChanged = "security.password_changed"

	MessageTypeAPIKeyCreated     = "apikey.created"
	MessageTypeAPIKeyDeleted     = "apikey.deleted"
	MessageTypeAPIKeyRegenerated = "apikey.regenerated"
	MessageTypeAPIKeyDisabled    = "apikey.disabled"
	MessageTypeAPIKeyEnabled     = "apikey.enabled"

	MessageTypeRandomAPICreated  = "random_api.created"
	MessageTypeRandomAPIDeleted  = "random_api.deleted"
	MessageTypeRandomAPIDisabled = "random_api.disabled"
	MessageTypeRandomAPIEnabled  = "random_api.enabled"

	MessageTypeShareExpiryWarning = "share.expiry_warning"
)

const (
	RelatedTypeFile         = "file"
	RelatedTypeReviewLog    = "review_log"
	RelatedTypeUser         = "user"
	RelatedTypeStorage      = "storage"
	RelatedTypeSystemNotice = "system_notice"
)

const (
	ActionTypeView     = "view"
	ActionTypeEdit     = "edit"
	ActionTypeManage   = "manage"
	ActionTypeDownload = "download"
)

const (
	FileInfoCacheExpire       = 15 * 60
	UserSettingsCacheExpire   = 30 * 60
	BandwidthUsageCacheExpire = 3 * 60
)

const (
	AIMaxTagsLimit        = 1000
	AITagPromptMaxDisplay = 1000
	AIDefaultConfidence   = 0.9

	AIUserTagsRatio     = 50
	AICategoryTagsRatio = 33
	AIPopularTagsRatio  = 20

	AICategorizationMaxRetries = 3
	AITaggingMaxRetries        = 3
	AIRequestTimeoutSeconds    = 30

	AIBatchProcessSize     = 10
	AIBatchProcessInterval = 5

	AIQueryUnificationEnabled = true
	AIFallbackOnQueryError    = true
)

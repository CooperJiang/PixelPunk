package types

// This package provides type aliases to centralize commonly used storage types
// without changing existing behaviors. It allows other packages to import a
// single location for UploadRequest/UploadResult/URLOptions/Capabilities while
// staying fully compatible with current adapter definitions.

import "pixelpunk/pkg/storage/adapter"

// Core request/response/options types
type UploadRequest = adapter.UploadRequest
type UploadResult = adapter.UploadResult
type URLOptions = adapter.URLOptions
type Capabilities = adapter.Capabilities

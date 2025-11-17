//go:build !windows
// +build !windows

package decode

import (
	_ "github.com/adrium/goheif" // HEIC/HEIF support (Linux/macOS only)
)

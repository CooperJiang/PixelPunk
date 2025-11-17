//go:build !windows
// +build !windows

package convert

import (
	_ "github.com/adrium/goheif" // Register HEIC decoder (Linux/macOS only)
)

package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/storage/adapter"
	pathutil "pixelpunk/pkg/storage/path"
)

// RemoteReadProvider is a minimal provider used by file serving helpers.
type RemoteReadProvider interface {
	IsDirectAccess() bool
	GetRemoteContent(objectPath string, isThumb bool, userID uint) (io.ReadCloser, string, error)
	GetFileURL(relativePath string, isThumb bool) (string, error)
}

type providerImpl struct {
	ad        adapter.StorageAdapter
	channelID string
}

func (p *providerImpl) IsDirectAccess() bool { return p.ad.GetType() == "local" }

func (p *providerImpl) GetFileURL(relativePath string, isThumb bool) (string, error) {
	// If already a full URL, return as-is
	if pathutil.IsHTTPURL(relativePath) {
		return relativePath, nil
	}
	rel := strings.TrimPrefix(relativePath, "/")
	// Normalize common logical prefixes
	if strings.HasPrefix(rel, "thumb/") {
		rel = strings.TrimPrefix(rel, "thumb/")
		isThumb = true
	}

	// If it already looks like an object key, delegate to adapter
	if strings.HasPrefix(rel, "files/") || strings.HasPrefix(rel, "thumbnails/") {
		return p.ad.GetURL(rel, &adapter.URLOptions{IsThumbnail: isThumb})
	}

	return "", fmt.Errorf("provider_compat已废弃，请使用FileID-based新架构")
}

func (p *providerImpl) GetRemoteContent(objectPath string, isThumb bool, userID uint) (io.ReadCloser, string, error) {
	// Normalize logical path to object key
	key := pathutil.EnsureObjectKey(userID, objectPath, isThumb)
	if key == "" {
		key = strings.TrimPrefix(objectPath, "/")
	}
	reader, err := p.ad.ReadFile(context.Background(), key)
	if err != nil {
		return nil, "", err
	}
	// Infer content type from extension when possible
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(key)), ".")
	ctype := formats.GetContentType(ext)
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	return reader, ctype, nil
}

// GetStorageProviderByChannelID returns a minimal provider backed by current StorageManager adapter.
func GetStorageProviderByChannelID(channelID string) (RemoteReadProvider, error) {
	mgr := New(&CompatChannelRepository{}).GetManager()
	ad, err := mgr.GetAdapter(channelID)
	if err != nil {
		return nil, err
	}
	return &providerImpl{ad: ad, channelID: channelID}, nil
}

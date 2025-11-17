package adapter

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/imagex/iox"
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/storage/utils"
)

// WebDAVAdapter 通过 HTTP WebDAV 实现
type WebDAVAdapter struct {
	httpClient   *http.Client
	endpoint     string // base, e.g. https://dav.example.com/webdav
	username     string
	password     string
	rootPath     string // optional prefix inside webdav share
	customDomain string // for direct link
	useHTTPS     bool
	allowDirect  bool // 是否允许生成直链（否则外层通过代理）
	autoMkdir    bool // 上传前自动 MKCOL 父目录
	initialized  bool
}

func NewWebDAVAdapter() StorageAdapter {
	return &WebDAVAdapter{httpClient: &http.Client{Timeout: 60 * time.Second}}
}
func (a *WebDAVAdapter) GetType() string { return "webdav" }

func (a *WebDAVAdapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)
	ep := strings.TrimSpace(cfg.GetStringWithDefault("endpoint", ""))
	if ep == "" {
		return NewStorageError(ErrorTypeInternal, "endpoint is required", nil)
	}
	a.username = cfg.GetStringWithDefault("username", "")
	a.password = cfg.GetStringWithDefault("password", "")
	a.rootPath = strings.Trim(strings.TrimSpace(cfg.GetStringWithDefault("root_path", "")), "/")
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)
	a.customDomain, a.useHTTPS = normalizeDomainAndScheme(cfg.GetString("custom_domain"), a.useHTTPS)
	a.allowDirect = cfg.GetBoolWithDefault("allow_direct", false)
	a.autoMkdir = cfg.GetBoolWithDefault("mkdir", true)

	// Normalize endpoint
	if !strings.HasPrefix(ep, "http://") && !strings.HasPrefix(ep, "https://") {
		if a.useHTTPS {
			ep = "https://" + ep
		} else {
			ep = "http://" + ep
		}
	}
	ep = strings.TrimRight(ep, "/")
	a.endpoint = ep
	a.initialized = true
	return nil
}

func (a *WebDAVAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	// Read original
	src, err := req.File.Open()
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to open file", err)
	}
	defer src.Close()
	original, err := iox.ReadAllWithLimit(src, iox.DefaultMaxReadBytes)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to read file", err)
	}

	processed, width, height, format := processUploadData(original, req)
	objectKey, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, req.FileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	logicalPath := utils.BuildLogicalPath(req.FolderPath, req.FileName)

	if err := a.webdavPut(ctx, a.fullKey(objectKey), processed, a.getContentType(format)); err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "webdav put failed", err)
	}

	// thumbnail optional
	var thumbPath, thumbLogical, thumbDirect string
	if req.Options != nil && req.Options.GenerateThumb {
		tb, tf := buildThumbnailBytes(original, req)
		thumbName := utils.MakeThumbName(req.FileName, tf)
		thumbKey, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbName)
		if err := a.webdavPut(ctx, a.fullKey(thumbKey), tb, formats.GetContentType(tf)); err == nil {
			thumbPath = thumbKey
			thumbLogical = utils.BuildLogicalPath(req.FolderPath, thumbName)
			if u, _ := a.GetURL(thumbKey, nil); u != "" {
				thumbDirect = u
			}
		}
	}

	sum := md5.Sum(processed)
	direct, _ := a.GetURL(objectKey, nil)
	return &UploadResult{
		OriginalPath:   objectKey,
		ThumbnailPath:  thumbPath,
		URL:            logicalPath,
		ThumbnailURL:   thumbLogical,
		FullURL:        direct,
		FullThumbURL:   thumbDirect,
		RemoteURL:      objectKey,
		RemoteThumbURL: thumbPath,
		Size:           int64(len(processed)),
		Width:          width,
		Height:         height,
		Hash:           fmt.Sprintf("%x", sum),
		ContentType:    a.getContentType(format),
		Format:         format,
	}, nil
}

func (a *WebDAVAdapter) Delete(ctx context.Context, key string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	u := a.resourceURL(a.fullKey(key))
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	a.basicAuth(req)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("webdav delete failed: %s: %s", resp.Status, string(b))
	}
	return nil
}

func (a *WebDAVAdapter) Exists(ctx context.Context, key string) (bool, error) {
	u := a.resourceURL(a.fullKey(key))
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
	a.basicAuth(req)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	// fallback GET 0-0
	gr, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	a.basicAuth(gr)
	gr.Header.Set("Range", "bytes=0-0")
	resp2, err2 := a.httpClient.Do(gr)
	if err2 != nil {
		return false, err2
	}
	defer resp2.Body.Close()
	if resp2.StatusCode/100 == 2 {
		return true, nil
	}
	return false, nil
}

func (a *WebDAVAdapter) ReadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	u := a.resourceURL(a.fullKey(key))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	a.basicAuth(req)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("webdav get failed: %s: %s", resp.Status, string(b))
	}
	return resp.Body, nil
}

func (a *WebDAVAdapter) GetURL(key string, options *URLOptions) (string, error) {
	if !a.initialized {
		return "", NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	if !a.allowDirect {
		return "", fmt.Errorf("direct url not enabled for webdav")
	}
	// prefer custom domain
	scheme := "https"
	if !a.useHTTPS {
		scheme = "http"
	}
	base := a.endpoint
	if a.customDomain != "" {
		base = fmt.Sprintf("%s://%s", scheme, strings.TrimSuffix(a.customDomain, "/"))
	}
	u := a.resourceURLWithBase(a.fullKey(key), base)
	return u.String(), nil
}

func (a *WebDAVAdapter) SetObjectACL(ctx context.Context, p string, acl string) error { return nil }

func (a *WebDAVAdapter) HealthCheck(ctx context.Context) error {
	// HEAD base root
	u := a.resourceURL(a.fullKey(""))
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
	a.basicAuth(req)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("webdav health failed: %s: %s", resp.Status, string(b))
	}
	return nil
}

func (a *WebDAVAdapter) GetCapabilities() Capabilities {
	return Capabilities{SupportsSignedURL: false, SupportsCDN: false, SupportsResize: false, SupportsWebP: true, MaxFileSize: 5 * 1024 * 1024 * 1024, SupportedFormats: []string{"jpg", "jpeg", "png", "gif", "webp"}}
}

// internal helpers
func (a *WebDAVAdapter) getContentType(format string) string { return formats.GetContentType(format) }

func (a *WebDAVAdapter) fullKey(key string) string {
	k := strings.TrimLeft(key, "/")
	if a.rootPath == "" {
		return k
	}
	return a.rootPath + "/" + k
}

func (a *WebDAVAdapter) resourceURL(key string) url.URL {
	return a.resourceURLWithBase(key, a.endpoint)
}

func (a *WebDAVAdapter) resourceURLWithBase(key string, baseStr string) url.URL {
	base, _ := url.Parse(baseStr)
	enc := encodePathSegments(key)
	raw, _ := url.PathUnescape(enc)
	if strings.Trim(raw, "/") == "" {
		return *base
	}
	base.Path = strings.TrimRight(base.Path, "/")
	base.RawPath = base.Path
	base.Path = path.Join(base.Path, raw)
	base.RawPath = strings.TrimRight(base.RawPath, "/") + "/" + enc
	return *base
}

func (a *WebDAVAdapter) basicAuth(req *http.Request) {
	if a.username != "" {
		req.SetBasicAuth(a.username, a.password)
	}
}

func (a *WebDAVAdapter) webdavPut(ctx context.Context, key string, data []byte, contentType string) error {
	// optional mkdir of parent
	if a.autoMkdir {
		_ = a.mkcolParents(ctx, key)
	}
	u := a.resourceURL(key)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), bytes.NewReader(data))
	a.basicAuth(req)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("webdav put failed: %s: %s", resp.Status, string(b))
	}
	return nil
}

func (a *WebDAVAdapter) mkcolParents(ctx context.Context, key string) error {
	// Make collections for each parent
	key = strings.TrimLeft(key, "/")
	parts := strings.Split(key, "/")
	if len(parts) <= 1 {
		return nil
	}
	// Build prefix path
	cur := ""
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "" {
			continue
		}
		if cur == "" {
			cur = parts[i]
		} else {
			cur = cur + "/" + parts[i]
		}
		u := a.resourceURL(cur)
		req, _ := http.NewRequestWithContext(ctx, "MKCOL", u.String(), nil)
		a.basicAuth(req)
		resp, err := a.httpClient.Do(req)
		if err != nil {
			return err
		}
		io.CopyN(io.Discard, resp.Body, 256)
		resp.Body.Close()
		// 201 Created or 405 Method Not Allowed (already exists) are both acceptable
		if resp.StatusCode == 201 || resp.StatusCode == 405 || resp.StatusCode == 301 || resp.StatusCode == 302 {
			continue
		}
		if resp.StatusCode/100 != 2 {
			return fmt.Errorf("mkcol failed at %s: %s", cur, resp.Status)
		}
	}
	return nil
}

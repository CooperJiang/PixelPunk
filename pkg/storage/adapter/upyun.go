package adapter

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"pixelpunk/pkg/imagex/compress"
	"pixelpunk/pkg/imagex/convert"
	"pixelpunk/pkg/imagex/decode"
	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/imagex/iox"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/pipeline"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/storage/utils"
)

// UpyunAdapter 又拍云（REST API）适配器
type UpyunAdapter struct {
	httpClient   *http.Client
	endpoint     string // e.g. https://v0.api.upyun.com
	service      string // bucket/service name
	operator     string
	password     string // raw password or md5 hex (config)
	hmacKey      []byte // HMAC key: md5(password) as 32-hex ASCII per Upyun doc
	customDomain string
	useHTTPS     bool
	initialized  bool
}

func NewUpyunAdapter() StorageAdapter {
	return &UpyunAdapter{httpClient: &http.Client{Timeout: 30 * time.Second}}
}
func (a *UpyunAdapter) GetType() string { return "upyun" }

// Initialize 初始化
func (a *UpyunAdapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)
	a.service = cfg.GetStringWithDefault("service", "")
	a.operator = cfg.GetStringWithDefault("operator", "")
	a.password = cfg.GetStringWithDefault("password", "")
	ep := cfg.GetStringWithDefault("endpoint", "https://v0.api.upyun.com")
	if !strings.HasPrefix(ep, "http") {
		ep = "https://" + ep
	}
	a.endpoint = strings.TrimRight(ep, "/")
	// Normalize custom domain: accept with or without scheme
	rawDomain := strings.TrimSpace(cfg.GetString("custom_domain"))
	if rawDomain != "" {
		if strings.HasPrefix(rawDomain, "http://") || strings.HasPrefix(rawDomain, "https://") {
			if u, err := url.Parse(rawDomain); err == nil {
				a.customDomain = strings.TrimSuffix(u.Host, "/")
				// honor provided scheme via useHTTPS
				a.useHTTPS = (u.Scheme == "https")
			} else {
				a.customDomain = strings.TrimPrefix(strings.TrimPrefix(rawDomain, "https://"), "http://")
			}
		} else {
			a.customDomain = strings.TrimSuffix(rawDomain, "/")
		}
	}
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)

	if a.service == "" || a.operator == "" || a.password == "" {
		return NewStorageError(ErrorTypeInternal, "service/operator/password are required", nil)
	}
	// Derive HMAC key from config password: if it's 32-hex, use as-is; else md5(password) hex
	pwd := strings.TrimSpace(a.password)
	isHex := func(s string) bool {
		if len(s) != 32 {
			return false
		}
		for i := 0; i < 32; i++ {
			c := s[i]
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
		return true
	}
	if isHex(pwd) {
		a.hmacKey = []byte(strings.ToLower(pwd))
	} else {
		sum := md5.Sum([]byte(pwd))
		a.hmacKey = []byte(fmt.Sprintf("%x", sum))
	}
	a.initialized = true
	return nil
}

// Upload 通过 REST API PUT 上传
func (a *UpyunAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	file, err := req.File.Open()
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to open file", err)
	}
	defer file.Close()
	data, err := iox.ReadAllWithLimit(file, iox.DefaultMaxReadBytes)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to read file", err)
	}

	processed := data
	var width, height int
	format := a.getFileFormat(req.FileName)
	if w, h, f, err := decode.DetectFormat(bytes.NewReader(data)); err == nil {
		width, height, format = w, h, func() string {
			if f != "" {
				return f
			}
			return format
		}()
	}

	if req.Options != nil && req.Options.Compress {
		if cr, err := compress.CompressToTargetSize(bytes.NewReader(data), 5.0, &compress.Options{MaxWidth: req.Options.MaxWidth, MaxHeight: req.Options.MaxHeight, Quality: req.Options.Quality}); err == nil {
			if buf, e := io.ReadAll(cr.Reader); e == nil {
				processed = buf
				width, height, format = cr.Width, cr.Height, cr.Format
			}
		}
	}
	if req.Options != nil && req.Options.WebPEnabled {
		if wr, err := convert.ToWebP(processed, convert.WebPOptions{Quality: req.Options.Quality}); err == nil && wr.Converted {
			if buf, e := io.ReadAll(wr.Reader); e == nil {
				processed = buf
				format = "webp"
			}
		}
	}

	objectKey, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, req.FileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	logicalPath := utils.BuildLogicalPath(req.FolderPath, req.FileName)

	// PUT /service/objectKey
	if err := a.restPut(ctx, objectKey, processed, a.getContentType(format)); err != nil {
		return nil, err
	}

	var thumbPath, thumbLogical, thumbDirect string
	if req.Options != nil && req.Options.GenerateThumb {
		tbytes, tformat, _ := pipeline.GenerateOrFallback(processed, pipeline.Options{Width: max1(req.Options.ThumbWidth, 1200), Height: max1(req.Options.ThumbHeight, 900), Quality: max1(req.Options.ThumbQuality, 85), EnableWebP: true, FallbackOnError: true})
		thumbName := utils.MakeThumbName(req.FileName, tformat)
		thumbKey, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbName)
		if err := a.restPut(ctx, thumbKey, tbytes, formats.GetContentType(tformat)); err == nil {
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

// Delete 删除
func (a *UpyunAdapter) Delete(ctx context.Context, pathKey string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	return a.restDelete(ctx, pathKey)
}

// ReadFile 读取用于代理
func (a *UpyunAdapter) ReadFile(ctx context.Context, pathKey string) (io.ReadCloser, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	// 走管理接口 GET 需要签名
	uri := a.restURI(pathKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	// Sign using encoded request path to match actual request
	a.fillAuthHeaders(req, http.MethodGet, a.encodedPath(pathKey), "", "")
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		logger.Error("Upyun ReadFile failed: %s, body=%s", resp.Status, string(b))
		return nil, fmt.Errorf("upyun get failed: %s", resp.Status)
	}
	return resp.Body, nil
}

// GetURL 返回直链（需要自定义域名）。未配置 custom_domain 时返回错误，外层会走代理回退
func (a *UpyunAdapter) GetURL(pathKey string, options *URLOptions) (string, error) {
	if !a.initialized {
		return "", NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	if strings.TrimSpace(a.customDomain) == "" {
		return "", fmt.Errorf("custom domain not configured for upyun")
	}
	scheme := "https"
	if !a.useHTTPS {
		scheme = "http"
	}
	domain := strings.TrimSuffix(a.customDomain, "/")
	// Encode each path segment to support Chinese and spaces
	parts := strings.Split(strings.TrimLeft(pathKey, "/"), "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return fmt.Sprintf("%s://%s/%s", scheme, domain, strings.Join(parts, "/")), nil
}

func (a *UpyunAdapter) SetObjectACL(ctx context.Context, p string, acl string) error { return nil }
func (a *UpyunAdapter) HealthCheck(ctx context.Context) error                        { return nil }

func (a *UpyunAdapter) Exists(ctx context.Context, pathKey string) (bool, error) {
	uri := a.restURI(pathKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
	a.fillAuthHeaders(req, http.MethodHead, a.encodedPath(pathKey), "", "")
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		logger.Warn("Upyun HEAD not ok: %s, body=%s", resp.Status, string(b))
	}
	return resp.StatusCode == http.StatusOK, nil
}

func (a *UpyunAdapter) GetCapabilities() Capabilities {
	return Capabilities{SupportsSignedURL: false, SupportsCDN: true, SupportsResize: false, SupportsWebP: true, MaxFileSize: 5 * 1024 * 1024 * 1024, SupportedFormats: []string{"jpg", "jpeg", "png", "gif", "webp"}}
}

// 内部：REST PUT/DELETE
func (a *UpyunAdapter) restPut(ctx context.Context, key string, data []byte, contentType string) error {
	uri := a.restURI(key)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPut, uri, bytes.NewReader(data))
	contentLen := strconv.Itoa(len(data))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", contentLen)
	// Upyun expects Content-MD5 as 32-char lowercase hex digest
	sum := md5.Sum(data)
	bodyMD5 := fmt.Sprintf("%x", sum)
	req.Header.Set("Content-MD5", bodyMD5)
	// Auto-create parent directories if not exist
	req.Header.Set("mkdir", "true")
	// Sign using encoded path and include optional Content-MD5
	a.fillAuthHeaders(req, http.MethodPut, a.encodedPath(key), bodyMD5, contentType)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		logger.Error("Upyun PUT failed: %s, body=%s", resp.Status, string(b))
		return fmt.Errorf("upyun put failed: %s", resp.Status)
	}
	return nil
}

func (a *UpyunAdapter) restDelete(ctx context.Context, key string) error {
	uri := a.restURI(key)
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	a.fillAuthHeaders(req, http.MethodDelete, a.encodedPath(key), "", "")
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		// 读取但不输出日志，由上层决定是否记录
		_, _ = io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("upyun delete failed: %s", resp.Status)
	}
	return nil
}

func (a *UpyunAdapter) restURI(key string) string {
	// endpoint/service/key (key should be URL-encoded per segment)
	u, _ := url.Parse(a.endpoint)
	enc := a.encodedPath(key)
	raw, _ := url.PathUnescape(enc)
	// Ensure actual HTTP request uses the encoded path exactly
	u.Path = raw
	u.RawPath = enc
	return u.String()
}

func (a *UpyunAdapter) encodedPath(key string) string {
	encodedKey := strings.TrimLeft(key, "/")
	parts := strings.Split(encodedKey, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return path.Join("/", a.service, strings.Join(parts, "/"))
}

func (a *UpyunAdapter) fillAuthHeaders(req *http.Request, method, uri, contentMD5, contentType string) {
	date := gmtDate()
	// Signature string per Upyun REST doc:
	// <Signature> = Base64(HMAC-SHA1(<Password>, <Method>&<URI>&<Date>&<Content-MD5>))
	// Note: Content-MD5 is optional; if empty, omit it and the preceding '&'. Content-Length is NOT part of signature.
	parts := []string{method, uri, date}
	if contentMD5 != "" {
		parts = append(parts, contentMD5)
	}
	raw := strings.Join(parts, "&")
	// key: md5(password) 32-hex ASCII as HMAC key per doc
	mac := hmac.New(sha1.New, a.hmacKey)
	mac.Write([]byte(raw))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	req.Header.Set("Authorization", fmt.Sprintf("UPYUN %s:%s", a.operator, sign))
	req.Header.Set("Date", date)
	if contentMD5 != "" {
		req.Header.Set("Content-MD5", contentMD5)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
}

// gmtDate 生成符合 Upyun 要求的 GMT 日期字符串
func gmtDate() string {
	// time.RFC1123 使用 UTC，需替换为 GMT
	s := time.Now().UTC().Format(time.RFC1123)
	return strings.ReplaceAll(s, "UTC", "GMT")
}

func (a *UpyunAdapter) getFileFormat(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	if len(ext) > 1 {
		return strings.TrimPrefix(ext, ".")
	}
	return "unknown"
}

func (a *UpyunAdapter) getContentType(format string) string { return formats.GetContentType(format) }

func max1(v, d int) int {
	if v > 0 {
		return v
	}
	return d
}

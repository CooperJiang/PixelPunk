package adapter

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/imagex/iox"
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/storage/utils"
)

// AzureBlobAdapter 原生 Azure Blob 存储适配器（基于 REST Shared Key）
type AzureBlobAdapter struct {
	httpClient    *http.Client
	accountName   string
	accountKeyB64 string
	accountKey    []byte // decoded from base64
	container     string
	endpoint      string // default: https://{account}.blob.core.windows.net
	customDomain  string
	useHTTPS      bool
	accessControl string // "public-read" / "private" (private: GetURL returns SAS)
	initialized   bool
}

func NewAzureBlobAdapter() StorageAdapter {
	return &AzureBlobAdapter{httpClient: &http.Client{Timeout: 30 * time.Second}}
}
func (a *AzureBlobAdapter) GetType() string { return "azureblob" }

// Initialize 读取配置
func (a *AzureBlobAdapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)
	a.accountName = strings.TrimSpace(cfg.GetStringWithDefault("account_name", ""))
	a.accountKeyB64 = strings.TrimSpace(cfg.GetStringWithDefault("account_key", ""))
	a.container = strings.TrimSpace(cfg.GetStringWithDefault("container", ""))
	a.endpoint = strings.TrimSpace(cfg.GetStringWithDefault("endpoint", ""))
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)
	a.customDomain, a.useHTTPS = normalizeDomainAndScheme(cfg.GetString("custom_domain"), a.useHTTPS)
	a.accessControl = strings.TrimSpace(cfg.GetString("access_control"))

	if a.accountName == "" || a.container == "" {
		return NewStorageError(ErrorTypeInternal, "account_name/container required", nil)
	}
	if a.accountKeyB64 == "" {
		return NewStorageError(ErrorTypeInternal, "account_key required", nil)
	}
	key, err := base64.StdEncoding.DecodeString(a.accountKeyB64)
	if err != nil {
		return NewStorageError(ErrorTypeInternal, "invalid account_key (not base64)", err)
	}
	a.accountKey = key
	if a.endpoint == "" {
		scheme := "https"
		if !a.useHTTPS {
			scheme = "http"
		}
		a.endpoint = fmt.Sprintf("%s://%s.blob.core.windows.net", scheme, a.accountName)
	}
	a.initialized = true
	return nil
}

// Upload 使用 Put Blob
func (a *AzureBlobAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	file, err := req.File.Open()
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to open file", err)
	}
	defer file.Close()
	original, err := iox.ReadAllWithLimit(file, iox.DefaultMaxReadBytes)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to read file", err)
	}

	processed, width, height, format := processUploadData(original, req)

	objectPath, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, req.FileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	logicalPath := utils.BuildLogicalPath(req.FolderPath, req.FileName)

	if err := a.putBlob(ctx, a.container, objectPath, processed, a.getContentType(format)); err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "azure put blob failed", err)
	}

	var thumbPath, thumbLogical, thumbDirect string
	if req.Options != nil && req.Options.GenerateThumb {
		tb, tf := buildThumbnailBytes(original, req)
		thumbName := utils.MakeThumbName(req.FileName, tf)
		thumbObject, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbName)
		if err := a.putBlob(ctx, a.container, thumbObject, tb, formats.GetContentType(tf)); err == nil {
			thumbPath = thumbObject
			thumbLogical = utils.BuildLogicalPath(req.FolderPath, thumbName)
			if u, _ := a.GetURL(thumbObject, nil); u != "" {
				thumbDirect = u
			}
		}
	}

	sum := md5.Sum(processed)
	direct, _ := a.GetURL(objectPath, nil)
	return &UploadResult{
		OriginalPath:   objectPath,
		ThumbnailPath:  thumbPath,
		URL:            logicalPath,
		ThumbnailURL:   thumbLogical,
		FullURL:        direct,
		FullThumbURL:   thumbDirect,
		RemoteURL:      objectPath,
		RemoteThumbURL: thumbPath,
		Size:           int64(len(processed)),
		Width:          width,
		Height:         height,
		Hash:           fmt.Sprintf("%x", sum),
		ContentType:    a.getContentType(format),
		Format:         format,
	}, nil
}

func (a *AzureBlobAdapter) Delete(ctx context.Context, pathKey string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	u := a.blobURL(pathKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	a.fillAuthSharedKey(req, "", "", a.canonicalizedResource(pathKey))
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("azure delete failed: %s: %s", resp.Status, string(b))
	}
	return nil
}

func (a *AzureBlobAdapter) Exists(ctx context.Context, pathKey string) (bool, error) {
	u := a.blobURL(pathKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
	a.fillAuthSharedKey(req, "", "", a.canonicalizedResource(pathKey))
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return resp.StatusCode/100 == 2, nil
}

func (a *AzureBlobAdapter) ReadFile(ctx context.Context, pathKey string) (io.ReadCloser, error) {
	u := a.blobURL(pathKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if a.accessControl == "private" { // sign via Shared Key on request
		a.fillAuthSharedKey(req, "", "", a.canonicalizedResource(pathKey))
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("azure get failed: %s: %s", resp.Status, string(b))
	}
	return resp.Body, nil
}

func (a *AzureBlobAdapter) GetURL(pathKey string, options *URLOptions) (string, error) {
	if !a.initialized {
		return "", NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	// private with SAS
	if a.accessControl == "private" {
		return a.generateSASURL(pathKey, options)
	}
	scheme := "https"
	if !a.useHTTPS {
		scheme = "http"
	}
	if a.customDomain != "" {
		return fmt.Sprintf("%s://%s/%s", scheme, strings.TrimSuffix(a.customDomain, "/"), encodePathSegments(pathKey)), nil
	}
	u := a.blobURL(pathKey)
	return u.String(), nil
}

func (a *AzureBlobAdapter) SetObjectACL(ctx context.Context, path string, acl string) error {
	return nil
}

func (a *AzureBlobAdapter) HealthCheck(ctx context.Context) error {
	// GET container properties
	u := a.containerURL()
	q := u.Query()
	q.Set("restype", "container")
	q.Set("comp", "properties")
	u.RawQuery = q.Encode()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	a.fillAuthSharedKey(req, "", "", fmt.Sprintf("/%s/%s", a.accountName, a.container))
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("azure health failed: %s: %s", resp.Status, string(b))
	}
	return nil
}

func (a *AzureBlobAdapter) GetCapabilities() Capabilities {
	return Capabilities{SupportsSignedURL: true, SupportsCDN: false, SupportsResize: false, SupportsWebP: true, MaxFileSize: 5 * 1024 * 1024 * 1024, SupportedFormats: []string{"jpg", "jpeg", "png", "gif", "webp"}}
}

// 内部
func (a *AzureBlobAdapter) getContentType(format string) string {
	return formats.GetContentType(format)
}

func (a *AzureBlobAdapter) containerURL() url.URL {
	base, _ := url.Parse(a.endpoint)
	base.Path = path.Join("/", a.container)
	base.RawPath = "/" + url.PathEscape(a.container)
	return *base
}

func (a *AzureBlobAdapter) blobURL(key string) url.URL {
	base := a.containerURL()
	enc := encodePathSegments(key)
	raw, _ := url.PathUnescape(enc)
	base.Path = path.Join(base.Path, raw)
	base.RawPath = base.RawPath + "/" + enc
	return base
}

func (a *AzureBlobAdapter) putBlob(ctx context.Context, container, key string, data []byte, contentType string) error {
	u := a.blobURL(key)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), bytes.NewReader(data))
	req.Header.Set("x-ms-version", "2021-08-06")
	req.Header.Set("x-ms-date", gmtDate())
	req.Header.Set("x-ms-blob-type", "BlockBlob")
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	a.fillAuthSharedKey(req, contentType, fmt.Sprintf("%d", len(data)), a.canonicalizedResource(key))
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("azure put failed: %s: %s", resp.Status, string(b))
	}
	return nil
}

func (a *AzureBlobAdapter) canonicalizedResource(key string) string {
	// "/{account}/{container}/{key}"
	rawKey := strings.TrimLeft(key, "/")
	return "/" + a.accountName + "/" + a.container + "/" + rawKey
}

// fillAuthSharedKey 计算 SharedKey 鉴权头
func (a *AzureBlobAdapter) fillAuthSharedKey(req *http.Request, contentType, contentLength, canonicalizedResource string) {
	// x-ms headers
	if req.Header.Get("x-ms-date") == "" {
		req.Header.Set("x-ms-date", gmtDate())
	}
	if req.Header.Get("x-ms-version") == "" {
		req.Header.Set("x-ms-version", "2021-08-06")
	}
	// Build CanonicalizedHeaders
	var xms []string
	for k, vals := range req.Header {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "x-ms-") {
			// join values with comma
			v := strings.Join(vals, ",")
			xms = append(xms, fmt.Sprintf("%s:%s", lk, strings.TrimSpace(v)))
		}
	}
	sort.Strings(xms)
	canonicalizedHeaders := strings.Join(xms, "\n")
	if canonicalizedHeaders != "" {
		canonicalizedHeaders += "\n"
	}

	// String-To-Sign (Blob Service)
	// VERB\nContent-Encoding\nContent-Language\nContent-Length\nContent-MD5\nContent-Type\nDate\nIf-Modified-Since\nIf-Match\nIf-None-Match\nIf-Unmodified-Since\nRange\nCanonicalizedHeaders+CanonicalizedResource
	// We use x-ms-date instead of Date, so Date is empty
	contentMD5 := ""
	ce, cl, clg := "", contentLength, contentLength
	// For non-zero length put, Content-Length included; for zero, Azure expects empty. Keep as passed.
	_ = ce
	_ = clg
	sts := strings.Join([]string{
		req.Method,
		"", // Content-Encoding
		"", // Content-Language
		cl, // Content-Length
		contentMD5,
		contentType,
		"", // Date
		"", // If-Modified-Since
		"", // If-Match
		"", // If-None-Match
		"", // If-Unmodified-Since
		"", // Range
		canonicalizedHeaders + canonicalizedResource,
	}, "\n")
	mac := hmac.New(sha256.New, a.accountKey)
	mac.Write([]byte(sts))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	req.Header.Set("Authorization", fmt.Sprintf("SharedKey %s:%s", a.accountName, sig))
}

// generateSASURL 生成仅读的 Blob SAS URL（私有访问时使用）
func (a *AzureBlobAdapter) generateSASURL(pathKey string, options *URLOptions) (string, error) {
	// minimal SAS: sp=r, spr=https,http, sv=2021-08-06, se=expiry, sr=b, sig=...
	expiry := time.Now().UTC().Add(time.Hour)
	if options != nil && options.Expires > 0 {
		expiry = time.Now().UTC().Add(time.Duration(options.Expires) * time.Second)
	}
	se := expiry.Format("2006-01-02T15:04:05Z")
	sp := "r"
	sv := "2021-08-06"
	sr := "b"
	// canonicalized resource: "/blob/{account}/{container}/{blob}"
	canonicalized := fmt.Sprintf("/blob/%s/%s/%s", a.accountName, a.container, strings.TrimLeft(pathKey, "/"))
	// StringToSign for service SAS blob (v 2020-02-10+): sp\nst\nse\ncanonicalizedResource\nsi\nip\nspr\nsv\nrscc\nrscd\nrsce\nrscs\nrsct\nse-signature-fields may vary; use common minimal set
	st := "" // no start time
	si, ip, spr := "", "", "https,http"
	rscc, rscd, rsce, rscs, rsct := "", "", "", "", ""
	stringToSign := strings.Join([]string{sp, st, se, canonicalized, si, ip, spr, sv, rscc, rscd, rsce, rscs, rsct}, "\n")
	mac := hmac.New(sha256.New, a.accountKey)
	mac.Write([]byte(stringToSign))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	u := a.blobURL(pathKey)
	q := u.Query()
	q.Set("sv", sv)
	q.Set("spr", spr)
	q.Set("se", se)
	q.Set("sr", sr)
	q.Set("sp", sp)
	q.Set("sig", url.QueryEscape(sig))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

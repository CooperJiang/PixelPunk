package adapter

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"path"
	"strconv"
	"strings"
	"time"

	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/imagex/iox"
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/storage/utils"
)

// FTPAdapter 简化 FTP/FTPS 客户端（被动模式；无外部依赖）
type FTPAdapter struct {
	host     string
	port     int
	username string
	password string
	// TLS/FTPS 相关
	useTLS        bool   // 是否启用 FTPS
	tlsMode       string // "explicit" (AUTH TLS, 端口21) 或 "implicit" (直接 TLS, 端口990)
	tlsSkipVerify bool   // 是否跳过证书校验（仅测试环境使用）
	serverName    string // SNI/证书校验的主机名（默认使用 host）
	passive       bool   // 始终使用被动
	rootPath      string
	customDomain  string
	allowDirect   bool
	useHTTPS      bool
	mkdir         bool
	timeout       time.Duration
	initialized   bool
}

func NewFTPAdapter() StorageAdapter {
	return &FTPAdapter{port: 21, passive: true, mkdir: true, timeout: 30 * time.Second, useHTTPS: true, tlsMode: "explicit"}
}
func (a *FTPAdapter) GetType() string { return "ftp" }

func (a *FTPAdapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)
	a.host = strings.TrimSpace(cfg.GetStringWithDefault("host", ""))
	a.port = cfg.GetIntWithDefault("port", 21)
	a.username = cfg.GetStringWithDefault("username", "anonymous")
	a.password = cfg.GetStringWithDefault("password", "anonymous@")
	a.rootPath = strings.Trim(strings.TrimSpace(cfg.GetStringWithDefault("root_path", "")), "/")
	a.allowDirect = cfg.GetBoolWithDefault("allow_direct", false)
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)
	a.customDomain, a.useHTTPS = normalizeDomainAndScheme(cfg.GetString("custom_domain"), a.useHTTPS)
	a.mkdir = cfg.GetBoolWithDefault("mkdir", true)
	a.useTLS = cfg.GetBoolWithDefault("use_tls", false)
	a.tlsMode = strings.ToLower(strings.TrimSpace(cfg.GetStringWithDefault("tls_mode", "explicit")))
	if a.tlsMode != "explicit" && a.tlsMode != "implicit" {
		a.tlsMode = "explicit"
	}
	a.tlsSkipVerify = cfg.GetBoolWithDefault("tls_skip_verify", false)
	a.serverName = strings.TrimSpace(cfg.GetStringWithDefault("server_name", ""))
	// 若启用隐式 TLS 且未显式设置端口，默认改为 990
	if a.useTLS && a.tlsMode == "implicit" && !cfg.Has("port") {
		a.port = 990
	}
	if a.host == "" {
		return NewStorageError(ErrorTypeInternal, "host required", nil)
	}
	a.initialized = true
	return nil
}

func (a *FTPAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
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
	key, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, req.FileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	logical := utils.BuildLogicalPath(req.FolderPath, req.FileName)
	remote := a.fullPath(key)
	if err := a.ftpStore(ctx, remote, processed); err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "ftp store failed", err)
	}

	var tpath, tlogical, tDirect string
	if req.Options != nil && req.Options.GenerateThumb {
		tb, tf := buildThumbnailBytes(original, req)
		tname := utils.MakeThumbName(req.FileName, tf)
		tkey, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, tname)
		if err := a.ftpStore(ctx, a.fullPath(tkey), tb); err == nil {
			tpath = tkey
			tlogical = utils.BuildLogicalPath(req.FolderPath, tname)
			if u, _ := a.GetURL(tkey, nil); u != "" {
				tDirect = u
			}
		}
	}

	sum := md5.Sum(processed)
	direct, _ := a.GetURL(key, nil)
	return &UploadResult{
		OriginalPath:   key,
		ThumbnailPath:  tpath,
		URL:            logical,
		ThumbnailURL:   tlogical,
		FullURL:        direct,
		FullThumbURL:   tDirect,
		RemoteURL:      key,
		RemoteThumbURL: tpath,
		Size:           int64(len(processed)),
		Width:          width,
		Height:         height,
		Hash:           fmt.Sprintf("%x", sum),
		ContentType:    formats.GetContentType(format),
		Format:         format,
	}, nil
}

func (a *FTPAdapter) Delete(ctx context.Context, key string) error {
	tp, ctrl, err := a.dialCtrl(ctx)
	if err != nil {
		return err
	}
	defer func() { tp.Close(); ctrl.Close() }()
	if err := a.writeLine(tp, "DELE "+a.fullPath(key)); err != nil {
		return err
	}
	code, _, err := a.readCode(tp)
	if err != nil {
		return err
	}
	if code != 250 && code != 200 {
		return fmt.Errorf("dele failed: %d", code)
	}
	return nil
}

func (a *FTPAdapter) Exists(ctx context.Context, key string) (bool, error) {
	tp, ctrl, err := a.dialCtrl(ctx)
	if err != nil {
		return false, err
	}
	defer func() { tp.Close(); ctrl.Close() }()
	if err := a.writeLine(tp, "SIZE "+a.fullPath(key)); err != nil {
		return false, err
	}
	code, _, err := a.readCode(tp)
	if err != nil {
		return false, err
	}
	if code == 213 {
		return true, nil
	}
	return false, nil
}

func (a *FTPAdapter) ReadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	return a.ftpRetrieve(ctx, a.fullPath(key))
}

func (a *FTPAdapter) GetURL(key string, options *URLOptions) (string, error) {
	if !a.allowDirect {
		return "", fmt.Errorf("direct url not enabled for ftp")
	}
	if a.customDomain == "" {
		return "", fmt.Errorf("custom domain required for direct url")
	}
	scheme := "https"
	if !a.useHTTPS {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/%s", scheme, strings.TrimSuffix(a.customDomain, "/"), encodePathSegments(key)), nil
}

func (a *FTPAdapter) SetObjectACL(ctx context.Context, p string, acl string) error { return nil }

func (a *FTPAdapter) HealthCheck(ctx context.Context) error {
	tp, ctrl, err := a.dialCtrl(ctx)
	if err != nil {
		return err
	}
	defer func() { tp.Close(); ctrl.Close() }()
	if err := a.writeLine(tp, "NOOP"); err != nil {
		return err
	}
	code, _, err := a.readCode(tp)
	if err != nil {
		return err
	}
	if code/100 != 2 {
		return fmt.Errorf("ftp health failed: %d", code)
	}
	return nil
}

func (a *FTPAdapter) GetCapabilities() Capabilities {
	return Capabilities{SupportsSignedURL: false, SupportsCDN: false, SupportsResize: false, SupportsWebP: true, MaxFileSize: 5 * 1024 * 1024 * 1024, SupportedFormats: []string{"jpg", "jpeg", "png", "gif", "webp"}}
}

// helpers
func (a *FTPAdapter) fullPath(key string) string {
	k := strings.TrimLeft(key, "/")
	if a.rootPath == "" {
		return "/" + k
	}
	return "/" + a.rootPath + "/" + k
}

// low-level
func (a *FTPAdapter) dialCtrl(ctx context.Context) (*textproto.Conn, net.Conn, error) {
	addr := net.JoinHostPort(a.host, strconv.Itoa(a.port))
	d := net.Dialer{Timeout: a.timeout}
	baseConn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	// 构造 tls.Config
	makeTLS := func() *tls.Config {
		serverName := a.serverName
		if serverName == "" {
			// 若 host 是域名则用于 SNI
			if net.ParseIP(a.host) == nil {
				serverName = a.host
			}
		}
		return &tls.Config{ServerName: serverName, InsecureSkipVerify: a.tlsSkipVerify}
	}

	var ctrl net.Conn = baseConn
	var tp *textproto.Conn

	// 隐式 TLS：先握手再读 banner
	if a.useTLS && a.tlsMode == "implicit" {
		tlsConn := tls.Client(baseConn, makeTLS())
		if err := tlsConn.Handshake(); err != nil {
			tlsConn.Close()
			return nil, nil, fmt.Errorf("ftps implicit handshake failed: %w", err)
		}
		ctrl = tlsConn
	}

	tp = textproto.NewConn(ctrl)
	if code, _, err := a.readCode(tp); err != nil || code/100 != 2 {
		tp.Close()
		ctrl.Close()
		return nil, nil, fmt.Errorf("ftp banner failed: %v", err)
	}

	// 显式 TLS：AUTH TLS -> TLS 握手
	if a.useTLS && a.tlsMode == "explicit" {
		if err := a.writeLine(tp, "AUTH TLS"); err != nil {
			tp.Close()
			ctrl.Close()
			return nil, nil, err
		}
		if code, _, err := a.readCode(tp); err != nil || code/100 != 2 {
			tp.Close()
			ctrl.Close()
			return nil, nil, fmt.Errorf("AUTH TLS failed: %v", err)
		}
		// 升级为 TLS
		// 关闭旧的 textproto，复用底层连接进行 TLS 包装
		// 注意：某些实现需要在 AUTH TLS 后立刻 TLS 握手
		tlsConn := tls.Client(ctrl, makeTLS())
		if err := tlsConn.Handshake(); err != nil {
			tlsConn.Close()
			return nil, nil, fmt.Errorf("ftps explicit handshake failed: %w", err)
		}
		ctrl = tlsConn
		tp = textproto.NewConn(ctrl)
	}

	if err := a.writeLine(tp, "USER "+a.username); err != nil {
		tp.Close()
		ctrl.Close()
		return nil, nil, err
	}
	if code, _, err := a.readCode(tp); err != nil {
		tp.Close()
		ctrl.Close()
		return nil, nil, err
	} else if code == 331 {
		if err := a.writeLine(tp, "PASS "+a.password); err != nil {
			tp.Close()
			ctrl.Close()
			return nil, nil, err
		}
		if code2, _, err := a.readCode(tp); err != nil || code2/100 != 2 {
			tp.Close()
			ctrl.Close()
			return nil, nil, fmt.Errorf("ftp login failed: %v", err)
		}
	} else if code/100 != 2 {
		tp.Close()
		ctrl.Close()
		return nil, nil, fmt.Errorf("ftp login failed: %d", code)
	}

	// 显式/隐式 TLS 下，设置数据通道保护
	if a.useTLS {
		_ = a.writeLine(tp, "PBSZ 0")
		_, _, _ = a.readCode(tp)
		// PROT P (保护数据通道)
		_ = a.writeLine(tp, "PROT P")
		_, _, _ = a.readCode(tp)
	}

	_ = a.writeLine(tp, "TYPE I")
	_, _, _ = a.readCode(tp)
	return tp, ctrl, nil
}

func (a *FTPAdapter) pasvAddr(tp *textproto.Conn) (string, error) {
	if err := a.writeLine(tp, "PASV"); err != nil {
		return "", err
	}
	code, msg, err := a.readCode(tp)
	if err != nil {
		return "", err
	}
	if code != 227 {
		return "", fmt.Errorf("PASV failed: %d", code)
	}
	// parse (h1,h2,h3,h4,p1,p2)
	start := strings.IndexByte(msg, '(')
	end := strings.IndexByte(msg, ')')
	if start == -1 || end == -1 || end <= start+1 {
		return "", fmt.Errorf("invalid PASV: %s", msg)
	}
	parts := strings.Split(msg[start+1:end], ",")
	if len(parts) != 6 {
		return "", fmt.Errorf("invalid PASV parts")
	}
	host := strings.Join(parts[0:4], ".")
	p1, _ := strconv.Atoi(parts[4])
	p2, _ := strconv.Atoi(parts[5])
	port := p1*256 + p2
	return net.JoinHostPort(host, strconv.Itoa(port)), nil
}

func (a *FTPAdapter) ftpStore(ctx context.Context, remotePath string, data []byte) error {
	tp, ctrl, err := a.dialCtrl(ctx)
	if err != nil {
		return err
	}
	defer func() { tp.Close(); ctrl.Close() }()
	if a.mkdir {
		_ = a.ftpMkdirAll(tp, remotePath)
	}
	addr, err := a.pasvAddr(tp)
	if err != nil {
		return err
	}
	if err := a.writeLine(tp, "STOR "+remotePath); err != nil {
		return err
	}
	// connect data
	d, err := net.DialTimeout("tcp", addr, a.timeout)
	if err != nil {
		return err
	}
	// FTPS 数据通道：根据 PROT 决定是否 TLS 包装
	var dataConn net.Conn = d
	if a.useTLS {
		cfg := &tls.Config{ServerName: func() string {
			if a.serverName != "" {
				return a.serverName
			}
			if net.ParseIP(a.host) == nil {
				return a.host
			}
			return ""
		}(), InsecureSkipVerify: a.tlsSkipVerify}
		tlsD := tls.Client(d, cfg)
		if err := tlsD.Handshake(); err != nil {
			tlsD.Close()
			return err
		}
		dataConn = tlsD
	}
	if _, err := dataConn.Write(data); err != nil {
		dataConn.Close()
		return err
	}
	dataConn.Close()
	// read completion
	code, _, err := a.readCode(tp)
	if err != nil {
		return err
	}
	if code/100 != 2 && code != 150 && code != 226 {
		return fmt.Errorf("stor failed: %d", code)
	}
	return nil
}

func (a *FTPAdapter) ftpRetrieve(ctx context.Context, remotePath string) (io.ReadCloser, error) {
	tp, ctrl, err := a.dialCtrl(ctx)
	if err != nil {
		return nil, err
	}
	addr, err := a.pasvAddr(tp)
	if err != nil {
		tp.Close()
		ctrl.Close()
		return nil, err
	}
	if err := a.writeLine(tp, "RETR "+remotePath); err != nil {
		tp.Close()
		ctrl.Close()
		return nil, err
	}
	d, err := net.DialTimeout("tcp", addr, a.timeout)
	if err != nil {
		tp.Close()
		ctrl.Close()
		return nil, err
	}
	// FTPS 数据通道：TLS 包装
	if a.useTLS {
		cfg := &tls.Config{ServerName: func() string {
			if a.serverName != "" {
				return a.serverName
			}
			if net.ParseIP(a.host) == nil {
				return a.host
			}
			return ""
		}(), InsecureSkipVerify: a.tlsSkipVerify}
		tlsD := tls.Client(d, cfg)
		if err := tlsD.Handshake(); err != nil {
			tp.Close()
			ctrl.Close()
			tlsD.Close()
			return nil, err
		}
		d = tlsD
	}
	// We return a ReadCloser that on Close finishes the control flow
	return &ftpReadCloser{Conn: d, tp: tp, ctrl: ctrl}, nil
}

type ftpReadCloser struct {
	net.Conn
	tp   *textproto.Conn
	ctrl net.Conn
}

func (f *ftpReadCloser) Close() error {
	f.Conn.Close()
	// read completion
	_, _, _ = readCode(f.tp)
	f.tp.Close()
	f.ctrl.Close()
	return nil
}

func (a *FTPAdapter) ftpMkdirAll(tp *textproto.Conn, remotePath string) error {
	dir := path.Dir(remotePath)
	if dir == "/" || dir == "." {
		return nil
	}
	// split and accumulate
	parts := strings.Split(strings.Trim(dir, "/"), "/")
	cur := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		cur = cur + "/" + p
		_ = a.writeLine(tp, "MKD "+cur)
		_, _, _ = a.readCode(tp) // 257 success / 550 exists
	}
	return nil
}

func (a *FTPAdapter) writeLine(tp *textproto.Conn, line string) error  { return tp.PrintfLine(line) }
func (a *FTPAdapter) readCode(tp *textproto.Conn) (int, string, error) { return readCode(tp) }

func readCode(tp *textproto.Conn) (int, string, error) {
	line, err := tp.ReadLine()
	if err != nil {
		return 0, "", err
	}
	if len(line) < 3 {
		return 0, line, fmt.Errorf("invalid reply: %s", line)
	}
	code, _ := strconv.Atoi(line[:3])
	// multiline not fully handled; basic support
	if len(line) > 3 && line[3] == '-' {
		// read until code + space
		for {
			l, err := tp.ReadLine()
			if err != nil {
				return code, line, err
			}
			if strings.HasPrefix(l, fmt.Sprintf("%03d ", code)) {
				line = l
				break
			}
		}
	}
	return code, line, nil
}

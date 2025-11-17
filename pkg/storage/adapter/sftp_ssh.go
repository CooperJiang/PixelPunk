package adapter

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"path"
	"strings"
	"time"

	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/imagex/iox"
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/storage/utils"

	"golang.org/x/crypto/ssh"
)

// SFTPAdapter 基于 SSH 的简化实现（通过远程 shell 命令代替 sftp 协议）
type SFTPAdapter struct {
	host         string
	port         int
	username     string
	password     string
	privateKey   string // PEM
	passphrase   string
	rootPath     string
	customDomain string
	useHTTPS     bool
	allowDirect  bool
	mkdir        bool
	timeout      time.Duration
	initialized  bool
}

func NewSFTPAdapter() StorageAdapter {
	return &SFTPAdapter{port: 22, mkdir: true, timeout: 30 * time.Second, useHTTPS: true}
}
func (a *SFTPAdapter) GetType() string { return "sftp" }

func (a *SFTPAdapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)
	a.host = strings.TrimSpace(cfg.GetStringWithDefault("host", ""))
	a.port = cfg.GetIntWithDefault("port", 22)
	a.username = cfg.GetStringWithDefault("username", "")
	a.password = cfg.GetStringWithDefault("password", "")
	a.privateKey = cfg.GetStringWithDefault("private_key", "")
	a.passphrase = cfg.GetStringWithDefault("passphrase", "")
	a.rootPath = strings.Trim(strings.TrimSpace(cfg.GetStringWithDefault("root_path", "")), "/")
	a.allowDirect = cfg.GetBoolWithDefault("allow_direct", false)
	a.mkdir = cfg.GetBoolWithDefault("mkdir", true)
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)
	a.customDomain, a.useHTTPS = normalizeDomainAndScheme(cfg.GetString("custom_domain"), a.useHTTPS)
	if a.host == "" || a.username == "" {
		return NewStorageError(ErrorTypeInternal, "host/username required", nil)
	}
	a.initialized = true
	return nil
}

func (a *SFTPAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
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

	objectKey, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, req.FileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	logicalPath := utils.BuildLogicalPath(req.FolderPath, req.FileName)

	remotePath := a.fullPath(objectKey)
	if a.mkdir {
		if err := a.sshMkdirP(ctx, path.Dir(remotePath)); err != nil {
			return nil, NewStorageError(ErrorTypeInternal, "mkdir failed", err)
		}
	}
	if err := a.sshWriteFile(ctx, remotePath, processed); err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "write failed", err)
	}

	var thumbPath, thumbLogical, thumbDirect string
	if req.Options != nil && req.Options.GenerateThumb {
		tb, tf := buildThumbnailBytes(original, req)
		thumbName := utils.MakeThumbName(req.FileName, tf)
		thumbKey, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbName)
		tRemote := a.fullPath(thumbKey)
		if a.mkdir {
			_ = a.sshMkdirP(ctx, path.Dir(tRemote))
		}
		if err := a.sshWriteFile(ctx, tRemote, tb); err == nil {
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
		ContentType:    formats.GetContentType(format),
		Format:         format,
	}, nil
}

func (a *SFTPAdapter) Delete(ctx context.Context, key string) error {
	return a.sshRun(ctx, fmt.Sprintf("rm -f -- %q", a.fullPath(key)), nil)
}

func (a *SFTPAdapter) Exists(ctx context.Context, key string) (bool, error) {
	err := a.sshRun(ctx, fmt.Sprintf("test -f %q", a.fullPath(key)), nil)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (a *SFTPAdapter) ReadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	return a.sshReadFile(ctx, a.fullPath(key))
}

func (a *SFTPAdapter) GetURL(key string, options *URLOptions) (string, error) {
	if !a.allowDirect {
		return "", fmt.Errorf("direct url not enabled for sftp")
	}
	scheme := "https"
	if !a.useHTTPS {
		scheme = "http"
	}
	if a.customDomain == "" {
		return "", fmt.Errorf("custom domain required for direct url")
	}
	return fmt.Sprintf("%s://%s/%s", scheme, strings.TrimSuffix(a.customDomain, "/"), encodePathSegments(key)), nil
}

func (a *SFTPAdapter) SetObjectACL(ctx context.Context, p string, acl string) error { return nil }

func (a *SFTPAdapter) HealthCheck(ctx context.Context) error {
	return a.sshRun(ctx, "echo ok", nil)
}

func (a *SFTPAdapter) GetCapabilities() Capabilities {
	return Capabilities{SupportsSignedURL: false, SupportsCDN: false, SupportsResize: false, SupportsWebP: true, MaxFileSize: 5 * 1024 * 1024 * 1024, SupportedFormats: []string{"jpg", "jpeg", "png", "gif", "webp"}}
}

// helpers
func (a *SFTPAdapter) fullPath(key string) string {
	k := strings.TrimLeft(key, "/")
	if a.rootPath == "" {
		return "/" + k
	}
	return "/" + a.rootPath + "/" + k
}

func (a *SFTPAdapter) sshClient(ctx context.Context) (*ssh.Client, error) {
	auths := []ssh.AuthMethod{}
	if strings.TrimSpace(a.privateKey) != "" {
		signer, err := a.parsePrivateKey([]byte(a.privateKey), a.passphrase)
		if err != nil {
			return nil, err
		}
		auths = append(auths, ssh.PublicKeys(signer))
	} else {
		auths = append(auths, ssh.Password(a.password))
	}
	cfg := &ssh.ClientConfig{
		User:            a.username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         a.timeout,
	}
	addr := net.JoinHostPort(a.host, fmt.Sprintf("%d", a.port))
	return ssh.Dial("tcp", addr, cfg)
}

func (a *SFTPAdapter) sshRun(ctx context.Context, cmd string, stdin []byte) error {
	cli, err := a.sshClient(ctx)
	if err != nil {
		return err
	}
	defer cli.Close()
	sess, err := cli.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	if len(stdin) > 0 {
		sess.Stdin = bytes.NewReader(stdin)
	}
	return sess.Run(cmd)
}

func (a *SFTPAdapter) sshMkdirP(ctx context.Context, dir string) error {
	return a.sshRun(ctx, fmt.Sprintf("mkdir -p -- %q", dir), nil)
}

func (a *SFTPAdapter) sshWriteFile(ctx context.Context, path string, data []byte) error {
	cmd := fmt.Sprintf("cat > %q", path)
	return a.sshRun(ctx, cmd, data)
}

type sshReadCloser struct {
	io.Reader
	sess *ssh.Session
	cli  *ssh.Client
}

func (s *sshReadCloser) Close() error { _ = s.sess.Close(); return s.cli.Close() }

func (a *SFTPAdapter) sshReadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	cli, err := a.sshClient(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := cli.NewSession()
	if err != nil {
		cli.Close()
		return nil, err
	}
	stdout, err := sess.StdoutPipe()
	if err != nil {
		sess.Close()
		cli.Close()
		return nil, err
	}
	if err := sess.Start(fmt.Sprintf("cat %q", path)); err != nil {
		sess.Close()
		cli.Close()
		return nil, err
	}
	return &sshReadCloser{Reader: stdout, sess: sess, cli: cli}, nil
}

func (a *SFTPAdapter) parsePrivateKey(pemBytes []byte, passphrase string) (ssh.Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid private key")
	}
	key := pemBytes
	if x509.IsEncryptedPEMBlock(block) {
		if passphrase == "" {
			return nil, fmt.Errorf("private key is encrypted, passphrase required")
		}
		decrypted, err := x509.DecryptPEMBlock(block, []byte(passphrase))
		if err != nil {
			return nil, err
		}
		key = pem.EncodeToMemory(&pem.Block{Type: block.Type, Bytes: decrypted})
		return ssh.ParsePrivateKey(key)
	}
	// 处理 OpenSSH 新格式加密私钥（"OPENSSH PRIVATE KEY"）
	if passphrase != "" {
		if signer, err := ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(passphrase)); err == nil {
			return signer, nil
		} else {
			// 回退到无口令解析，返回原始错误提示
			if signer2, err2 := ssh.ParsePrivateKey(pemBytes); err2 == nil {
				return signer2, nil
			}
			return nil, fmt.Errorf("failed to parse private key with passphrase: %v", err)
		}
	}
	return ssh.ParsePrivateKey(key)
}

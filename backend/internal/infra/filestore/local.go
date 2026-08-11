package filestore

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalStore struct {
	root string
}

func NewLocal(root string) (*LocalStore, error) {
	if root == "" {
		root = "./data/files"
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, err
	}
	return &LocalStore{root: abs}, nil
}

func (s *LocalStore) Put(ctx context.Context, bucket, key string, r io.Reader, meta PutMeta) error {
	path, err := s.path(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, readerWithContext{ctx: ctx, r: r})
	return err
}

func (s *LocalStore) Get(ctx context.Context, bucket, key string) (io.ReadCloser, FileMeta, error) {
	path, err := s.path(key)
	if err != nil {
		return nil, FileMeta{}, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, FileMeta{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, FileMeta{}, err
	}
	return file, FileMeta{SizeBytes: info.Size()}, nil
}

func (s *LocalStore) Delete(ctx context.Context, bucket, key string) error {
	path, err := s.path(key)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func (s *LocalStore) Exists(ctx context.Context, bucket, key string) (bool, error) {
	path, err := s.path(key)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (s *LocalStore) SignGetURL(ctx context.Context, bucket, key string, ttl time.Duration, contentDisposition string) (string, bool, error) {
	return "", false, nil
}

// SignPostURL 不受 local 驱动支持：local 驱动仅面向开发环境。生产部署必须使用
// s3 驱动（RustFS/S3），OpenAPI 直传才能签发预签名 POST 表单。
func (s *LocalStore) SignPostURL(ctx context.Context, bucket, key, contentType string, sizeBytes int64, ttl time.Duration) (string, map[string]string, bool, error) {
	return "", nil, false, nil
}

// 分片上传控制面不受 local 驱动支持：local 驱动仅面向开发环境，无法发起分片会话、
// 无法签发分片 PUT 地址。生产部署必须使用 s3 驱动（RustFS/S3）。
func (s *LocalStore) CreateMultipartUpload(ctx context.Context, bucket, key, contentType string) (string, error) {
	return "", ErrMultipartUnsupported
}

func (s *LocalStore) CompleteMultipartUpload(ctx context.Context, bucket, key, uploadID string, parts []CompletePart) error {
	return ErrMultipartUnsupported
}

func (s *LocalStore) AbortMultipartUpload(ctx context.Context, bucket, key, uploadID string) error {
	return ErrMultipartUnsupported
}

func (s *LocalStore) ListParts(ctx context.Context, bucket, key, uploadID string) ([]PartInfo, error) {
	return nil, ErrMultipartUnsupported
}

func (s *LocalStore) SignPutPartURL(ctx context.Context, bucket, key, uploadID string, partNumber int32, ttl time.Duration) (string, bool, error) {
	return "", false, nil
}

func (s *LocalStore) path(key string) (string, error) {
	if key == "" || strings.Contains(key, "\x00") {
		return "", errors.New("invalid file key")
	}
	clean := filepath.Clean(filepath.FromSlash(key))
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") {
		return "", errors.New("invalid file key")
	}
	path := filepath.Join(s.root, clean)
	if !strings.HasPrefix(path, s.root+string(os.PathSeparator)) && path != s.root {
		return "", errors.New("file key escapes root")
	}
	return path, nil
}

type readerWithContext struct {
	ctx context.Context
	r   io.Reader
}

func (r readerWithContext) Read(p []byte) (int, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.r.Read(p)
	}
}

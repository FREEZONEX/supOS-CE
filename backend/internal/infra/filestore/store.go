package filestore

import (
	"context"
	"errors"
	"io"
	"time"
)

type PutMeta struct {
	ContentType string
	SizeBytes   int64
	Sha256      string
}

type FileMeta struct {
	ContentType string
	SizeBytes   int64
	ETag        string
}

// Store abstracts object storage. All methods take an explicit bucket so the
// two-bucket layout (public bucket for anonymous-read files, private bucket
// for presigned-only files) can address the right bucket per file. An empty
// bucket means the driver's default (private) bucket; the local driver
// ignores the bucket entirely.
type Store interface {
	Put(ctx context.Context, bucket, key string, r io.Reader, meta PutMeta) error
	Get(ctx context.Context, bucket, key string) (io.ReadCloser, FileMeta, error)
	Delete(ctx context.Context, bucket, key string) error
	Exists(ctx context.Context, bucket, key string) (bool, error)
	SignGetURL(ctx context.Context, bucket, key string, ttl time.Duration, contentDisposition string) (string, bool, error)
	// SignPostURL 签发预签名 POST（POST policy）上传表单：postURL 为客户端直接 POST
	// （multipart/form-data）的地址，fields 为需随表单原样提交的隐藏字段（含 policy
	// 与签名）。policy 将对象 key 精确匹配与 content-length-range [1, sizeBytes]
	// 签入，对象超限由对象存储直接拒绝；contentType 非空时同时签入 Content-Type
	// 条件（随 fields 下发，客户端回显后对象保留声明的媒体类型），为空则不约束；
	// ok=false 表示当前驱动不支持签发。
	SignPostURL(ctx context.Context, bucket, key, contentType string, sizeBytes int64, ttl time.Duration) (postURL string, fields map[string]string, ok bool, err error)

	// CreateMultipartUpload 发起一次分片上传并返回对象存储侧的 uploadId，后续
	// 分片 PUT、complete、abort 均以 (bucket, key, uploadId) 定位本次上传会话。
	// local 驱动不支持分片，返回 ErrMultipartUnsupported。
	CreateMultipartUpload(ctx context.Context, bucket, key, contentType string) (uploadID string, err error)
	// CompleteMultipartUpload 用客户端回传的 (partNumber, etag) 列表合并已上传分片。
	// parts 的 etag 来自各分片 PUT 响应，合并由对象存储按 uploadId 校验。
	CompleteMultipartUpload(ctx context.Context, bucket, key, uploadID string, parts []CompletePart) error
	// AbortMultipartUpload 丢弃未完成的分片上传会话，已上传分片一并清理。
	AbortMultipartUpload(ctx context.Context, bucket, key, uploadID string) error
	// ListParts 列出当前上传会话已上传的分片（含真实大小），complete 前据此核对
	// 实际上传字节数，不信任客户端自报大小。local 驱动返回 ErrMultipartUnsupported。
	ListParts(ctx context.Context, bucket, key, uploadID string) ([]PartInfo, error)
	// SignPutPartURL 签发单个分片的预签名 PUT 地址（query 携带 partNumber 与
	// uploadId）；ok=false 表示当前驱动不支持签发。
	SignPutPartURL(ctx context.Context, bucket, key, uploadID string, partNumber int32, ttl time.Duration) (url string, ok bool, err error)
}

// CompletePart 是 complete 时客户端回传的分片元数据：partNumber 对应上传时的
// 分片序号，ETag 取自该分片 PUT 响应头。
type CompletePart struct {
	PartNumber int32
	ETag       string
}

// PartInfo 是 ListParts 返回的单个已上传分片信息，SizeBytes 为对象存储记录的真实大小。
type PartInfo struct {
	PartNumber int32
	SizeBytes  int64
	ETag       string
}

// ErrMultipartUnsupported 表示当前存储驱动不支持分片上传控制面（local 仅面向开发环境）。
var ErrMultipartUnsupported = errors.New("multipart upload is not supported by the local storage driver; configure the s3/RustFS storage driver")

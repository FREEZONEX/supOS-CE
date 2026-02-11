package types

import (
	"mime/multipart"
)

// UploadAttachmentRequest 上传附件请求
type UploadAttachmentRequest struct {
	FileHeader  *multipart.FileHeader `form:"file"`
	Category    string                `form:"category,optional"`
	Description string                `form:"description,optional"`
	Tags        []string              `form:"tags,optional"`
}

// DownloadAttachmentRequest 下载附件请求
type DownloadAttachmentRequest struct {
	FileID   string `form:"fileId"`
	Filename string `form:"filename,optional"`
}

// DeleteAttachmentRequest 删除附件请求
type DeleteAttachmentRequest struct {
	FileID string `form:"fileId"`
}

// ListAttachmentsRequest 列出附件请求
type ListAttachmentsRequest struct {
	Category string `form:"category,optional"`
	Tag      string `form:"tag,optional"`
	Page     int    `form:"page,optional,default=1"`
	PageSize int    `form:"pageSize,optional,default=20"`
}

// UploadAttachmentResponse 上传附件响应
type UploadAttachmentResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"msg,optional"`
	Data    UploadResult `json:"data"`
}

// DeleteAttachmentResponse 删除附件响应
type DeleteAttachmentResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message,optional"`
	Data    DeleteResult `json:"data"`
}

// ListAttachmentsResponse 列出附件响应
type ListAttachmentsResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message,optional"`
	Data    ListResult `json:"data"`
	Total   int        `json:"total"`
}

// UploadResult 上传结果
type UploadResult struct {
	FileID      string `json:"fileId"`
	Filename    string `json:"filename"`
	FileSize    int64  `json:"fileSize"`
	FileType    string `json:"fileType"`
	MD5         string `json:"md5"`
	UploadTime  string `json:"uploadTime"`
	StoragePath string `json:"storagePath"`
}

// DeleteResult 删除结果
type DeleteResult struct {
	FileID  string `json:"fileId"`
	Success bool   `json:"success"`
}

// ListResult 列表结果
type ListResult struct {
	Items      []AttachmentInfo `json:"items"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalPages int              `json:"totalPages"`
}

// AttachmentInfo 附件信息
type AttachmentInfo struct {
	FileID        string   `json:"fileId"`
	Filename      string   `json:"filename"`
	FileSize      int64    `json:"fileSize"`
	FileType      string   `json:"fileType"`
	Category      string   `json:"category,optional"`
	Description   string   `json:"description,optional"`
	Tags          []string `json:"tags,optional"`
	MD5           string   `json:"md5"`
	UploadTime    string   `json:"uploadTime"`
	Uploader      string   `json:"uploader,optional"`
	DownloadCount int64    `json:"downloadCount"`
}

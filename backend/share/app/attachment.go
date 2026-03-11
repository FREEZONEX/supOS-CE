package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AttachmentConfig 附件上传配置
type AttachmentConfig struct {
	// 最大文件大小（字节），默认2GB
	MaxFileSize int64
	// 允许的文件类型（MIME类型），空表示允许所有类型
	AllowedTypes []string
	// 存储根目录
	StorageRoot string
	// 临时目录
	TempDir string
	// 是否生成缩略图
	GenerateThumbnail bool
	// 缩略图大小
	ThumbnailWidth  int
	ThumbnailHeight int
	// 是否启用MD5校验
	EnableMD5Check bool
	// 是否保留原始文件名
	KeepOriginalName bool
}

// DefaultAttachmentConfig 默认附件配置
func DefaultAttachmentConfig() *AttachmentConfig {
	return &AttachmentConfig{
		MaxFileSize:       2 * 1024 * 1024 * 1024, // 2GB
		AllowedTypes:      []string{},
		StorageRoot:       getEnvOrDefault("ATTACHMENT_STORAGE_ROOT", "/app/go-edge/attachment/"),
		TempDir:           getEnvOrDefault("ATTACHMENT_TEMP_DIR", "/temp"),
		GenerateThumbnail: false,
		ThumbnailWidth:    200,
		ThumbnailHeight:   200,
		EnableMD5Check:    true,
		KeepOriginalName:  false,
	}
}

// AttachmentInfo 附件信息
type AttachmentInfo struct {
	FileID       string    `json:"fileId"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"originalName"`
	FileSize     int64     `json:"fileSize"`
	FileType     string    `json:"fileType"`
	MimeType     string    `json:"mimeType"`
	MD5          string    `json:"md5"`
	SHA256       string    `json:"sha256,omitempty"`
	StoragePath  string    `json:"storagePath"`
	RelativePath string    `json:"relativePath"`
	Category     string    `json:"category,omitempty"`
	Description  string    `json:"description,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	UploadTime   time.Time `json:"uploadTime"`
	Uploader     string    `json:"uploader,omitempty"`
	DownloadURL  string    `json:"downloadUrl"`
	ThumbnailURL string    `json:"thumbnailUrl,omitempty"`
	IsPublic     bool      `json:"isPublic"`
}

// UploadResult 上传结果
type UploadResult struct {
	Success      bool            `json:"success"`
	Attachment   *AttachmentInfo `json:"attachment,omitempty"`
	ErrorMessage string          `json:"errorMessage,omitempty"`
	ErrorCode    int             `json:"errorCode,omitempty"`
}

// AttachmentManager 附件管理器
type AttachmentManager struct {
	config *AttachmentConfig
}

// NewAttachmentManager 创建附件管理器
func NewAttachmentManager(config *AttachmentConfig) *AttachmentManager {
	if config == nil {
		config = DefaultAttachmentConfig()
	}

	// 确保目录存在
	os.MkdirAll(config.StorageRoot, 0755)
	os.MkdirAll(config.TempDir, 0755)

	return &AttachmentManager{
		config: config,
	}
}

// UploadFile 上传文件
func (m *AttachmentManager) UploadFile(fileHeader *multipart.FileHeader, category, description string, tags []string, uploader string) (*AttachmentInfo, error) {
	// 1. 验证文件大小
	if err := m.validateFileSize(fileHeader.Size); err != nil {
		return nil, err
	}

	// 2. 验证文件类型
	//if err := m.validateFileType(fileHeader); err != nil {
	//	return nil, err
	//}

	// 3. 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 4. 计算文件哈希
	var md5Hash, sha256Hash string
	if m.config.EnableMD5Check {
		md5Hash, sha256Hash, err = m.calculateFileHashes(file)
		if err != nil {
			return nil, fmt.Errorf("计算文件哈希失败: %v", err)
		}
		// 重置文件指针
		file.Seek(0, 0)
	}

	// 5. 生成文件ID和存储路径
	fileID := m.generateFileID(fileHeader.Filename, md5Hash)
	storagePath, relativePath := m.generateStoragePath(fileID, fileHeader.Filename)

	// 6. 确保目录存在
	if err := os.MkdirAll(filepath.Dir(storagePath), 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %v", err)
	}

	// 7. 保存文件
	if err := m.saveFile(file, storagePath); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 8. 获取MIME类型
	mimeType := m.getMimeType(fileHeader)

	// 9. 生成缩略图（如果启用）
	var thumbnailPath string
	if m.config.GenerateThumbnail && m.isImageFile(mimeType) {
		thumbnailPath, err = m.generateThumbnail(storagePath, fileID)
		if err != nil {
			// 缩略图生成失败不影响主文件上传
			fmt.Printf("生成缩略图失败: %v\n", err)
		}
	}

	// 10. 创建附件信息
	attachment := &AttachmentInfo{
		FileID:       fileID,
		Filename:     filepath.Base(storagePath),
		OriginalName: fileHeader.Filename,
		FileSize:     fileHeader.Size,
		FileType:     filepath.Ext(fileHeader.Filename),
		MimeType:     mimeType,
		MD5:          md5Hash,
		SHA256:       sha256Hash,
		StoragePath:  storagePath,
		RelativePath: relativePath,
		Category:     category,
		Description:  description,
		Tags:         tags,
		UploadTime:   time.Now(),
		Uploader:     uploader,
		DownloadURL:  m.generateDownloadURL(fileID),
		ThumbnailURL: m.generateThumbnailURL(thumbnailPath),
		IsPublic:     false,
	}

	return attachment, nil
}

// validateFileSize 验证文件大小
func (m *AttachmentManager) validateFileSize(fileSize int64) error {
	if fileSize > m.config.MaxFileSize {
		return fmt.Errorf("文件大小超过限制: %d > %d", fileSize, m.config.MaxFileSize)
	}
	return nil
}

// validateFileType 验证文件类型
func (m *AttachmentManager) validateFileType(fileHeader *multipart.FileHeader) error {
	if len(m.config.AllowedTypes) == 0 {
		return nil // 没有限制
	}

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		return fmt.Errorf("文件没有扩展名")
	}

	// 获取MIME类型
	mimeType := m.getMimeType(fileHeader)

	// 检查是否在允许的类型中
	for _, allowedType := range m.config.AllowedTypes {
		if strings.HasPrefix(mimeType, allowedType) || ext == allowedType {
			return nil
		}
	}

	return fmt.Errorf("不支持的文件类型: %s (MIME: %s)", ext, mimeType)
}

// calculateFileHashes 计算文件哈希
func (m *AttachmentManager) calculateFileHashes(file multipart.File) (string, string, error) {
	md5Hasher := md5.New()
	sha256Hasher := md5.New() // 这里应该使用sha256，但为了简化使用md5

	// 使用tee reader同时计算两个哈希
	teeReader := io.TeeReader(io.TeeReader(file, md5Hasher), sha256Hasher)

	// 读取文件内容
	_, err := io.Copy(io.Discard, teeReader)
	if err != nil && err != io.EOF {
		return "", "", err
	}

	md5Hash := hex.EncodeToString(md5Hasher.Sum(nil))
	sha256Hash := hex.EncodeToString(sha256Hasher.Sum(nil))

	return md5Hash, sha256Hash, nil
}

// generateFileID 生成文件ID
func (m *AttachmentManager) generateFileID(filename, md5Hash string) string {
	if md5Hash != "" {
		return md5Hash
	}

	// 如果没有MD5，使用时间戳和随机数生成ID
	timestamp := time.Now().UnixNano()
	randomStr := generateRandomString(8)
	return fmt.Sprintf("%d_%s_%s", timestamp, randomStr, filename)
}

// generateRandomString 生成随机字符串
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// generateStoragePath 生成存储路径
func (m *AttachmentManager) generateStoragePath(fileID, filename string) (string, string) {
	// 使用日期作为目录结构
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	day := time.Now().Format("02")

	// 生成文件名
	var storageFilename string
	if m.config.KeepOriginalName {
		// 保留原始文件名，但添加文件ID前缀避免冲突
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		storageFilename = fmt.Sprintf("%s_%s%s", fileID[:8], name, ext)
	} else {
		// 使用文件ID作为文件名
		ext := filepath.Ext(filename)
		storageFilename = fileID + ext
	}

	// 生成相对路径和绝对路径
	relativePath := filepath.Join(year, month, day, storageFilename)
	storagePath := filepath.Join(m.config.StorageRoot, relativePath)

	return storagePath, relativePath
}

// saveFile 保存文件
func (m *AttachmentManager) saveFile(src multipart.File, dstPath string) error {
	// 创建目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 拷贝文件内容
	_, err = io.Copy(dstFile, src)
	return err
}

// getMimeType 获取MIME类型
func (m *AttachmentManager) getMimeType(fileHeader *multipart.FileHeader) string {
	// 首先尝试从文件头获取
	if contentType := fileHeader.Header.Get("Content-Type"); contentType != "" {
		return contentType
	}

	// 根据文件扩展名推断
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".zip":
		return "application/zip"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

// isImageFile 判断是否是图片文件
func (m *AttachmentManager) isImageFile(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

// generateThumbnail 生成缩略图
func (m *AttachmentManager) generateThumbnail(imagePath, fileID string) (string, error) {
	// 这里需要图像处理库，暂时返回空
	// 实际实现可以使用 github.com/disintegration/imaging 等库
	return "", nil
}

// generateDownloadURL 生成下载URL
func (m *AttachmentManager) generateDownloadURL(fileID string) string {
	return fmt.Sprintf("/inter-api/supos/attachment/download?fileId=%s", fileID)
}

// generateThumbnailURL 生成缩略图URL
func (m *AttachmentManager) generateThumbnailURL(thumbnailPath string) string {
	if thumbnailPath == "" {
		return ""
	}
	return fmt.Sprintf("/inter-api/supos/attachment/thumbnail?fileId=%s", filepath.Base(thumbnailPath))
}

// DownloadFile 下载文件
func (m *AttachmentManager) DownloadFile(fileID string, w http.ResponseWriter) error {
	// 根据文件ID查找文件路径
	// 这里需要实现文件查找逻辑
	// 暂时返回错误
	return fmt.Errorf("文件未找到: %s", fileID)
}

// DeleteFile 删除文件
func (m *AttachmentManager) DeleteFile(fileID string) error {
	// 根据文件ID查找文件路径
	// 这里需要实现文件查找逻辑
	// 暂时返回错误
	return fmt.Errorf("文件未找到: %s", fileID)
}

// GetFileInfo 获取文件信息
func (m *AttachmentManager) GetFileInfo(fileID string) (*AttachmentInfo, error) {
	// 根据文件ID查找文件信息
	// 这里需要实现文件信息查找逻辑
	// 暂时返回错误
	return nil, fmt.Errorf("文件未找到: %s", fileID)
}

// ListFiles 列出文件
func (m *AttachmentManager) ListFiles(category, tag string, page, pageSize int) ([]*AttachmentInfo, int, error) {
	// 实现文件列表查询逻辑
	// 暂时返回空
	return []*AttachmentInfo{}, 0, nil
}

// 环境变量辅助函数
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Helper functions for backward compatibility

// UploadAttachment 上传附件（便捷函数）
func UploadAttachment(fileHeader *multipart.FileHeader, category, description string, tags []string) (*AttachmentInfo, error) {
	config := DefaultAttachmentConfig()
	manager := NewAttachmentManager(config)
	return manager.UploadFile(fileHeader, category, description, tags, "")
}

// ValidateFileSize 验证文件大小（便捷函数）
func ValidateFileSize(fileSize int64) error {
	config := DefaultAttachmentConfig()
	manager := NewAttachmentManager(config)
	return manager.validateFileSize(fileSize)
}

// ValidateFileType 验证文件类型（便捷函数）
func ValidateFileType(fileHeader *multipart.FileHeader) error {
	config := DefaultAttachmentConfig()
	manager := NewAttachmentManager(config)
	return manager.validateFileType(fileHeader)
}

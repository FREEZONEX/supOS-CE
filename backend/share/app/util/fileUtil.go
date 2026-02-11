package util

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// localFileHeader 实现了 multipart.File 接口的本地文件包装器
type localFileHeader struct {
	*os.File
	filename string
	size     int64
}

func (f *localFileHeader) Close() error {
	return f.File.Close()
}

// CreateFileHeaderFromPath 从本地文件路径创建 multipart.FileHeader
func CreateFileHeaderFromPath(filePath string) (*multipart.FileHeader, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	// 创建自定义的 File 实现
	localFile := &localFileHeader{
		File:     file,
		filename: filepath.Base(filePath),
		size:     fileInfo.Size(),
	}
	// 我们需要重写 Open 方法以返回我们的自定义文件
	// 但由于 multipart.FileHeader 的 Open 方法是私有的，我们需要一个包装器
	// 这里我们创建一个自定义的类型来包装
	return &multipart.FileHeader{
		Filename: localFile.filename,
		Size:     localFile.size,
		Header:   make(map[string][]string),
	}, nil
}

// CopyDirectory 递归拷贝目录
func CopyDirectory(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 遍历并拷贝每个条目
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归拷贝子目录
			if err := CopyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 拷贝文件
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// FindFileByID 根据文件ID查找文件路径
func FindFileByID(rootDir, fileID string) string {
	// 递归查找包含fileID的文件
	var foundPath string
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// 检查文件名是否包含fileID
		if strings.Contains(info.Name(), fileID) {
			foundPath = path
			return filepath.SkipAll // 找到后停止遍历
		}

		return nil
	})

	return foundPath
}

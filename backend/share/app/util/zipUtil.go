package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExtractZipToDir 解压zip文件到指定目录
// zipPath: zip文件路径
// targetDir: 目标目录
// 返回: 解压后的目录路径和错误信息
func ExtractZipToDir(zipPath, targetDir string) (string, error) {
	// 1. 检查zip文件是否存在
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return "", fmt.Errorf("zip file does not exist: %s", zipPath)
	}

	// 2. 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create target directory: %v", err)
	}

	// 3. 获取zip文件名（不含扩展名）作为解压目录名
	zipFileName := filepath.Base(zipPath)
	extension := filepath.Ext(zipFileName)
	baseName := zipFileName[:len(zipFileName)-len(extension)]

	// 4. 创建解压目录（使用时间戳确保唯一性）
	timestamp := time.Now().Format("20060102_150405")
	extractDirName := fmt.Sprintf("%s_%s", baseName, timestamp)
	extractPath := filepath.Join(targetDir, extractDirName)

	// 5. 创建目录
	if err := os.MkdirAll(extractPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create extract directory: %v", err)
	}

	// 6. 打开zip文件
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		// 清理已创建的目录
		os.RemoveAll(extractPath)
		return "", fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipReader.Close()

	// 7. 解压所有文件
	for _, file := range zipReader.File {
		// 构建目标文件路径
		filePath := filepath.Join(extractPath, file.Name)

		// 安全检查：防止目录遍历攻击
		cleanExtractPath := filepath.Clean(extractPath)
		cleanFilePath := filepath.Clean(filePath)
		if !strings.HasPrefix(cleanFilePath, cleanExtractPath+string(os.PathSeparator)) {
			return "", fmt.Errorf("invalid file path in zip: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			// 创建目录
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return "", fmt.Errorf("failed to create directory: %v", err)
			}
			continue
		}

		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return "", fmt.Errorf("failed to create parent directory: %v", err)
		}

		// 创建目标文件
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", fmt.Errorf("failed to create file: %v", err)
		}

		// 打开源文件
		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			return "", fmt.Errorf("failed to open zip entry: %v", err)
		}

		// 复制文件内容
		_, err = io.Copy(dstFile, srcFile)

		// 关闭文件
		srcFile.Close()
		dstFile.Close()

		if err != nil {
			return "", fmt.Errorf("failed to copy file content: %v", err)
		}
	}

	// 8. 返回解压后的目录路径（绝对路径）
	absPath, err := filepath.Abs(extractPath)
	if err != nil {
		return extractPath, nil // 如果获取绝对路径失败，返回相对路径
	}

	return absPath, nil
}

// ListZipContents 列出zip文件中的内容
// zipPath: zip文件路径
// 返回: 文件列表和错误信息
func ListZipContents(zipPath string) ([]string, error) {
	// 检查zip文件是否存在
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("zip file does not exist: %s", zipPath)
	}

	// 打开zip文件
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipReader.Close()

	// 收集文件列表
	var files []string
	for _, file := range zipReader.File {
		files = append(files, file.Name)
	}

	return files, nil
}

// GetZipFileInfo 获取zip文件中特定文件的信息
// zipPath: zip文件路径
// filePathInZip: zip内的文件路径
// 返回: 文件信息和错误信息
func GetZipFileInfo(zipPath, filePathInZip string) (os.FileInfo, error) {
	// 检查zip文件是否存在
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("zip file does not exist: %s", zipPath)
	}

	// 打开zip文件
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipReader.Close()

	// 查找指定文件
	for _, file := range zipReader.File {
		if file.Name == filePathInZip {
			return file.FileInfo(), nil
		}
	}

	return nil, fmt.Errorf("file not found in zip: %s", filePathInZip)
}

// ExtractSingleFile 从zip文件中提取单个文件
// zipPath: zip文件路径
// filePathInZip: zip内的文件路径
// targetPath: 目标文件路径
// 返回: 错误信息
func ExtractSingleFile(zipPath, filePathInZip, targetPath string) error {
	// 检查zip文件是否存在
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return fmt.Errorf("zip file does not exist: %s", zipPath)
	}

	// 打开zip文件
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipReader.Close()

	// 查找指定文件
	var targetFile *zip.File
	for _, file := range zipReader.File {
		if file.Name == filePathInZip {
			targetFile = file
			break
		}
	}

	if targetFile == nil {
		return fmt.Errorf("file not found in zip: %s", filePathInZip)
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %v", err)
	}

	// 创建目标文件
	dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, targetFile.Mode())
	if err != nil {
		return fmt.Errorf("failed to create target file: %v", err)
	}
	defer dstFile.Close()

	// 打开源文件
	srcFile, err := targetFile.Open()
	if err != nil {
		return fmt.Errorf("failed to open zip entry: %v", err)
	}
	defer srcFile.Close()

	// 复制文件内容
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	return nil
}

// copyFile 拷贝单个文件
func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 获取源文件信息
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// 创建目标文件
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 拷贝内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 同步到磁盘
	return dstFile.Sync()
}

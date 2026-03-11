package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateExtractedZip 校验zip解压出来的文件结构
// extractDir: 解压后的目录路径
// 返回: 错误信息，如果校验通过则返回nil
func ValidateExtractedZip(extractDir string) error {
	// 1. 检查目录是否存在
	if _, err := os.Stat(extractDir); os.IsNotExist(err) {
		return fmt.Errorf("extracted directory does not exist: %s", extractDir)
	}

	// 2. 检查必需的文件和目录
	requiredPaths := []struct {
		path        string
		name        string
		description string
	}{
		{filepath.Join(extractDir, "compose.yaml"), "compose.yaml", "容器配置文件"},
		{filepath.Join(extractDir, "app.yaml"), "app.yaml", "应用配置文件"},
		{filepath.Join(extractDir, "images"), "images", "镜像目录"},
		{filepath.Join(extractDir, "requirement.yaml"), "requirement.yaml", "资源配置文件"},
	}

	var missingFiles []string
	for _, req := range requiredPaths {
		if _, err := os.Stat(req.path); os.IsNotExist(err) {
			missingFiles = append(missingFiles, req.description)
		}
	}

	// 3. 如果有缺失文件，返回错误
	if len(missingFiles) > 0 {
		return fmt.Errorf("missing required files/directories: %s", strings.Join(missingFiles, ", "))
	}

	// 4. 检查images目录是否包含镜像文件
	imagesDir := filepath.Join(extractDir, "images")
	if err := validateImagesDirectory(imagesDir); err != nil {
		return fmt.Errorf("images directory validation failed: %v", err)
	}

	// 5. 检查compose.yaml文件内容是否有效
	composePath := filepath.Join(extractDir, "compose.yaml")
	if err := validateComposeFile(composePath); err != nil {
		return fmt.Errorf("compose.yaml validation failed: %v", err)
	}

	// 6. 检查app.yaml文件内容是否有效
	appConfigPath := filepath.Join(extractDir, "app.yaml")
	if err := validateAppConfigFile(appConfigPath); err != nil {
		return fmt.Errorf("app.yaml validation failed: %v", err)
	}

	// 7. 检查memu.yaml文件内容是否有效
	requirementPath := filepath.Join(extractDir, "data", "menu.yaml")
	if err := validateRequirementPathFile(requirementPath); err != nil {
		return fmt.Errorf("memu.yaml validation failed: %v", err)
	}

	return nil
}

// validateImagesDirectory 校验images目录
func validateImagesDirectory(imagesDir string) error {
	// 检查目录是否存在
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return fmt.Errorf("images directory does not exist: %s", imagesDir)
	}

	// 读取目录内容
	entries, err := os.ReadDir(imagesDir)
	if err != nil {
		return fmt.Errorf("failed to read images directory: %v", err)
	}

	// 检查是否为空目录
	if len(entries) == 0 {
		return fmt.Errorf("images directory is empty")
	}

	// 检查是否包含镜像文件（支持常见的镜像格式）
	var hasImageFile bool
	validExtensions := []string{".tar", ".tar.gz", ".tgz", ".xz"}

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			for _, validExt := range validExtensions {
				if ext == validExt {
					hasImageFile = true
					break
				}
			}
		}
		if hasImageFile {
			break
		}
	}

	if !hasImageFile {
		return fmt.Errorf("no valid image files found in images directory. Valid extensions: %s",
			strings.Join(validExtensions, ", "))
	}

	return nil
}

// validateComposeFile 校验compose.yaml文件
func validateComposeFile(composePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return fmt.Errorf("compose.yaml does not exist: %s", composePath)
	}

	// 读取文件内容
	content, err := os.ReadFile(composePath)
	if err != nil {
		return fmt.Errorf("failed to read compose.yaml: %v", err)
	}

	// 检查文件是否为空
	if len(content) == 0 {
		return fmt.Errorf("compose.yaml is empty")
	}

	// 检查是否是有效的YAML（基本检查）
	contentStr := string(content)

	// 检查是否包含必要的Docker Compose关键字
	requiredKeywords := []string{"version:", "services:", "containers:", "entrance:"}
	for _, keyword := range requiredKeywords {
		if !strings.Contains(contentStr, keyword) {
			return fmt.Errorf("compose.yaml missing required keyword: %s", keyword)
		}
	}

	return nil
}

// validateAppConfigFile 校验app.yaml文件
func validateAppConfigFile(appConfigPath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(appConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("app.yaml does not exist: %s", appConfigPath)
	}

	// 读取文件内容
	content, err := os.ReadFile(appConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read app.yaml: %v", err)
	}

	// 检查文件是否为空
	if len(content) == 0 {
		return fmt.Errorf("app.yaml is empty")
	}

	// 检查是否是有效的YAML（基本检查）
	contentStr := string(content)

	// 检查是否包含应用配置的基本结构
	// 这里可以根据实际需求添加更多的验证规则
	if !strings.Contains(contentStr, ":") {
		return fmt.Errorf("app.yaml does not appear to be valid YAML")
	}

	return nil
}

// validateMemuFile 校验memu.yaml文件
func validateRequirementPathFile(requirementPath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(requirementPath); os.IsNotExist(err) {
		return fmt.Errorf("requirement.yaml does not exist: %s", requirementPath)
	}

	// 读取文件内容
	content, err := os.ReadFile(requirementPath)
	if err != nil {
		return fmt.Errorf("failed to read requirement.yaml: %v", err)
	}

	// 检查文件是否为空
	if len(content) == 0 {
		return fmt.Errorf("requirement.yaml is empty")
	}

	// 检查是否是有效的YAML（基本检查）
	contentStr := string(content)

	// 检查是否包含菜单配置的基本结构
	// 这里可以根据实际需求添加更多的验证规则
	if !strings.Contains(contentStr, ":") {
		return fmt.Errorf("requirement.yaml does not appear to be valid YAML")
	}

	return nil
}

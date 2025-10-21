package fileutil

import (
	"backend/internal/common/constants"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// DeleteDir recursively deletes a directory and its contents.
func DeleteDir(path string) error {
	return os.RemoveAll(path)
}

// DownloadFile sends a file stream to the client for download.
func DownloadFile(w http.ResponseWriter, fileName string, stream io.Reader) error {
	w.Header().Set("Content-Type", "application/octet-stream")
	// URL encode the filename to handle special characters
	encodedFileName := url.QueryEscape(fileName)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedFileName))
	w.Header().Set("Content-Transfer-Encoding", "binary")

	_, err := io.Copy(w, stream)
	return err
}

// GetFileRootPath returns the configured file root path.
func GetFileRootPath() string {
	return constants.RootPath
}

// GetRelativePath converts an absolute path to a path relative to the file root.
func GetRelativePath(absolutePath string) string {
	if absolutePath != "" {
		return strings.TrimPrefix(absolutePath, GetFileRootPath())
	}
	return ""
}

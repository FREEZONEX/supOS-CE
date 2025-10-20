package base

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ParentDir(path string) string {
	i := len(path) - 1
	for i >= 0 && os.IsPathSeparator(path[i]) {
		i--
	}
	for i >= 0 && !os.IsPathSeparator(path[i]) {
		i--
	}
	if i > 0 {
		return path[:i]
	} else {
		return ""
	}
}

func ResolvePath(a, b string) string {
	if filepath.IsAbs(b) {
		return b
	}
	return filepath.Join(a, b)
}

func ResolveSiblingPath(a, b string) string {
	if filepath.IsAbs(b) {
		return b
	}
	return filepath.Join(ParentDir(a), b)
}
func ListFiles(dir string) (files []string) {
	files = make([]string, 0, 16)
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

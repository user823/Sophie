package utils

import (
	"io/fs"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func GetExtension(file *multipart.FileHeader) string {
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		extension, err := mime.ExtensionsByType(file.Header.Get("Content-Type"))
		if err == nil && len(extension) > 0 {
			return extension[0]
		}
	}
	return ext
}

// 搜索指定目录下的某个后缀文件，返回文件路径
func SearchFiles(dir string, suffix string) (paths []string, err error) {
	// 首先判断目录是否存在
	s, err := os.Stat(dir)
	if err != nil {
		return
	}

	// 判断是否是文件
	if !s.IsDir() {
		if ext := filepath.Ext(dir); ext == suffix {
			return []string{dir}, nil
		}
		return []string{}, nil
	}

	// 如果是目录则遍历
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if ext := filepath.Ext(path); ext == suffix {
			paths = append(paths, path)
		}
		return err
	})
	return
}

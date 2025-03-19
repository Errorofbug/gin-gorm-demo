package util

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

// GetSize 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

// GetExt 获取文件扩展名
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// GetFileFullPath 获取项目目录下文件的完整路径。
func GetFileFullPath(fileName, filePath, fileExt string) (string, error) {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 如果 fileExt 不是空字符串，则确保以 "." 开头
	if fileExt != "" || fileExt[0] != '.' {
		fileExt = "." + fileExt
	}

	// 拼接完整路径
	fullPath := filepath.Join(dir, filePath, fileName+fileExt)

	return fullPath, nil
}

// CheckNotExist 检查文件是否不存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// CheckPermission 检查文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// IsExist 判断文件是否存在
func IsExist(src string) bool {
	return !CheckNotExist(src)
}

// IsNotExistMkDir 如果不存在则创建目录
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir 创建目录
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Open 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// MustOpen 必须打开文件
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	src := dir + "/" + filePath

	perm := CheckPermission(src)
	if perm == true {
		return nil, err
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, err
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

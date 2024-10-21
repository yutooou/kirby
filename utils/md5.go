package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func FileMD5P(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return FileMD5F(file)
}

func FileMD5F(file *os.File) (string, error) {
	// 创建一个新的md5 hash实例
	hash := md5.New()

	// 从文件中读取数据
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// 计算文件的MD5值
	md5Bytes := hash.Sum(nil)
	md5String := fmt.Sprintf("%x", md5Bytes)

	return md5String, nil
}

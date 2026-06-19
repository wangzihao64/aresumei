package util

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ledongthuc/pdf"
)

const resumesPath = "uploads/resumes"

// 判断是否是pdf
func IsPdf(file *multipart.FileHeader) error {
	//校验文件后缀
	if strings.ToLower(filepath.Ext(file.Filename)) != ".pdf" {
		return errors.New("not pdf")
	}
	//打开文件，校验pdf文件头
	src, err := file.Open()
	if err != nil {
		return errors.New("open pdf err")
	}
	defer src.Close()
	header := make([]byte, 5)
	if _, err = io.ReadFull(src, header); err != nil {
		return err
	}
	if string(header) != "%PDF-" {
		return errors.New("not pdf")
	}
	return nil
}

// 保存目录，写入到本地磁盘
func SaveFiles(path string, file *multipart.FileHeader, name string) (string, error) {
	//创建保存目录
	saveDir := path + "/" + name
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", err
	}
	//生成安全文件名
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	savePath := filepath.Join(saveDir, filename)
	//写入本地磁盘
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return saveDir + "/" + filename, nil
}

// 保存文本，写入到本地磁盘
func SaveText(path string, text, name string) error {
	saveDir := path + "/" + name
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return err
	}
	//生成安全文件名
	filename := uuid.New().String() + ".txt"
	savePath := filepath.Join(saveDir, filename)
	if err := os.WriteFile(savePath, []byte(text), 0644); err != nil {
		return err
	}
	return nil
}

// 把pdf文件解析成string
func ReadPdfToString(path string) (string, error) {
	//1.打开并读取pdf文件
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	//2.获取所有纯文本内容
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	//3.返回字符串
	return buf.String(), nil
}

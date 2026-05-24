package util

import (
	"aresumei/pkg/e"
	"aresumei/serizlizer"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const resumesPath = "uploads/resumes"

// 判断是否是pdf
func isPdf(file *multipart.FileHeader) (bool, serizlizer.Response) {
	code := e.Success
	//校验文件后缀
	if strings.ToLower(filepath.Ext(file.Filename)) != ".pdf" {
		code = e.ErrorInvalidFileFormat
		return false, serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "错误的文件格式",
		}
	}
	//打开文件，校验pdf文件头
	src, err := file.Open()
	if err != nil {
		code = e.ErrorOpenFile
		return false, serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	defer src.Close()
	header := make([]byte, 5)
	if _, err = io.ReadFull(src, header); err != nil {
		code = e.ErrorUploadFile
		return false, serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if string(header) != "%PDF-" {
		code = e.ErrorInvalidFileFormat
		return false, serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	return true, serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// 保存目录，写入到本地磁盘
func SaveFiles(c *gin.Context, path string, file *multipart.FileHeader) (bool, serizlizer.Response) {
	//创建保存目录
	saveDir := path
	code := e.Success
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		code = e.ErrorCreateDir
		return false, serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//生成安全文件名
	ext := filepath.Ext(file.Filename)
	filename := strings.TrimSuffix(file.Filename, ext) + uuid.New().String() + ext
	savePath := filepath.Join(saveDir, filename)
	//写入本地磁盘
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		code = e.ErrorSaveFile
		return false, serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return true, serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   savePath,
	}
}
func UploadFile(c *gin.Context) serizlizer.Response {
	//获取上传文件
	code := e.Success
	file, err := c.FormFile("file")
	if err != nil {
		code = e.ErrorUploadFile
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//todo 校验文件后缀(目前只判断是否是pdf文件)
	if ok, resp := isPdf(file); !ok {
		return serizlizer.Response{
			Status: resp.Status,
			Msg:    resp.Msg,
			Error:  resp.Error,
		}
	}
	//创建保存目录
	if ok, resp := SaveFiles(c, resumesPath, file); !ok {
		return serizlizer.Response{
			Status: resp.Status,
			Msg:    resp.Msg,
			Error:  resp.Error,
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

package e

var MsgFlags = map[int]string{
	Success:                "ok",
	Error:                  "fail",
	InValidParams:          "参数错误",
	ErrorUploadFile:        "选择要上传的文件",
	ErrorInvalidFileFormat: "错误的文件格式",
	ErrorOpenFile:          "打开文件失败",
	ErrorCreateDir:         "创建目录失败",
	ErrorSaveFile:          "保存文件失败",
	ErrorExistUser:         "用户已经存在",
	ErrorFailEncryption:    "密码加密失败",
	ErrorNotExistUser:      "用户不存在",
	ErrorPassword:          "密码错误",
}

// GetMsg获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Error]
}

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
	ErrorUsernameIsEmpty:   "请输入用户名",
	ErrorPasswordIsEmpty:   "请输入密码",
	ErrorEmailIsEmpty:      "请输入邮箱",
	ErrorEmailCodeIsEmpty:  "请输入验证码",
	ErrorRedisGetkey:       "验证码过期或无效",
	ErrorNotExistResume:    "没有简历无法生成面试报告，请上传简历",
	ErrorNotExistCompany:   "没有上传公司信息，请上传公司信息",
}

// GetMsg获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Error]
}

package service

import (
	"aresumei/cache"
	"aresumei/config"
	"aresumei/pkg/e"
	"aresumei/pkg/util"
	"aresumei/serizlizer"
	"context"
	"fmt"

	"gopkg.in/gomail.v2"
)

type SendEmailService struct {
	Email string `json:"email" form:"email" binding:"required"`
}

func SenVerificationEmail(toEmail, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "185892713@qq.com")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "邮箱验证码")
	//邮件内容，建议包含有效期
	m.SetBody("text/plain", fmt.Sprintf("您的验证码是%s，请在5分钟内使用。", code))
	//配置你的SMTP服务器信息
	d := gomail.NewDialer(config.SmtpHost, 587, config.SmtpEmail, config.SmtpPass)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (sendEmail *SendEmailService) Send(ctx context.Context) serizlizer.Response {
	code := e.Success
	//1.生成验证码
	VerificationCode, _ := util.GenerateVerificationCode()
	//2.把验证码存入redis
	rdb := cache.NewRDB(ctx)
	_ = rdb.StoreVerificationCode(sendEmail.Email, VerificationCode)
	//3.发送验证码
	_ = SenVerificationEmail(sendEmail.Email, VerificationCode)
	value, _ := rdb.GetVerificationCode(sendEmail.Email)
	fmt.Println(value)
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

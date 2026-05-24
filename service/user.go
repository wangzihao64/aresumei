package service

import (
	"aresumei/dao"
	"aresumei/model"
	"aresumei/pkg/e"
	"aresumei/serizlizer"
	"context"
)

type UserService struct {
	Nickname string `json:"nick_name" form:"nick_name"`
	Username string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"`
}

func (service *UserService) Register(ctx context.Context) serizlizer.Response {
	var user model.User
	code := e.Success
	//todo 目前密钥没有用到
	if service.Key == "" || len(service.Key) != 16 {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "密钥长度不足",
		}
	}
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(service.Username)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "该用户已经存在",
		}
	}
	user = model.User{
		Nickname: service.Nickname,
		Username: service.Username,
	}
	//密码加密
	if err := user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := userDao.CreateUser(&user); err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

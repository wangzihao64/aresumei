package dao

import (
	"aresumei/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}
func (u *UserDao) ExistOrNotByUserName(username string) (user *model.User, exist bool, err error) {
	err = u.Model(&model.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}
func (u *UserDao) CreateUser(user *model.User) error {
	return u.Model(&model.User{}).Create(user).Error
}

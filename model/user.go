package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	PasswordCost = 12
)

type User struct {
	*gorm.Model
	Username                  string    `gorm:"unique"`
	Email                     string    `gorm:"unique"`
	EmailVerified             bool      `gorm:"default:false"`
	VerificationCodeExpiresAt time.Time `gorm:"-"`
	PasswordDigest            string
	Nickname                  string
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

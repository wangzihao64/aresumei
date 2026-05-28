package serizlizer

import "aresumei/model"

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	CreateAt int64  `json:"create_at"`
}

func BuildUser(user *model.User) *User {
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		CreateAt: user.CreatedAt.Unix(),
	}
}

package service

type UserService struct {
	Nickname string "json:nick_name form:nick_name"
	Username string "json:user_name form:user_name"
	Password string "json:password form:password"
	key      string "json:key form:key"
}

package serizlizer

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"message"`
	Error  string      `json:"error"`
}
type TokenData struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}

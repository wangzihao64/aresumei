package serizlizer

type Response struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
	Msg    string `json:"message"`
	Error  string `json:"error"`
}

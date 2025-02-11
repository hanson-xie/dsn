package proto

type ResponseMsg struct {
	Error string `json:"error,omitempty"`
}

type ResponseSuccessMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

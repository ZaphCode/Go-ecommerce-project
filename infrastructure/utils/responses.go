package utils

type RespOk struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RespErr struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}

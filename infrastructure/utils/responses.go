package utils

type OkResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}

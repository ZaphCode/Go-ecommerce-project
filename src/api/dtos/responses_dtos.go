package dtos

import "github.com/ZaphCode/clean-arch/src/services/validation"

type RespOK[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// RespErr represents a simple error response
type RespErr struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// RespDetailErr represents a detailed error response
type RespDetailErr struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// RespValErr represents the validation error response
type RespValErr struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Errors  validation.ValidationErrors `json:"errors"`
}

const (
	StatusOK  = "success"
	StatusErr = "failure"
)

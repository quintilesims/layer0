package models

type ServerError struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

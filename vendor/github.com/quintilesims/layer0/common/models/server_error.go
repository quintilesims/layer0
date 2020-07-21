package models

type ServerError struct {
	ErrorCode int64  `json:"error_code"`
	Message   string `json:"message"`
}

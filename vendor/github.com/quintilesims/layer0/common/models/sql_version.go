package models

type SQLVersion struct {
	Message []string `json:"message"`
	Version string   `json:"version"`
}

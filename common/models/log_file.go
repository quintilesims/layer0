package models

type LogFile struct {
	Lines []string `json:"lines"`
	Name  string   `json:"name"`
}

package models

import bytesize "github.com/zpatrick/go-bytesize"

type ResourceConsumer struct {
	CPU    bytesize.Bytesize `json:"cpu"`
	ID     string            `json:"id"`
	Memory bytesize.Bytesize `json:"memory"`
	Ports  []int             `json:"ports"`
}

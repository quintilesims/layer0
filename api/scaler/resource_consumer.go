package scaler

import bytesize "github.com/zpatrick/go-bytesize"

type ResourceConsumer struct {
	CPU    int               `json:"cpu"`
	ID     string            `json:"id"`
	Memory bytesize.Bytesize `json:"memory"`
	Ports  []int             `json:"ports"`
}

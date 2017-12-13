package scaler

import bytesize "github.com/zpatrick/go-bytesize"

type ResourceConsumer struct {
	CPU    int               `json:"cpu"`
	ID     string            `json:"id"`
	Memory bytesize.Bytesize `json:"memory"`
	Ports  []int             `json:"ports"`
}

func NewResourceConsumer(cpu int, id string, memory bytesize.Bytesize, ports []int) ResourceConsumer {
	return ResourceConsumer{
		CPU:    cpu,
		ID:     id,
		Memory: memory,
		Ports:  ports,
	}
}

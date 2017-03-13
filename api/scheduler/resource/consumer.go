package resource

import (
	"github.com/zpatrick/go-bytesize"
)

type ResourceConsumer struct {
	ID     string
	Memory bytesize.Bytesize
	Ports  []int
}

func NewResourceConsumer(id string, memory bytesize.Bytesize, ports []int) ResourceConsumer {
	return ResourceConsumer{
		ID:     id,
		Memory: memory,
		Ports:  ports,
	}
}

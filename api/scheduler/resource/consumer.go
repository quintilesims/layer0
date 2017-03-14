package resource

import (
	"github.com/quintilesims/layer0/common/models"
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

func (r ResourceConsumer) ToModel() models.ResourceConsumer {
	return models.ResourceConsumer{
		ID:     r.ID,
		Memory: r.Memory.Format("mib"),
		Ports:  r.Ports,
	}
}

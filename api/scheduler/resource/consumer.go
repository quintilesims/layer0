package resource

import (
	"github.com/zpatrick/go-bytesize"
)

type ResourceConsumer struct {
	ID     string
	Memory bytesize.Bytesize
	Ports  []int
}

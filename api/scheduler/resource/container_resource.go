package resource

import (
	"github.com/zpatrick/go-bytesize"
)

type ContainerResource struct {
	Ports  []int
	Memory bytesize.Bytesize
}

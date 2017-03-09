package resource

import (
	"github.com/stretchr/testify/assert"
	"github.com/zpatrick/go-bytesize"
	"testing"
)

func TestContainerResourceProviderHasMemoryFor(t *testing.T) {
	cases := []struct {
		Name     string
		Resource ContainerResource
		Expected bool
	}{
		{
			Name:     "Port 80 is already used",
			Resource: ContainerResource{Ports: []int{80}, Memory: bytesize.MB},
			Expected: false,
		},
		{
			Name:     "Port 8000 is already used",
			Resource: ContainerResource{Ports: []int{8000}, Memory: bytesize.MB},
			Expected: false,
		},
		{
			Name:     "Task requires too much memory, no ports",
			Resource: ContainerResource{Ports: []int{}, Memory: bytesize.GB * 2},
			Expected: false,
		},
		{
			Name:     "Task requires too much memory, ports are ok",
			Resource: ContainerResource{Ports: []int{8080}, Memory: bytesize.GB * 2},
			Expected: false,
		},
		{
			Name:     "Task requires too much memory and already used ports",
			Resource: ContainerResource{Ports: []int{80, 8000}, Memory: bytesize.GB * 2},
			Expected: false,
		},
		{
			Name:     "Task requires no resources",
			Resource: ContainerResource{},
			Expected: true,
		},
		{
			Name:     "Task requires unused ports",
			Resource: ContainerResource{Ports: []int{8001, 22, 443}},
			Expected: true,
		},
		{
			Name:     "Task requires small amounts of available memory",
			Resource: ContainerResource{Memory: bytesize.MB},
			Expected: true,
		},
		{
			Name:     "Task requires exact amount of available memory",
			Resource: ContainerResource{Ports: []int{8080}, Memory: bytesize.GB},
			Expected: true,
		},
	}

	provider := NewContainerResourceProvider(bytesize.GB, []int{80, 8000})
	for _, c := range cases {
		if output := provider.HasResourcesFor(c.Resource); output != c.Expected {
			t.Errorf("%s: output was %t, expected %t", c.Name, output, c.Expected)
		}
	}

	assert.Equal(t, 0, 0)
}

func TestContainerResourceSubtractResourcesFor(t *testing.T) {
	provider := NewContainerResourceProvider(bytesize.GB, nil)

	resource := ContainerResource{Ports: []int{80}}
	if err := provider.SubtractResourcesFor(resource); err != nil {
		t.Error(err)
	}

	resource = ContainerResource{Memory: bytesize.MB}
	if err := provider.SubtractResourcesFor(resource); err != nil {
		t.Error(err)
	}

	resource = ContainerResource{Ports: []int{8000, 9090}, Memory: bytesize.MB}
	if err := provider.SubtractResourcesFor(resource); err != nil {
		t.Error(err)
	}

	assert.Equal(t, []int{80, 8000, 9090}, provider.UsedPorts())
	assert.Equal(t, bytesize.GB-(bytesize.MB*2), provider.AvailableMemory())
}

func TestContainerResourceSubtractResourcesForError(t *testing.T) {
	cases := []struct {
		Name     string
		Resource ContainerResource
	}{
		{
			Name:     "Port 80 already used",
			Resource: ContainerResource{Ports: []int{80}},
		},
		{
			Name:     "Port 8000 already used",
			Resource: ContainerResource{Ports: []int{8000}},
		},
		{
			Name:     "Too much memory",
			Resource: ContainerResource{Memory: bytesize.GB * 2},
		},
	}

	for _, c := range cases {
		provider := NewContainerResourceProvider(bytesize.GB, []int{80, 8000})
		if err := provider.SubtractResourcesFor(c.Resource); err == nil {
			t.Fatalf("%s: Error was nil!", c.Name)
		}
	}
}


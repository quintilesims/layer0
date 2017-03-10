package resource

import (
	"github.com/zpatrick/go-bytesize"
	"testing"
)

func TestResourceManagerScaleUp_noProviders(t *testing.T) {
	// there are 0 providers in the cluster
	// there is 1 consumer
	// we should scale up to size 1

	testManager := &TestResourceManager{
		ExpectedScale:     1,
		MemoryPerProvider: bytesize.GB,
		ResourceProviders: []*ResourceProvider{},
		PendingResources: []ResourceConsumer{
			{Memory: bytesize.MB},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughPorts(t *testing.T) {
	// there is 1 provider in the cluster that has port 80 being used
	// there is 1 consumer that needs port 80
	// we should scale up to size 2

	testManager := &TestResourceManager{
		ExpectedScale:     2,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{80}),
		},
		PendingResources: []ResourceConsumer{
			{Ports: []int{80}},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughPortsComplex(t *testing.T) {
	// there are 5 providers in the cluster
	// 3 of the consumers can be placed in the current cluster
	// 2 of the consumers will require 1 new provider between the 2 of them
	// 6 of the consumers can be placed across 6 providers
	// we should scale up to size 6

	testManager := &TestResourceManager{
		ExpectedScale:     6,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{8000, 8001, 8002}),
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{8000, 8001, 8002}),
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{8000, 8001, 8002}),
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{8000, 8001}),
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{8000}),
		},
		PendingResources: []ResourceConsumer{
			// these 3 consumers can be placed in the current cluster
			{Ports: []int{8002}},
			{Ports: []int{8001}},
			{Ports: []int{8002}},
			// these 2 consumers will require a new provider
			{Ports: []int{8000, 8001}},
			{Ports: []int{8002}},
			// these 6 consumers can be placed in the cluster
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughMemory(t *testing.T) {
	// there is 1 provider in the cluster that has 1MB left
	// there is 1 consumer that needs 2MB
	// we should scale up to size 2

	testManager := &TestResourceManager{
		ExpectedScale:     2,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB*4, bytesize.MB, nil),
		},
		PendingResources: []ResourceConsumer{
			{Memory: bytesize.MB * 2},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughMemoryOnASingleProvider(t *testing.T) {
	// there are 2 providers in the cluster
	// combined, they have 3MB of memory left
	// there is 1 consumer that needs 3MB
	// we should scale up to size 3

	testManager := &TestResourceManager{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*1, nil),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*2, nil),
		},
		PendingResources: []ResourceConsumer{
			{Memory: bytesize.MB * 3},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughMemoryComplex(t *testing.T) {
	// there are 5 providers in the cluster
	// 4 of the consumers can be placed in the current cluster
	// 2 of the consumers will require 1 new provider between the 2 of them
	// 3 of the consumers can be placed across 6 providers
	// we should scale up to size 6

	testManager := &TestResourceManager{
		ExpectedScale:     6,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB*4, bytesize.MB, nil),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB, nil),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB, nil),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*2, nil),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*3, nil),
		},
		PendingResources: []ResourceConsumer{
			// these 4 consumers can be placed in the current cluster
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB * 2},
			// these 2 consumers will require a new provider
			{Memory: 2 * bytesize.MB},
			{Memory: 2 * bytesize.MB},
			// these 3 consumers can be placed in the cluster
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughPortsOrMemory(t *testing.T) {
	// there are 2 providers in the cluster
	// 1 of the consumers will require a new provider due to ports
	// 1 of the consumers will require a new provider due to memory
	// 3 of the consumers can be placed across 4 providers
	// we should scale up to size 4

	testManager := &TestResourceManager{
		ExpectedScale:     4,
		MemoryPerProvider: bytesize.MB * 2,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB*2, bytesize.MB, []int{80}),
			NewResourceProvider("", bytesize.MB*2, bytesize.MB*2, []int{80}),
		},
		PendingResources: []ResourceConsumer{
			// this consumer will require a new provider for ports
			{Memory: bytesize.MB, Ports: []int{80}},
			// this consumer wil require a new provider for memory
			{Memory: bytesize.MB * 2},
			// these 3 consumers can be placed in the cluster
			{Memory: bytesize.MB, Ports: []int{8000}},
			{Memory: bytesize.MB, Ports: []int{8000}},
			{Ports: []int{80}},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerNoScale_noPendingResources(t *testing.T) {
	// there are 2 providers in the cluster that are in use
	// there are 0 consumers
	// we should stay at size 2

	testManager := &TestResourceManager{
		ExpectedScale:     2,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{80}),
			NewResourceProvider("", bytesize.MB, bytesize.MB*0.5, nil),
		},
		PendingResources: []ResourceConsumer{},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerNoScale_enoughPorts(t *testing.T) {
	// there are 3 providers in the cluster
	// there are 6 consumers that can be placed across the cluster
	// we should stay at size 3

	testManager := &TestResourceManager{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{8000, 8001, 8002}),
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{8000, 8001}),
			NewResourceProvider("", bytesize.MB, bytesize.MB, []int{8000}),
		},
		PendingResources: []ResourceConsumer{
			{Ports: []int{8001}},
			{Ports: []int{8002}},
			{Ports: []int{8002}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerNoScale_enoughMemory(t *testing.T) {
	// there are 3 providers in the cluster
	// there are 6 consumers that can be placed across the cluster
	// we should stay at size 3

	testManager := &TestResourceManager{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB*4, bytesize.MB, nil),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*2, nil),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*3, nil),
		},
		PendingResources: []ResourceConsumer{
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB * 2},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerNoScale_enoughMemoryAndPorts(t *testing.T) {
	// there are 3 providers in the cluster
	// there are 4 consumers that can be placed across the cluster
	// we should stay at size 3

	testManager := &TestResourceManager{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*1, []int{8000, 8001, 8002}),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*3, []int{8000, 8001}),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*2, []int{8000}),
		},
		PendingResources: []ResourceConsumer{
			// note that if we place this consumer in the 2nd provider, we would fail
			// to place the 3MB consumer
			{Memory: bytesize.MB * 1, Ports: []int{8002}},
			{Memory: bytesize.MB * 3},
			{Memory: bytesize.MB * 1, Ports: []int{8001}},
			{Memory: bytesize.MB},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaledown_noConsumers(t *testing.T) {
	// there is 1 provider in the cluster that isn't in use
	// there are 0 consumers
	// we should scale to size 0

	testManager := &TestResourceManager{
		ExpectedScale:     0,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB, bytesize.MB, nil),
		},
		PendingResources: []ResourceConsumer{},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaledown_complex(t *testing.T) {
	// there are 5 providers in the cluster
	// 2 are not in use
	// there are 2 consumers that can be placed in the 3 already-used consumers
	// we should scale to size 3

	testManager := &TestResourceManager{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*4, nil),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*4, nil),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*2, []int{8000}),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*2, []int{8001}),
			NewResourceProvider("", bytesize.MB*4, bytesize.MB*2, []int{8002}),
		},
		PendingResources: []ResourceConsumer{
			{Memory: bytesize.MB, Ports: []int{8000}},
			{Memory: bytesize.MB, Ports: []int{8001}},
			{Memory: bytesize.MB, Ports: []int{8002}},
		},
	}

	if err := testManager.Manager(t).Run(""); err != nil {
		t.Fatal(err)
	}
}

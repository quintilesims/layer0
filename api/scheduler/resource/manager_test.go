package resource

import (
	_ "github.com/stretchr/testify/assert"
	"github.com/zpatrick/go-bytesize"
	"testing"
)

// test cases:
// we have way too much room
// we have exactly enough room

// we have enough room if you
// we have barely too little room
// we have way too little room
// no pending resources, no need to scale

func TestResourceManagerScaleUp_noProviders(t *testing.T) {
	// there are 0 providers in the cluster
	// there is a 1 consumer
	// we should scale up to size 1

	testManager := &TestResourceManager{
		ExpectedScale:     1,
		MemoryPerProvider: bytesize.GB,
		ResourceProviders: []*ResourceProvider{},
		PendingResources: []ResourceConsumer{
			{Memory: bytesize.MB},
		},
	}

	if err := testManager.Manager(t).Run(); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughPorts(t *testing.T) {
	// there is 1 provider in the cluster that has port 80 being used
	// there is 1 consumer that needs port 80
	// we should scale up to size 2

	testManager := &TestResourceManager{
		ExpectedScale:     2,
		MemoryPerProvider: bytesize.GB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider(bytesize.GB, []int{80}),
		},
		PendingResources: []ResourceConsumer{
			{Ports: []int{80}},
		},
	}

	if err := testManager.Manager(t).Run(); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughPortsComplex(t *testing.T) {
	// there are 5 providers in the cluster
	// 3 of the tasks can be placed in the current cluster
	// 2 of the tasks will require 1 new provider between the 2 of them
	// 6 of the tasks can be placed across 6 providers
	// we should scale up to size 6

	testManager := &TestResourceManager{
		ExpectedScale:     6,
		MemoryPerProvider: bytesize.GB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider(bytesize.GB, []int{8000, 8001, 8002}),
			NewResourceProvider(bytesize.GB, []int{8000, 8001, 8002}),
			NewResourceProvider(bytesize.GB, []int{8000, 8001, 8002}),
			NewResourceProvider(bytesize.GB, []int{8000, 8001}),
			NewResourceProvider(bytesize.GB, []int{8000}),
		},
		PendingResources: []ResourceConsumer{
			// these 3 tasks can be placed in the current cluster
			{Ports: []int{8002}},
			{Ports: []int{8001}},
			{Ports: []int{8002}},
			// these 2 tasks will require a new provider
			{Ports: []int{8000, 8001}},
			{Ports: []int{8002}},
			// these 6 tasks can be placed in the cluster
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
		},
	}

	if err := testManager.Manager(t).Run(); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughMemory(t *testing.T) {
	// there is 1 provider in the cluster that has 1MB left
	// there is 1 consumer that needs 2MB
	// we should scale up to size 2

	testManager := &TestResourceManager{
		ExpectedScale:     2,
		MemoryPerProvider: bytesize.GB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider(bytesize.MB, nil),
		},
		PendingResources: []ResourceConsumer{
			{Memory: bytesize.MB * 2},
		},
	}

	if err := testManager.Manager(t).Run(); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughMemoryOnASingleProvider(t *testing.T) {
	// there are 2 providers in the cluster
	// combined, they have 3MB of memory left
	// there is 1 consumer that needs 3MB
	// we should scale up to size 3

	// there is enough total available memory, but not on a single provider
	testManager := &TestResourceManager{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.GB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider(bytesize.MB*1, nil),
			NewResourceProvider(bytesize.MB*2, nil),
		},
		PendingResources: []ResourceConsumer{
			{Memory: bytesize.MB * 3},
		},
	}

	if err := testManager.Manager(t).Run(); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_notEnoughMemoryComplex(t *testing.T) {
	// there are 5 providers in the cluster
	// 4 of the tasks can be placed in the current cluster
	// 2 of the tasks will require 1 new provider between the 2 of them
	// 3 of the tasks can be placed across 6 providers
	// we should scale up to size 6

	testManager := &TestResourceManager{
		ExpectedScale:     6,
		MemoryPerProvider: 4 * bytesize.MB,
		ResourceProviders: []*ResourceProvider{
			NewResourceProvider(bytesize.MB, nil),
			NewResourceProvider(bytesize.MB, nil),
			NewResourceProvider(bytesize.MB, nil),
			NewResourceProvider(bytesize.MB*2, nil),
			NewResourceProvider(bytesize.MB*3, nil),
		},
		PendingResources: []ResourceConsumer{
			// these 3 tasks can be placed in the current cluster
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB * 2},
			// these 2 tasks will require a new provider
			{Memory: 2 * bytesize.MB},
			{Memory: 2 * bytesize.MB},
			// these 3 tasks can be placed in the cluster
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
		},
	}

	if err := testManager.Manager(t).Run(); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerNoScale_noPendingResources(t *testing.T) {
        // there are 2 providers in the cluster
        // there is 0 consumers
        // we should stay at size 2

        testManager := &TestResourceManager{
                ExpectedScale:     2,
                MemoryPerProvider: bytesize.MB,
                ResourceProviders: []*ResourceProvider{
                        NewResourceProvider(bytesize.MB, nil),
                        NewResourceProvider(bytesize.MB, nil),
                },
                PendingResources: []ResourceConsumer{},
        }

        if err := testManager.Manager(t).Run(); err != nil {
                t.Fatal(err)
        }
}



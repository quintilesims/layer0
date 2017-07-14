package scheduler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/scheduler/resource"
	"github.com/quintilesims/layer0/api/scheduler/resource/mock_resource"
	"github.com/zpatrick/go-bytesize"
)

type MockProviderManager struct {
	*mock_resource.MockProviderManager
	MemoryPerProvider bytesize.Bytesize
}

func (m *MockProviderManager) CalculateNewProvider(environmentID string) (*resource.ResourceProvider, error) {
	return resource.NewResourceProvider("", false, m.MemoryPerProvider, nil), nil
}

type EnvironmentScalerUnitTest struct {
	ExpectedScale     int
	MemoryPerProvider bytesize.Bytesize
	ResourceProviders []*resource.ResourceProvider
	ResourceConsumers []resource.ResourceConsumer
}

func (e *EnvironmentScalerUnitTest) Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := mock_resource.NewMockConsumerGetter(ctrl)
	mockGetter.EXPECT().
		GetConsumers("eid").
		Return(e.ResourceConsumers, nil)

	mockProvider := &MockProviderManager{
		mock_resource.NewMockProviderManager(ctrl),
		e.MemoryPerProvider,
	}

	mockProvider.EXPECT().
		GetProviders("eid").
		Return(e.ResourceProviders, nil)

	mockProvider.EXPECT().
		ScaleTo("eid", e.ExpectedScale, gomock.Any()).
		Return(0, nil)

	environmentScaler := NewL0EnvironmentScaler(mockGetter, mockProvider)

	if _, err := environmentScaler.Scale("eid"); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_noProviders(t *testing.T) {
	// there are 0 providers in the cluster
	// there is 1 consumer
	// we should scale up to size 1
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     1,
		MemoryPerProvider: bytesize.GB,
		ResourceProviders: []*resource.ResourceProvider{},
		ResourceConsumers: []resource.ResourceConsumer{
			{Memory: bytesize.MB},
		},
	}

	test.Run(t)
}

func TestResourceManagerScaleUp_notEnoughPorts(t *testing.T) {
	// there is 1 provider in the cluster that has port 80 being used
	// there is 1 consumer that needs port 80
	// we should scale up to size 2
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     2,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB, []int{80}),
		},
		ResourceConsumers: []resource.ResourceConsumer{
			{Ports: []int{80}},
		},
	}

	test.Run(t)
}

func TestResourceManagerScaleUp_notEnoughPortsComplex(t *testing.T) {
	// there are 5 providers in the cluster
	// 3 of the consumers can be placed in the current cluster
	// 2 of the consumers will require 1 new provider between the 2 of them
	// 6 of the consumers can be placed across 6 providers
	// we should scale up to size 6
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     6,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB, []int{8000, 8001, 8002}),
			resource.NewResourceProvider("", true, bytesize.MB, []int{8000, 8001, 8002}),
			resource.NewResourceProvider("", true, bytesize.MB, []int{8000, 8001, 8002}),
			resource.NewResourceProvider("", true, bytesize.MB, []int{8000, 8001}),
			resource.NewResourceProvider("", true, bytesize.MB, []int{8000}),
		},
		ResourceConsumers: []resource.ResourceConsumer{
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

	test.Run(t)
}
func TestResourceManagerScaleUp_notEnoughMemory(t *testing.T) {
	// there is 1 provider in the cluster that has 1MB left
	// there is 1 consumer that needs 2MB
	// we should scale up to size 2
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     2,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB, nil),
		},
		ResourceConsumers: []resource.ResourceConsumer{
			{Memory: bytesize.MB * 2},
		},
	}
	test.Run(t)
}
func TestResourceManagerScaleUp_notEnoughMemoryOnASingleProvider(t *testing.T) {
	// there are 2 providers in the cluster
	// combined, they have 3MB of memory left
	// there is 1 consumer that needs 3MB
	// we should scale up to size 3
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB*1, nil),
			resource.NewResourceProvider("", true, bytesize.MB*2, nil),
		},
		ResourceConsumers: []resource.ResourceConsumer{
			{Memory: bytesize.MB * 3},
		},
	}

	test.Run(t)
}
func TestResourceManagerScaleUp_notEnoughMemoryComplex(t *testing.T) {
	// there are 5 providers in the cluster
	// 4 of the consumers can be placed in the current cluster
	// 2 of the consumers will require 1 new provider between the 2 of them
	// 3 of the consumers can be placed across 6 providers
	// we should scale up to size 6
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     6,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB, nil),
			resource.NewResourceProvider("", true, bytesize.MB, nil),
			resource.NewResourceProvider("", true, bytesize.MB, nil),
			resource.NewResourceProvider("", true, bytesize.MB*2, nil),
			resource.NewResourceProvider("", true, bytesize.MB*3, nil),
		},
		ResourceConsumers: []resource.ResourceConsumer{
			// these 4 consumers can be placed in the current cluster
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB * 2},
			// these 2 consumers will require a new provider
			{Memory: bytesize.MB * 2},
			{Memory: bytesize.MB * 2},
			// these 3 consumers can be placed in the cluster
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
		},
	}

	test.Run(t)
}
func TestResourceManagerScaleUp_notEnoughPortsOrMemory(t *testing.T) {
	// there are 2 providers in the cluster
	// 1 of the consumers will require a new provider due to ports
	// 1 of the consumers will require a new provider due to memory
	// 3 of the consumers can be placed across 4 providers
	// we should scale up to size 4
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     4,
		MemoryPerProvider: bytesize.MB * 2,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB, []int{80}),
			resource.NewResourceProvider("", true, bytesize.MB, []int{80}),
		},
		ResourceConsumers: []resource.ResourceConsumer{
			// this consumer will require a new provider for ports
			{Memory: bytesize.MB, Ports: []int{80}},
			// this consumer wil require a new provider for memory
			{Memory: bytesize.MB * 2},
			// these 4 consumers can be placed in the cluster
			{Ports: []int{80}},
			{Memory: bytesize.MB, Ports: []int{8000}},
			{Memory: bytesize.MB, Ports: []int{8000}},
			{Memory: bytesize.MB, Ports: []int{8000}},
		},
	}

	test.Run(t)
}
func TestResourceManagerNoScale_noResourceConsumers(t *testing.T) {
	// there are 2 providers in the cluster that are in use
	// there are 0 consumers
	// we should stay at size 2
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     2,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB, []int{80}),
			resource.NewResourceProvider("", true, bytesize.MB*0.5, nil),
		},
		ResourceConsumers: []resource.ResourceConsumer{},
	}

	test.Run(t)
}
func TestResourceManagerNoScale_enoughPorts(t *testing.T) {
	// there are 3 providers in the cluster
	// there are 6 consumers that can be placed across the cluster
	// we should stay at size 3
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB, []int{8000, 8001, 8002}),
			resource.NewResourceProvider("", true, bytesize.MB, []int{8000, 8001}),
			resource.NewResourceProvider("", true, bytesize.MB, []int{8000}),
		},
		ResourceConsumers: []resource.ResourceConsumer{
			{Ports: []int{8001}},
			{Ports: []int{8002}},
			{Ports: []int{8002}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
			{Ports: []int{8003}},
		},
	}
	test.Run(t)
}
func TestResourceManagerNoScale_enoughMemory(t *testing.T) {
	// there are 3 providers in the cluster
	// there are 6 consumers that can be placed across the cluster
	// we should stay at size 3
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB, nil),
			resource.NewResourceProvider("", true, bytesize.MB*2, nil),
			resource.NewResourceProvider("", true, bytesize.MB*3, nil),
		},
		ResourceConsumers: []resource.ResourceConsumer{
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB * 2},
			{Memory: bytesize.MB},
			{Memory: bytesize.MB},
		},
	}

	test.Run(t)
}
func TestResourceManagerNoScale_enoughMemoryAndPorts(t *testing.T) {
	// there are 3 providers in the cluster
	// there are 4 consumers that can be placed across the cluster
	// we should stay at size 3
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", true, bytesize.MB*1, []int{8000, 8001, 8002}),
			resource.NewResourceProvider("", true, bytesize.MB*3, []int{8000, 8001}),
			resource.NewResourceProvider("", true, bytesize.MB*2, []int{8000}),
		},
		ResourceConsumers: []resource.ResourceConsumer{
			// note that if we place this consumer in the 2nd provider, we would fail
			// to place the 3MB consumer
			{Memory: bytesize.MB * 1, Ports: []int{8002}},
			{Memory: bytesize.MB * 3},
			{Memory: bytesize.MB * 1, Ports: []int{8001}},
			{Memory: bytesize.MB},
		},
	}

	test.Run(t)
}
func TestResourceManagerScaledown_noConsumers(t *testing.T) {
	// there is 1 provider in the cluster that isn't in use
	// there are 0 consumers
	// we should scale to size 0
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     0,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", false, bytesize.MB, nil),
		},
		ResourceConsumers: []resource.ResourceConsumer{},
	}

	test.Run(t)
}
func TestResourceManagerScaledown_complex(t *testing.T) {
	// there are 5 providers in the cluster
	// 2 are not in use
	// there are 2 consumers that can be placed in the 3 already-used consumers
	// we should scale to size 3
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     3,
		MemoryPerProvider: bytesize.MB * 4,
		ResourceProviders: []*resource.ResourceProvider{
			resource.NewResourceProvider("", false, bytesize.MB*4, nil),
			resource.NewResourceProvider("", false, bytesize.MB*4, nil),
			resource.NewResourceProvider("", true, bytesize.MB*2, []int{8000}),
			resource.NewResourceProvider("", true, bytesize.MB*2, []int{8001}),
			resource.NewResourceProvider("", true, bytesize.MB*2, []int{8002}),
		},
		ResourceConsumers: []resource.ResourceConsumer{
			{Memory: bytesize.MB, Ports: []int{8000}},
			{Memory: bytesize.MB, Ports: []int{8001}},
			{Memory: bytesize.MB, Ports: []int{8002}},
		},
	}

	test.Run(t)
}

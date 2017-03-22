package ecsbackend

import (
	awsasg "github.com/aws/aws-sdk-go/service/autoscaling"
	awsecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/api/scheduler/resource"
	"github.com/quintilesims/layer0/common/aws/autoscaling"
	"github.com/quintilesims/layer0/common/aws/autoscaling/mock_autoscaling"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/aws/ecs/mock_ecs"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/zpatrick/go-bytesize"
	"testing"
)

type MockResourceManager struct {
	ECS         *mock_ecs.MockProvider
	Autoscaling *mock_autoscaling.MockProvider
}

func newMockResourceManager(t *testing.T) (*MockResourceManager, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rm := &MockResourceManager{
		ECS:         mock_ecs.NewMockProvider(ctrl),
		Autoscaling: mock_autoscaling.NewMockProvider(ctrl),
	}

	return rm, ctrl
}

func (m *MockResourceManager) ResourceManager() *ECSResourceManager {
	return NewECSResourceManager(m.ECS, m.Autoscaling)
}

func TestResourceManager_GetProviders(t *testing.T) {
	rm, ctrl := newMockResourceManager(t)
	defer ctrl.Finish()

	environmentID := id.L0EnvironmentID("eid")

	rm.ECS.EXPECT().
		ListContainerInstances(environmentID.ECSEnvironmentID().String()).
		Return([]*string{stringp("i1")}, nil)

	containerInstances := []*ecs.ContainerInstance{
		{
			&awsecs.ContainerInstance{
				Status:            stringp("ACTIVE"),
				RunningTasksCount: int64p(1),
				PendingTasksCount: int64p(1),
				RemainingResources: []*awsecs.Resource{
					{
						Name:         stringp("MEMORY"),
						IntegerValue: int64p(500),
					},
					{
						Name: stringp("PORTS"),
						StringSetValue: []*string{
							stringp("80"),
							stringp("8000"),
						},
					},
				},
			},
		},
		{
			&awsecs.ContainerInstance{
				Status:            stringp("ACTIVE"),
				RunningTasksCount: int64p(0),
				PendingTasksCount: int64p(0),
				RemainingResources: []*awsecs.Resource{
					{
						Name:         stringp("MEMORY"),
						IntegerValue: int64p(1000),
					},
					{
						Name: stringp("PORTS"),
						StringSetValue: []*string{
							stringp("80"),
						},
					},
				},
			},
		},
		{
			&awsecs.ContainerInstance{
				Status: stringp("INACTIVE"),
				RemainingResources: []*awsecs.Resource{
					{
						Name:         stringp("MEMORY"),
						IntegerValue: int64p(1000),
					},
				},
			},
		},
	}

	rm.ECS.EXPECT().
		DescribeContainerInstances(environmentID.ECSEnvironmentID().String(), gomock.Any()).
		Return(containerInstances, nil)

	providers, err := rm.ResourceManager().GetProviders("eid")
	if err != nil {
		t.Fatal(err)
	}

	expected := []*resource.ResourceProvider{
		resource.NewResourceProvider("", true, bytesize.MiB*500, []int{80, 8000}),
		resource.NewResourceProvider("", false, bytesize.MiB*1000, []int{80}),
	}

	testutils.AssertEqual(t, expected, providers)
}

func TestResourceManager_scaleUp(t *testing.T) {
	rm, ctrl := newMockResourceManager(t)
	defer ctrl.Finish()

	environmentID := id.L0EnvironmentID("eid")
	rm.Autoscaling.EXPECT().
		UpdateAutoScalingGroupMaxSize("asg_name", 5).
		Return(nil)

	rm.Autoscaling.EXPECT().
		SetDesiredCapacity("asg_name", 5).
		Return(nil)

	asg := &autoscaling.Group{
		&awsasg.Group{
			AutoScalingGroupName: stringp("asg_name"),
			MaxSize:              int64p(3),
			MinSize:              int64p(0),
		},
	}

	scale, err := rm.ResourceManager().scaleUp(environmentID.ECSEnvironmentID(), 5, asg)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, scale, 5)
}

func TestResourceManager_scaleDown(t *testing.T) {
        rm, ctrl := newMockResourceManager(t)
        defer ctrl.Finish()

        environmentID := id.L0EnvironmentID("eid")

        rm.Autoscaling.EXPECT().
                SetDesiredCapacity("asg_name", 0).
                Return(nil)

        asg := &autoscaling.Group{
                &awsasg.Group{
                        AutoScalingGroupName: stringp("asg_name"),
                        MaxSize:              int64p(3),
                        MinSize:              int64p(0),
                        DesiredCapacity:      int64p(3),
                },
        }

        scale, err := rm.ResourceManager().scaleDown(environmentID.ECSEnvironmentID(), 0, asg, nil)
        if err != nil {
                t.Fatal(err)
        }

        testutils.AssertEqual(t, scale, 0)
}

func TestResourceManager_scaleDownStayAboveMin(t *testing.T) {
	rm, ctrl := newMockResourceManager(t)
	defer ctrl.Finish()

	environmentID := id.L0EnvironmentID("eid")

	rm.Autoscaling.EXPECT().
		SetDesiredCapacity("asg_name", 1).
		Return(nil)

	asg := &autoscaling.Group{
		&awsasg.Group{
			AutoScalingGroupName: stringp("asg_name"),
			MaxSize:              int64p(3),
			MinSize:              int64p(1),
			DesiredCapacity:      int64p(3),
		},
	}

	scale, err := rm.ResourceManager().scaleDown(environmentID.ECSEnvironmentID(), 0, asg, nil)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, scale, 1)
}

package ecsbackend

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	aws_ec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/api/backend/mock_backend"
	"github.com/quintilesims/layer0/common/aws/ec2"
	"github.com/quintilesims/layer0/common/aws/ec2/mock_ec2"
	"github.com/quintilesims/layer0/common/aws/elb"
	"github.com/quintilesims/layer0/common/aws/elb/mock_elb"
	"github.com/quintilesims/layer0/common/aws/iam/mock_iam"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

type MockECSLoadBalancerManager struct {
	EC2     *mock_ec2.MockProvider
	ELB     *mock_elb.MockProvider
	IAM     *mock_iam.MockProvider
	Backend *mock_backend.MockBackend
}

func NewMockECSLoadBalancerManager(ctrl *gomock.Controller) *MockECSLoadBalancerManager {
	return &MockECSLoadBalancerManager{
		EC2:     mock_ec2.NewMockProvider(ctrl),
		ELB:     mock_elb.NewMockProvider(ctrl),
		IAM:     mock_iam.NewMockProvider(ctrl),
		Backend: mock_backend.NewMockBackend(ctrl),
	}
}

func (this *MockECSLoadBalancerManager) LoadBalancer() *ECSLoadBalancerManager {
	return NewECSLoadBalancerManager(this.EC2, this.ELB, this.IAM, this.Backend)
}

func makeSubnet(az string) *ec2.Subnet {
	return &ec2.Subnet{
		&aws_ec2.Subnet{
			AvailabilityZone: &az,
		},
	}
}

func TestGetLoadBalancer(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should call ecs.DescribeLoadBalancer with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				loadBalancerID := id.L0LoadBalancerID("lbid").ECSLoadBalancerID()
				loadBalancer := elb.NewLoadBalancerDescription(loadBalancerID.String(), "", nil)

				mockLB.ELB.EXPECT().
					DescribeLoadBalancer(loadBalancerID.String()).
					Return(loadBalancer, nil)

				loadBalancerAttributes := elb.NewLoadBalancerAttributes()

				mockLB.ELB.EXPECT().
					DescribeLoadBalancerAttributes(gomock.Any()).
					Return(loadBalancerAttributes, nil)

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)
				manager.GetLoadBalancer("lbid")
			},
		},
		{
			Name: "Should return layer0-formatted load balancer id",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				loadBalancerID := id.L0LoadBalancerID("lbid").ECSLoadBalancerID()
				loadBalancer := elb.NewLoadBalancerDescription(loadBalancerID.String(), "", nil)

				mockLB.ELB.EXPECT().
					DescribeLoadBalancer(gomock.Any()).
					Return(loadBalancer, nil)

				loadBalancerAttributes := elb.NewLoadBalancerAttributes()

				mockLB.ELB.EXPECT().
					DescribeLoadBalancerAttributes(gomock.Any()).
					Return(loadBalancerAttributes, nil)

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)

				loadBalancer, err := manager.GetLoadBalancer("lbid")
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(loadBalancer.LoadBalancerID, "lbid")
			},
		},
		{
			Name: "Should propagate elb.DescribeLoadBalancer error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				mockLB.ELB.EXPECT().
					DescribeLoadBalancer(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)

				if _, err := manager.GetLoadBalancer("lbid"); err == nil {
					reporter.Fatalf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListLoadBalancers(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should return layer0-formatted load balancer ids",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				loadBalancerID := id.L0LoadBalancerID("lbid").ECSLoadBalancerID()
				loadBalancer := elb.NewLoadBalancerDescription(loadBalancerID.String(), "", nil)

				mockLB.ELB.EXPECT().
					DescribeLoadBalancers().
					Return([]*elb.LoadBalancerDescription{loadBalancer}, nil)

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)

				loadBalancers, err := manager.ListLoadBalancers()
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(len(loadBalancers), 1)
				reporter.AssertEqual(loadBalancers[0].LoadBalancerID, "lbid")
			},
		},
		{
			Name: "Should propagate elb.DescribeLoadBalancers error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				mockLB.ELB.EXPECT().
					DescribeLoadBalancers().
					Return(nil, fmt.Errorf("some error"))

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)

				if _, err := manager.ListLoadBalancers(); err == nil {
					reporter.Fatalf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateLoadBalancer(t *testing.T) {
	defer id.StubIDGeneration("lbid")()

	sgList := &ec2.SecurityGroup{
		&aws_ec2.SecurityGroup{
			GroupId: stringp("some_group_id"),
		},
	}

	testCases := []testutils.TestCase{
		{
			Name: "Should use ECS-formatted IDs in AWS calls",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				loadBalancerID := id.L0LoadBalancerID("lbid").ECSLoadBalancerID()
				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				roleName := loadBalancerID.RoleName()

				mockLB.IAM.EXPECT().
					CreateRole(roleName, "ecs.amazonaws.com")

				mockLB.IAM.EXPECT().
					PutRolePolicy(roleName, gomock.Any())

				mockLB.IAM.EXPECT().
					GetAccountId().
					Return("100", nil)

				// getSubnetsAndAvailZones
				mockLB.EC2.EXPECT().
					DescribeSubnet(gomock.Any()).
					Return(makeSubnet("a"), nil)

				mockLB.EC2.EXPECT().
					DescribeSubnet(gomock.Any()).
					Return(makeSubnet("b"), nil)

				mockLB.EC2.EXPECT().
					DescribeSecurityGroup(loadBalancerID.SecurityGroupName()).
					Return(sgList, nil)

				mockLB.EC2.EXPECT().
					DescribeSecurityGroup(environmentID.SecurityGroupName()).
					Return(sgList, nil)

				mockLB.ELB.EXPECT().
					CreateLoadBalancer(loadBalancerID.String(), "internet-facing", gomock.Any(), gomock.Any(), gomock.Any())

				elbHealthCheck := elb.NewHealthCheck("", 0, 0, 0, 0)
				mockLB.ELB.EXPECT().
					ConfigureHealthCheck(loadBalancerID.String(), elbHealthCheck)

				mockLB.ELB.EXPECT().
					SetIdleTimeout(loadBalancerID.String(), 60)

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)
				manager.CreateLoadBalancer("lb_name", "envid", true, nil, models.HealthCheck{}, 60)
			},
		},
		{
			Name: "Should not create security group for private load balancers",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				mockLB.IAM.EXPECT().
					CreateRole(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				mockLB.IAM.EXPECT().
					PutRolePolicy(gomock.Any(), gomock.Any()).
					Return(nil)

				mockLB.IAM.EXPECT().
					GetAccountId().
					Return("100", nil)

				// getSubnetsAndAvailZones
				mockLB.EC2.EXPECT().
					DescribeSubnet(gomock.Any()).
					Return(makeSubnet("a"), nil)

				mockLB.EC2.EXPECT().
					DescribeSubnet(gomock.Any()).
					Return(makeSubnet("b"), nil)

				mockLB.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any()).
					Return(sgList, nil)

				mockLB.ELB.EXPECT().
					CreateLoadBalancer(gomock.Any(), "internal", gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil)

				mockLB.ELB.EXPECT().
					ConfigureHealthCheck(gomock.Any(), gomock.Any()).
					Return(nil)

				mockLB.ELB.EXPECT().
					SetIdleTimeout(gomock.Any(), gomock.Any()).
					Return(nil)

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)
				manager.CreateLoadBalancer("lb_name", "envid", false, nil, models.HealthCheck{}, 60)
			},
		},
		{
			Name: "Should pass through idempotent aws errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				mockLB.IAM.EXPECT().
					CreateRole(gomock.Any(), gomock.Any()).
					Return(nil, awserr.New("EntityAlreadyExists", "some message", nil))

				mockLB.IAM.EXPECT().
					PutRolePolicy(gomock.Any(), gomock.Any())

				mockLB.IAM.EXPECT().
					GetAccountId().
					Return("100", nil)

				mockLB.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any()).
					Return(sgList, nil).
					AnyTimes()

				// getSubnetsAndAvailZones
				mockLB.EC2.EXPECT().
					DescribeSubnet(gomock.Any()).
					Return(makeSubnet("a"), nil)

				mockLB.EC2.EXPECT().
					DescribeSubnet(gomock.Any()).
					Return(makeSubnet("b"), nil)

				mockLB.ELB.EXPECT().
					CreateLoadBalancer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

				mockLB.ELB.EXPECT().
					ConfigureHealthCheck(gomock.Any(), gomock.Any())

				mockLB.ELB.EXPECT().
					SetIdleTimeout(gomock.Any(), gomock.Any())

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)
				manager.CreateLoadBalancer("lb_name", "envid", true, nil, models.HealthCheck{}, 60)
			},
		},
		{
			Name: "Should return L0-formatted IDs",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				mockLB.IAM.EXPECT().
					CreateRole(gomock.Any(), gomock.Any())

				mockLB.IAM.EXPECT().
					PutRolePolicy(gomock.Any(), gomock.Any())

				mockLB.IAM.EXPECT().
					GetAccountId().
					Return("100", nil)

				mockLB.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any()).
					Return(sgList, nil).
					AnyTimes()

				// getSubnetsAndAvailZones
				mockLB.EC2.EXPECT().
					DescribeSubnet(gomock.Any()).
					Return(makeSubnet("a"), nil)

				mockLB.EC2.EXPECT().
					DescribeSubnet(gomock.Any()).
					Return(makeSubnet("b"), nil)

				mockLB.ELB.EXPECT().
					CreateLoadBalancer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

				mockLB.ELB.EXPECT().
					ConfigureHealthCheck(gomock.Any(), gomock.Any())

				mockLB.ELB.EXPECT().
					SetIdleTimeout(gomock.Any(), gomock.Any())

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)

				loadBalancerID := id.L0LoadBalancerID("lbid")
				environmentID := id.L0EnvironmentID("envid")

				model, err := manager.CreateLoadBalancer("lb_name", environmentID.String(), true, nil, models.HealthCheck{}, 60)
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(model.LoadBalancerID, loadBalancerID.String())
				reporter.AssertEqual(model.EnvironmentID, environmentID.String())
			},
		},
		{
			Name: "Should propagate unexpected aws errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				return func(g testutils.ErrorGenerator) interface{} {
					mockLB := NewMockECSLoadBalancerManager(ctrl)

					mockLB.IAM.EXPECT().
						CreateRole(gomock.Any(), gomock.Any()).
						Return(nil, g.Error())

					mockLB.IAM.EXPECT().
						PutRolePolicy(gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockLB.IAM.EXPECT().
						GetAccountId().
						Return("100", g.Error()).
						AnyTimes()

					// getSubnetsAndAvailZones
					mockLB.EC2.EXPECT().
						DescribeSubnet(gomock.Any()).
						Return(makeSubnet("a"), g.Error()).
						AnyTimes()

					mockLB.EC2.EXPECT().
						DescribeSubnet(gomock.Any()).
						Return(makeSubnet("b"), g.Error()).
						AnyTimes()

					// test this failure for both lb and env calls
					mockLB.EC2.EXPECT().
						DescribeSecurityGroup(gomock.Any()).
						Return(sgList, g.Error()).
						AnyTimes()

					mockLB.ELB.EXPECT().
						CreateLoadBalancer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil, g.Error()).
						AnyTimes()

					return mockLB.LoadBalancer()
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				setup := target.(func(testutils.ErrorGenerator) interface{})

				for i := 0; i < 8; i++ {
					var g testutils.ErrorGenerator
					g.Set(i+1, fmt.Errorf("some error"))

					manager := setup(g).(*ECSLoadBalancerManager)
					if _, err := manager.CreateLoadBalancer("", "", true, nil, models.HealthCheck{}, 60); err == nil {
						reporter.Errorf("Error on variation %d, Error was nil!", i)
					}
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteLoadBalancer(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should delete role policies and role correctly",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				loadBalancerID := id.L0LoadBalancerID("lbid").ECSLoadBalancerID()
				roleName := loadBalancerID.RoleName()
				policyName := stringp("some_policy")

				mockLB.IAM.EXPECT().
					ListRolePolicies(roleName).
					Return([]*string{policyName}, nil)

				// for waiters
				mockLB.IAM.EXPECT().
					ListRolePolicies(gomock.Any()).
					Return([]*string{}, nil).
					AnyTimes()

				mockLB.IAM.EXPECT().
					DeleteRolePolicy(roleName, *policyName).
					Return(nil)

				mockLB.IAM.EXPECT().
					DeleteRole(roleName).
					Return(nil)

				mockLB.ELB.EXPECT().
					DeleteLoadBalancer(gomock.Any())

				mockLB.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any())

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)
				manager.DeleteLoadBalancer("lbid")
			},
		},
		{
			Name: "Should call manager.DeleteLoadBalancer with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				loadBalancerID := id.L0LoadBalancerID("lbid").ECSLoadBalancerID()

				mockLB.IAM.EXPECT().
					ListRolePolicies(gomock.Any()).
					AnyTimes()

				mockLB.IAM.EXPECT().
					DeleteRole(gomock.Any())

				mockLB.ELB.EXPECT().
					DeleteLoadBalancer(loadBalancerID.String())

				mockLB.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any())

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)
				manager.DeleteLoadBalancer("lbid")
			},
		},
		{
			Name: "Should delete security groups correctly",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				loadBalancerID := id.L0LoadBalancerID("lbid").ECSLoadBalancerID()
				sgName := loadBalancerID.SecurityGroupName()
				sg := ec2.NewSecurityGroup("some_id")

				mockLB.IAM.EXPECT().
					ListRolePolicies(gomock.Any()).
					AnyTimes()

				mockLB.IAM.EXPECT().
					DeleteRole(gomock.Any())

				mockLB.ELB.EXPECT().
					DeleteLoadBalancer(gomock.Any())

				mockLB.EC2.EXPECT().
					DescribeSecurityGroup(sgName).
					Return(sg, nil)

				mockLB.EC2.EXPECT().
					DeleteSecurityGroup(sg).
					Return(nil)

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)
				manager.DeleteLoadBalancer("lbid")
			},
		},
		{
			Name: "Should pass through idempotent aws errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				sg := ec2.NewSecurityGroup("some_id")

				mockLB.IAM.EXPECT().
					ListRolePolicies(gomock.Any()).
					Return(nil, awserr.New("NoSuchEntity", "some message", nil)).
					AnyTimes()

				mockLB.IAM.EXPECT().
					DeleteRole(gomock.Any()).
					Return(awserr.New("NoSuchEntity", "some message", nil))

				mockLB.ELB.EXPECT().
					DeleteLoadBalancer(gomock.Any()).
					Return(awserr.New("NoSuchEntity", "some message", nil))

				mockLB.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any()).
					Return(sg, nil)

				// do retry for DeleteSecurityGroup. We don't check error codes here
				mockLB.EC2.EXPECT().
					DeleteSecurityGroup(sg).
					Return(fmt.Errorf(""))

				mockLB.EC2.EXPECT().
					DeleteSecurityGroup(sg).
					Return(nil)

				manager := mockLB.LoadBalancer()
				manager.Clock = &testutils.StubClock{}
				return manager
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)

				if err := manager.DeleteLoadBalancer("lbid"); err != nil {
					reporter.Fatal(err)
				}
			},
		},
		{
			Name: "Should propagate unexpected aws errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				return func(g testutils.ErrorGenerator) interface{} {
					mockLB := NewMockECSLoadBalancerManager(ctrl)

					mockLB.IAM.EXPECT().
						ListRolePolicies(gomock.Any()).
						Return(nil, g.Error()).
						AnyTimes()

					mockLB.IAM.EXPECT().
						DeleteRole(gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockLB.ELB.EXPECT().
						DeleteLoadBalancer(gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockLB.EC2.EXPECT().
						DescribeSecurityGroup(gomock.Any()).
						Return(nil, g.Error()).
						AnyTimes()

					return mockLB.LoadBalancer()
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				setup := target.(func(testutils.ErrorGenerator) interface{})

				for i := 0; i < 4; i++ {
					var g testutils.ErrorGenerator
					g.Set(i+1, fmt.Errorf("some error"))

					manager := setup(g).(*ECSLoadBalancerManager)
					if err := manager.DeleteLoadBalancer("lbid"); err == nil {
						reporter.Errorf("Error on variation %d, Error was nil!", i)
					}
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestUpdateLoadBalancerHealthCheck(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should pass proper params to ELB.ConfigureHealthCheck.",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				loadBalancerID := id.L0LoadBalancerID("lbid").ECSLoadBalancerID()
				loadBalancer := elb.NewLoadBalancerDescription(loadBalancerID.String(), "", nil)

				mockLB.ELB.EXPECT().
					DescribeLoadBalancer(loadBalancerID.String()).
					Return(loadBalancer, nil)

				loadBalancerAttributes := elb.NewLoadBalancerAttributes()

				mockLB.ELB.EXPECT().
					DescribeLoadBalancerAttributes(gomock.Any()).
					Return(loadBalancerAttributes, nil)

				mockLB.ELB.EXPECT().
					ConfigureHealthCheck(loadBalancerID.String(), elb.NewHealthCheck("TCP:80", 30, 5, 2, 2))

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)

				healthCheck := models.HealthCheck{
					Target:             "TCP:80",
					Interval:           30,
					Timeout:            5,
					HealthyThreshold:   2,
					UnhealthyThreshold: 2,
				}

				model, err := manager.UpdateLoadBalancerHealthCheck("lbid", healthCheck)
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(healthCheck, model.HealthCheck)
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestUpdateLoadBalancerIdleTimeout(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should pass proper params to ELB.SetIdleTimeout.",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLB := NewMockECSLoadBalancerManager(ctrl)

				loadBalancerID := id.L0LoadBalancerID("lbid").ECSLoadBalancerID()
				loadBalancer := elb.NewLoadBalancerDescription(loadBalancerID.String(), "", nil)

				mockLB.ELB.EXPECT().
					DescribeLoadBalancer(loadBalancerID.String()).
					Return(loadBalancer, nil)

				loadBalancerAttributes := elb.NewLoadBalancerAttributes()

				mockLB.ELB.EXPECT().
					DescribeLoadBalancerAttributes(gomock.Any()).
					Return(loadBalancerAttributes, nil)

				mockLB.ELB.EXPECT().
					SetIdleTimeout(loadBalancerID.String(), 60)

				return mockLB.LoadBalancer()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSLoadBalancerManager)

				idleTimeout := 60

				model, err := manager.UpdateLoadBalancerIdleTimeout("lbid", idleTimeout)
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(idleTimeout, model.IdleTimeout)
			},
		},
	}

	testutils.RunTests(t, testCases)
}

// todo: UpdateLoadBalancerPorts

package test_aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadBalancerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name",
		},
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "env_id",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	deleteLBInput := &elb.DeleteLoadBalancerInput{}
	deleteLBInput.SetLoadBalancerName("l0-test-lb_id")

	mockAWS.ELB.EXPECT().
		DeleteLoadBalancer(deleteLBInput).
		Return(&elb.DeleteLoadBalancerOutput{}, nil)

	deleteRolePolicyInput := &iam.DeleteRolePolicyInput{}
	deleteRolePolicyInput.SetRoleName("l0-test-lb_id-lb")
	deleteRolePolicyInput.SetPolicyName("l0-test-lb_id-lb")

	mockAWS.IAM.EXPECT().
		DeleteRolePolicy(deleteRolePolicyInput).
		Return(&iam.DeleteRolePolicyOutput{}, nil)

	deleteRoleInput := &iam.DeleteRoleInput{}
	deleteRoleInput.SetRoleName("l0-test-lb_id-lb")

	mockAWS.IAM.EXPECT().
		DeleteRole(deleteRoleInput).
		Return(&iam.DeleteRoleOutput{}, nil)

	readSGHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg")
	deleteSGHelper(mockAWS, "lb_sg")

	getRoleInput := &iam.GetRoleInput{}
	getRoleInput.SetRoleName("l0-test-lb_id-lb")

	mockAWS.IAM.EXPECT().
		GetRole(getRoleInput).
		Return(&iam.GetRoleOutput{}, nil).
		AnyTimes()

	getRolePolicyInput := &iam.GetRolePolicyInput{}
	getRolePolicyInput.SetRoleName("l0-test-lb_id-lb")
	getRolePolicyInput.SetPolicyName("l0-test-lb_id-lb")

	mockAWS.IAM.EXPECT().
		GetRolePolicy(getRolePolicyInput).
		Return(&iam.GetRolePolicyOutput{}, nil).
		AnyTimes()

	describeLoadBalancersInput := &elb.DescribeLoadBalancersInput{}
	describeLoadBalancersInput.SetLoadBalancerNames([]*string{aws.String("l0-test-lb_id")})
	describeLoadBalancersInput.SetPageSize(1)

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(describeLoadBalancersInput).
		Return(&elb.DescribeLoadBalancersOutput{}, nil).
		AnyTimes()

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("lb_id"); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, tagStore.Tags(), 0)
}

func TestLoadBalancerDeleteIdempotence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	mockAWS.ELB.EXPECT().
		DeleteLoadBalancer(gomock.Any()).
		Return(nil, awserr.New("NoSuchEntity", "", nil))

	mockAWS.IAM.EXPECT().
		DeleteRolePolicy(gomock.Any()).
		Return(nil, awserr.New("NoSuchEntity", "", nil))

	mockAWS.IAM.EXPECT().
		DeleteRole(gomock.Any()).
		Return(nil, awserr.New("NoSuchEntity", "", nil))

	mockAWS.EC2.EXPECT().
		DescribeSecurityGroups(gomock.Any()).
		Return(&ec2.DescribeSecurityGroupsOutput{}, nil).
		AnyTimes()

	mockAWS.IAM.EXPECT().
		GetRole(gomock.Any()).
		Return(&iam.GetRoleOutput{}, nil).
		AnyTimes()

	mockAWS.IAM.EXPECT().
		GetRolePolicy(gomock.Any()).
		Return(&iam.GetRolePolicyOutput{}, nil).
		AnyTimes()

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(gomock.Any()).
		Return(&elb.DescribeLoadBalancersOutput{}, nil).
		AnyTimes()

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("lb_id"); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerDelete_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name",
		},
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "env_id",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	deleteLBInput := &elb.DeleteLoadBalancerInput{}
	deleteLBInput.SetLoadBalancerName("l0-test-lb_id")

	mockAWS.ELB.EXPECT().
		DeleteLoadBalancer(deleteLBInput).
		Return(&elb.DeleteLoadBalancerOutput{}, nil)

	deleteRolePolicyInput := &iam.DeleteRolePolicyInput{}
	deleteRolePolicyInput.SetRoleName("l0-test-lb_id-lb")
	deleteRolePolicyInput.SetPolicyName("l0-test-lb_id-lb")

	mockAWS.IAM.EXPECT().
		DeleteRolePolicy(deleteRolePolicyInput).
		Return(&iam.DeleteRolePolicyOutput{}, nil)

	deleteRoleInput := &iam.DeleteRoleInput{}
	deleteRoleInput.SetRoleName("l0-test-lb_id-lb")

	mockAWS.IAM.EXPECT().
		DeleteRole(deleteRoleInput).
		Return(&iam.DeleteRoleOutput{}, nil)

	readSGHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg")
	deleteSGHelper(mockAWS, "lb_sg")

	getRoleInput := &iam.GetRoleInput{}
	getRoleInput.SetRoleName("l0-test-lb_id-lb")

	mockAWS.IAM.EXPECT().
		GetRole(getRoleInput).
		Return(&iam.GetRoleOutput{}, nil).
		AnyTimes()

	getRolePolicyInput := &iam.GetRolePolicyInput{}
	getRolePolicyInput.SetRoleName("l0-test-lb_id-lb")
	getRolePolicyInput.SetPolicyName("l0-test-lb_id-lb")

	mockAWS.IAM.EXPECT().
		GetRolePolicy(getRolePolicyInput).
		Return(&iam.GetRolePolicyOutput{}, nil).
		AnyTimes()

	describeLoadBalancersInput := &elb.DescribeLoadBalancersInput{}
	describeLoadBalancersInput.SetLoadBalancerNames([]*string{aws.String("l0-test-lb_id")})
	describeLoadBalancersInput.SetPageSize(1)

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(describeLoadBalancersInput).
		Return(&elb.DescribeLoadBalancersOutput{}, fmt.Errorf("Error")).
		AnyTimes()

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("lb_id"); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() != "FailedRequestTimeout" {
			t.Fatal(err)
		}
	}

	assert.Len(t, tagStore.Tags(), 2)
}

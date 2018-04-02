package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestClassicLoadBalancerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	describeLBInput := &elb.DescribeLoadBalancersInput{}
	describeLBInput.LoadBalancerNames = []*string{aws.String("l0-test-lb_id")}
	describeLBInput.SetPageSize(1)

	describeLBOutput := &elb.DescribeLoadBalancersOutput{}
	describeLBOutput.LoadBalancerDescriptions = []*elb.LoadBalancerDescription{
		{
			LoadBalancerName: aws.String("l0-test-lb_id"),
		},
	}

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(describeLBInput).
		Return(describeLBOutput, nil)

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
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "type",
			Value:      string(models.ClassicLoadBalancerType),
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

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(gomock.Any()).
		Return(&elb.DescribeLoadBalancersOutput{}, awserr.New("LoadBalancerNotFound", "", nil)).
		AnyTimes()

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

	describeTGInput := &alb.DescribeTargetGroupsInput{}
	describeTGInput.SetNames([]*string{aws.String("l0-test-lb_id")})
	describeTGOutput := &alb.DescribeTargetGroupsOutput{}
	describeTGOutput.SetTargetGroups([]*alb.TargetGroup{
		{
			TargetGroupArn: aws.String("arn:l0-test-lb_id"),
		},
	})

	mockAWS.ALB.EXPECT().
		DescribeTargetGroups(describeTGInput).
		Return(describeTGOutput, nil)

	deleteTGInput := &alb.DeleteTargetGroupInput{}
	deleteTGInput.SetTargetGroupArn("arn:l0-test-lb_id")

	mockAWS.ALB.EXPECT().
		DeleteTargetGroup(deleteTGInput).
		Return(&alb.DeleteTargetGroupOutput{}, nil)

	readSGHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg")
	deleteSGHelper(mockAWS, "lb_sg")

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName("l0-test-lb_id-lb")
	securityGroup.SetGroupId("lb_sg")

	describeSecurityGroupsOutput := &ec2.DescribeSecurityGroupsOutput{}
	describeSecurityGroupsOutput.SetSecurityGroups([]*ec2.SecurityGroup{securityGroup})

	gomock.InOrder(
		mockAWS.EC2.EXPECT().
			DescribeSecurityGroups(gomock.Any()).
			Return(describeSecurityGroupsOutput, nil),
		mockAWS.EC2.EXPECT().
			DescribeSecurityGroups(gomock.Any()).
			Return(&ec2.DescribeSecurityGroupsOutput{}, nil),
	)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("lb_id"); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, tagStore.Tags(), 0)
}

func TestApplicationLoadBalancerDelete(t *testing.T) {
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
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "type",
			Value:      models.ApplicationLoadBalancerType,
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	describeELBInput := &elb.DescribeLoadBalancersInput{}
	describeELBInput.LoadBalancerNames = []*string{aws.String("l0-test-lb_id")}
	describeELBInput.SetPageSize(1)
	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(describeELBInput).
		Return(nil, awserr.New("LoadBalancerNotFound", "", nil)).
		AnyTimes()

	describeALBInput := &alb.DescribeLoadBalancersInput{}
	describeALBInput.SetNames([]*string{aws.String("l0-test-lb_id")})
	describeALBOutput := &alb.DescribeLoadBalancersOutput{
		LoadBalancers: []*alb.LoadBalancer{
			{
				LoadBalancerArn: aws.String("arn:l0-test-lb_id"),
			},
		},
	}

	mockAWS.ALB.EXPECT().
		DescribeLoadBalancers(describeALBInput).
		Return(describeALBOutput, nil).
		AnyTimes()

	deleteApplicationLBInput := &alb.DeleteLoadBalancerInput{}
	deleteApplicationLBInput.SetLoadBalancerArn("arn:l0-test-lb_id")
	mockAWS.ALB.EXPECT().
		DeleteLoadBalancer(deleteApplicationLBInput).
		Return(&alb.DeleteLoadBalancerOutput{}, nil)

	waitInput := &alb.DescribeLoadBalancersInput{}
	waitInput.SetLoadBalancerArns([]*string{aws.String("arn:l0-test-lb_id")})
	mockAWS.ALB.EXPECT().
		WaitUntilLoadBalancersDeleted(waitInput).
		Return(nil)

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

	describeTGInput := &alb.DescribeTargetGroupsInput{}
	describeTGInput.SetNames([]*string{aws.String("l0-test-lb_id")})
	describeTGOutput := &alb.DescribeTargetGroupsOutput{}
	describeTGOutput.SetTargetGroups([]*alb.TargetGroup{
		{
			TargetGroupArn: aws.String("arn:l0-test-lb_id"),
		},
	})

	mockAWS.ALB.EXPECT().
		DescribeTargetGroups(describeTGInput).
		Return(describeTGOutput, nil)

	deleteTGInput := &alb.DeleteTargetGroupInput{}
	deleteTGInput.SetTargetGroupArn("arn:l0-test-lb_id")

	mockAWS.ALB.EXPECT().
		DeleteTargetGroup(deleteTGInput).
		Return(&alb.DeleteTargetGroupOutput{}, nil)

	readSGHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg")
	deleteSGHelper(mockAWS, "lb_sg")

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName("l0-test-lb_id-lb")
	securityGroup.SetGroupId("lb_sg")

	describeSecurityGroupsOutput := &ec2.DescribeSecurityGroupsOutput{}
	describeSecurityGroupsOutput.SetSecurityGroups([]*ec2.SecurityGroup{securityGroup})

	gomock.InOrder(
		mockAWS.EC2.EXPECT().
			DescribeSecurityGroups(gomock.Any()).
			Return(describeSecurityGroupsOutput, nil),
		mockAWS.EC2.EXPECT().
			DescribeSecurityGroups(gomock.Any()).
			Return(&ec2.DescribeSecurityGroupsOutput{}, nil),
	)

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
		DescribeLoadBalancers(gomock.Any()).
		Return(nil, awserr.New("LoadBalancerNotFound", "", nil))

	mockAWS.ALB.EXPECT().
		DescribeLoadBalancers(gomock.Any()).
		Return(nil, awserr.New("LoadBalancerNotFound", "", nil))

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(gomock.Any()).
		Return(&elb.DescribeLoadBalancersOutput{}, awserr.New("LoadBalancerNotFound", "", nil))

	mockAWS.ALB.EXPECT().
		DescribeTargetGroups(gomock.Any()).
		Return(nil, awserr.New(alb.ErrCodeTargetGroupNotFoundException, "", nil))

	mockAWS.IAM.EXPECT().
		DeleteRolePolicy(gomock.Any()).
		Return(nil, awserr.New("NoSuchEntity", "", nil))

	mockAWS.IAM.EXPECT().
		DeleteRole(gomock.Any()).
		Return(nil, awserr.New("NoSuchEntity", "", nil))

	mockAWS.EC2.EXPECT().
		DescribeSecurityGroups(gomock.Any()).
		Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName("l0-test-lb_id-lb")
	securityGroup.SetGroupId("lb_sg")

	describeSecurityGroupsOutput := &ec2.DescribeSecurityGroupsOutput{}
	describeSecurityGroupsOutput.SetSecurityGroups([]*ec2.SecurityGroup{securityGroup})

	gomock.InOrder(
		mockAWS.EC2.EXPECT().
			DescribeSecurityGroups(gomock.Any()).
			Return(describeSecurityGroupsOutput, nil),
		mockAWS.EC2.EXPECT().
			DescribeSecurityGroups(gomock.Any()).
			Return(&ec2.DescribeSecurityGroupsOutput{}, nil),
	)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("lb_id"); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerDeleteRetry(t *testing.T) {
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

	mockAWS.ELB.EXPECT().
		DeleteLoadBalancer(gomock.Any()).
		Return(&elb.DeleteLoadBalancerOutput{}, nil)

	lb := &elb.LoadBalancerDescription{}
	lb.SetLoadBalancerName("l0-test-lb_id")

	describeLoadBalancersOutput := &elb.DescribeLoadBalancersOutput{}
	describeLoadBalancersOutput.SetLoadBalancerDescriptions([]*elb.LoadBalancerDescription{lb})

	gomock.InOrder(
		mockAWS.ELB.EXPECT().
			DescribeLoadBalancers(gomock.Any()).
			Return(describeLoadBalancersOutput, nil),
		mockAWS.ELB.EXPECT().
			DescribeLoadBalancers(gomock.Any()).
			Return(&elb.DescribeLoadBalancersOutput{}, awserr.New("LoadBalancerNotFound", "", nil)),
	)

	mockAWS.IAM.EXPECT().
		DeleteRolePolicy(gomock.Any()).
		Return(&iam.DeleteRolePolicyOutput{}, nil)

	mockAWS.IAM.EXPECT().
		DeleteRole(gomock.Any()).
		Return(&iam.DeleteRoleOutput{}, nil)

	readSGHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg")
	deleteSGHelper(mockAWS, "lb_sg")

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName("l0-test-lb_id-lb")
	securityGroup.SetGroupId("lb_sg")

	describeSecurityGroupsOutput := &ec2.DescribeSecurityGroupsOutput{}
	describeSecurityGroupsOutput.SetSecurityGroups([]*ec2.SecurityGroup{securityGroup})

	gomock.InOrder(
		mockAWS.EC2.EXPECT().
			DescribeSecurityGroups(gomock.Any()).
			Return(describeSecurityGroupsOutput, nil),
		mockAWS.EC2.EXPECT().
			DescribeSecurityGroups(gomock.Any()).
			Return(&ec2.DescribeSecurityGroupsOutput{}, nil),
	)

	mockAWS.ALB.EXPECT().
		DescribeTargetGroups(gomock.Any()).
		Return(nil, awserr.New("TargetGroupNotFound", "", nil))

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("lb_id"); err != nil {
		t.Fatal(err)
	}
}

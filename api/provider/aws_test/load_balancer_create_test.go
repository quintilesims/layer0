package aws_test

import (
	"testing"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/golang/mock/gomock"
	. "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	mock_aws "github.com/quintilesims/layer0/common/aws/mock_aws"
	models "github.com/quintilesims/layer0/common/models"
)

func TestLoadBalancer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock_aws.NewMockClient(ctrl)
	//tagStore := tag.NewMemoryStore()
	//apiConfig := config.NewContextAPIConfig(&cli.Context{})

	defer ctrl.Finish()

	lbp := mock_provider.NewMockLoadBalancerProvider(ctrl)

	environmentSGName := "env_id-env"
	filter := &ec2.Filter{}
	filter.SetName("group-name")
	filter.SetValues([]*string{aws.String(environmentSGName)})

	describeEnvSGInput := &ec2.DescribeSecurityGroupsInput{}
	describeEnvSGInput.SetFilters([]*ec2.Filter{filter})

	environmentSG := &ec2.DescribeSecurityGroupsOutput{}
	client.EC2.EXPECT().
		DescribeSecurityGroups(describeEnvSGInput).
		Return(environmentSG, nil)

	loadBalancerSGName := "lb_id-lb"
	createLBSGInput := &ec2.CreateSecurityGroupInput{}
	createLBSGInput.SetGroupName(loadBalancerSGName)
	createLBSGInput.SetDescription("SG for Layer0 load balancer lb_id")
	createLBSGInput.SetVpcId("vpc-id")

	client.EC2.EXPECT().
		CreateSecurityGroup(createLBSGInput)

	filter = &ec2.Filter{}
	filter.SetName("group-name")
	filter.SetValues([]*string{aws.String(loadBalancerSGName)})

	describeLBSGInput := &ec2.DescribeSecurityGroupsInput{}
	describeLBSGInput.SetFilters([]*ec2.Filter{filter})

	loadBalancerSG := &ec2.SecurityGroup{}
	client.EC2.EXPECT().
		DescribeSecurityGroups(describeLBSGInput).
		Return(loadBalancerSG, nil)

	AuthorizeSGIngressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
	AuthorizeSGIngressInput.SetGroupId("lb_id")
	AuthorizeSGIngressInput.SetCidrIp("0.0.0.0/0")
	AuthorizeSGIngressInput.SetIpProtocol("TCP")
	AuthorizeSGIngressInput.SetFromPort(80)
	AuthorizeSGIngressInput.SetToPort(80)

	client.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(AuthorizeSGIngressInput)

	CreateRoleInput := &iam.CreateRoleInput{}
	CreateRoleInput.SetRoleName("lb_id-lb")
	CreateRoleInput.SetAssumeRolePolicyDocument(DEFAULT_ASSUME_ROLE_POLICY)

	lbIAMRole := iam.CreateRoleOutput{}
	client.IAM.EXPECT().
		CreateRole(CreateRoleInput).
		Return(lbIAMRole, nil)

	putRolePolicyInput := &iam.PutRolePolicyInput{}
	putRolePolicyInput.SetPolicyName("lb_id-lb")
	putRolePolicyInput.SetRoleName("lb_id-lb")
	putRolePolicyInput.SetPolicyDocument("the_policy")

	client.IAM.EXPECT().
		PutRolePolicy(putRolePolicyInput)

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		IsPublic:         true,
		Ports:            []models.Port{},
		HealthCheck:      models.HealthCheck{},
	}

	if _, err := lbp.Create(req); err != nil {
		t.Fatal(err)
	}
}

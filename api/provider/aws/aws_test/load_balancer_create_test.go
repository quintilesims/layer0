package aws_test

import (
	"testing"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
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

	// Read the env security group
	environmentSGName := "env_id-env"
	filter := &ec2.Filter{}
	filter.SetName("group-name")
	filter.SetValues([]*string{aws.String(environmentSGName)})

	describeEnvSGInput := &ec2.DescribeSecurityGroupsInput{}
	describeEnvSGInput.SetFilters([]*ec2.Filter{filter})

	securityGroups := make([]*ec2.SecurityGroup, 1)
	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName(environmentSGName)
	securityGroups[0] = securityGroup
	environmentSG := &ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: securityGroups,
	}

	client.EC2.EXPECT().
		DescribeSecurityGroups(describeEnvSGInput).
		Return(environmentSG, nil)

	// Create LB security group
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

	// Read the LB security group
	describeLBSGInput := &ec2.DescribeSecurityGroupsInput{}
	describeLBSGInput.SetFilters([]*ec2.Filter{filter})

	loadBalancerSG := &ec2.DescribeSecurityGroupsOutput{}
	client.EC2.EXPECT().
		DescribeSecurityGroups(describeLBSGInput).
		Return(loadBalancerSG, nil)

	// Authorize SG as ingress to the LB
	AuthorizeSGIngressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
	AuthorizeSGIngressInput.SetGroupId("lb_id")
	AuthorizeSGIngressInput.SetCidrIp("0.0.0.0/0")
	AuthorizeSGIngressInput.SetIpProtocol("HTTPS")
	AuthorizeSGIngressInput.SetFromPort(443)
	AuthorizeSGIngressInput.SetToPort(443)

	client.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(AuthorizeSGIngressInput)

	// Create IAM role
	CreateRoleInput := &iam.CreateRoleInput{}
	CreateRoleInput.SetRoleName("lb_id-lb")
	CreateRoleInput.SetAssumeRolePolicyDocument(DEFAULT_ASSUME_ROLE_POLICY)

	lbIAMRole := &iam.CreateRoleOutput{}
	client.IAM.EXPECT().
		CreateRole(CreateRoleInput).
		Return(lbIAMRole, nil)

	// Add inline policy to IAM role
	putRolePolicyInput := &iam.PutRolePolicyInput{}
	putRolePolicyInput.SetPolicyName("lb_id-lb")
	putRolePolicyInput.SetRoleName("lb_id-lb")
	putRolePolicyInput.SetPolicyDocument("the_policy")

	client.IAM.EXPECT().
		PutRolePolicy(putRolePolicyInput)

	ports := []models.Port{
		models.Port{
			CertificateName: "cert",
			ContainerPort:   443,
			HostPort:        443,
			Protocol:        "https",
		},
	}
	certARNs := &iam.ListServerCertificatesOutput{}
	// Ports to listeners
	for _, port := range ports {
		listener := &elb.Listener{}
		listener.SetProtocol(port.Protocol)
		listener.SetLoadBalancerPort(port.HostPort)
		listener.SetInstancePort(port.ContainerPort)

		if port.CertificateName != "" {
			// If SSL Cert, lookup cert ARN
			client.IAM.EXPECT().
				ListServerCertificates(&iam.ListServerCertificatesInput{}).
				Return(certARNs, nil)

			listener.SetSSLCertificateId("arn")
		}
	}

	// Create load balancer request
	createLoadBalancerInput := &elb.CreateLoadBalancerInput{}
	createLoadBalancerInput.SetLoadBalancerName("lb_name")
	createLoadBalancerInput.SetScheme("internet-facing")
	createLoadBalancerInput.SetSecurityGroups([]*string{})
	createLoadBalancerInput.SetSubnets([]*string{})
	createLoadBalancerInput.SetListeners([]*elb.Listener{})

	client.ELB.EXPECT().
		CreateLoadBalancer(createLoadBalancerInput)

	healthCheck := &elb.HealthCheck{}
	// Update health check
	input := &elb.ConfigureHealthCheckInput{}
	input.SetLoadBalancerName("lb_name")
	input.SetHealthCheck(healthCheck)

	client.ELB.EXPECT().
		ConfigureHealthCheck(input)

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		IsPublic:         true,
		Ports:            ports,
		HealthCheck:      models.HealthCheck{},
	}

	lbp.EXPECT().
		Create(req).
		Return("lb_id", nil)

	if _, err := lbp.Create(req); err != nil {
		t.Fatal(err)
	}
}

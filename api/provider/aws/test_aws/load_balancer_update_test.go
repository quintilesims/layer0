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
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
)

func TestClassicLoadBalancerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	tags := models.Tags{
		{
			EntityID:   "lb_name",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name",
		},
		{
			EntityID:   "lb_name",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "lb_name",
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

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	requestPorts := []models.Port{
		models.Port{
			CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name",
			ContainerPort:  8080,
			HostPort:       8088,
			Protocol:       "http",
		},
		models.Port{
			CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name",
			ContainerPort:  4444,
			HostPort:       444,
			Protocol:       "https",
		},
	}

	requestHealthCheck := &models.HealthCheck{
		Target:             "HTTPS:444/path/to/site",
		Path:               "/",
		Interval:           15,
		Timeout:            10,
		HealthyThreshold:   5,
		UnhealthyThreshold: 4,
	}

	req := models.UpdateLoadBalancerRequest{
		Ports:       &requestPorts,
		HealthCheck: requestHealthCheck,
	}

	configureHealthCheckInput := &elb.ConfigureHealthCheckInput{}
	configureHealthCheckInput.SetLoadBalancerName("l0-test-lb_name")
	configureHealthCheckInput.SetHealthCheck(healthCheckHelper(requestHealthCheck))

	configureHealthCheckOutput := &elb.ConfigureHealthCheckOutput{}
	configureHealthCheckOutput.SetHealthCheck(healthCheckHelper(requestHealthCheck))

	mockAWS.ELB.EXPECT().
		ConfigureHealthCheck(configureHealthCheckInput).
		Return(configureHealthCheckOutput, nil)

	certificateARN := "arn:aws:iam::12345:server-certificate/crt_name"
	serverCertificateMetadata := &iam.ServerCertificateMetadata{}
	serverCertificateMetadata.SetArn(certificateARN)
	serverCertificateMetadata.SetServerCertificateName("cert")
	serverCertificateMetadataList := []*iam.ServerCertificateMetadata{serverCertificateMetadata}

	listServerCertificatesOutput := &iam.ListServerCertificatesOutput{}
	listServerCertificatesOutput.SetServerCertificateMetadataList(serverCertificateMetadataList)

	mockAWS.IAM.EXPECT().
		ListServerCertificates(&iam.ListServerCertificatesInput{}).
		Return(listServerCertificatesOutput, nil).
		AnyTimes()

	filter := &ec2.Filter{}
	filter.SetName("group-name")
	filter.SetValues([]*string{aws.String("l0-test-lb_name-lb")})

	input := &ec2.DescribeSecurityGroupsInput{}
	input.SetFilters([]*ec2.Filter{filter})

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName("l0-test-lb_name-lb")
	securityGroup.SetGroupId("lb_sg")
	securityGroup.IpPermissions = []*ec2.IpPermission{
		{
			FromPort: aws.Int64(config.DefaultLoadBalancerPort().HostPort),
		},
	}

	output := &ec2.DescribeSecurityGroupsOutput{}
	output.SetSecurityGroups([]*ec2.SecurityGroup{securityGroup})

	mockAWS.EC2.EXPECT().
		DescribeSecurityGroups(input).
		Return(output, nil)

	listenerDescription := &elb.ListenerDescription{}
	listenerDescription.SetListener(listenerHelper(config.DefaultLoadBalancerPort()))

	hc := config.DefaultLoadBalancerHealthCheck()
	lb := &elb.LoadBalancerDescription{}
	lb.SetLoadBalancerName("l0-test-lb_name")
	lb.SetHealthCheck(healthCheckHelper(&hc))
	lb.SetListenerDescriptions([]*elb.ListenerDescription{listenerDescription})

	describeLoadBalancersInput := &elb.DescribeLoadBalancersInput{}
	describeLoadBalancersInput.SetLoadBalancerNames([]*string{aws.String("l0-test-lb_name")})
	describeLoadBalancersInput.SetPageSize(1)

	describeLoadBalancersOutput := &elb.DescribeLoadBalancersOutput{}
	describeLoadBalancersOutput.SetLoadBalancerDescriptions([]*elb.LoadBalancerDescription{lb})

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(describeLoadBalancersInput).
		Return(describeLoadBalancersOutput, nil)

	revokeIngressInput := revokeSGIngressHelper(config.DefaultLoadBalancerPort())
	revokeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		RevokeSecurityGroupIngress(revokeIngressInput).
		Return(&ec2.RevokeSecurityGroupIngressOutput{}, nil)

	port := int64(config.DefaultLoadBalancerPort().HostPort)
	deleteLoadBalancerListenersInput := &elb.DeleteLoadBalancerListenersInput{}
	deleteLoadBalancerListenersInput.SetLoadBalancerName("l0-test-lb_name")
	deleteLoadBalancerListenersInput.SetLoadBalancerPorts([]*int64{&port})

	mockAWS.ELB.EXPECT().
		DeleteLoadBalancerListeners(deleteLoadBalancerListenersInput).
		Return(&elb.DeleteLoadBalancerListenersOutput{}, nil)

	listener1 := listenerHelper(requestPorts[0])
	listener1.SetSSLCertificateId(certificateARN)
	listener2 := listenerHelper(requestPorts[1])
	listener2.SetSSLCertificateId(certificateARN)

	createLoadBalancerListenersInput := &elb.CreateLoadBalancerListenersInput{}
	createLoadBalancerListenersInput.SetLoadBalancerName("l0-test-lb_name")
	createLoadBalancerListenersInput.SetListeners([]*elb.Listener{listener1, listener2})

	mockAWS.ELB.EXPECT().
		CreateLoadBalancerListeners(createLoadBalancerListenersInput).
		Return(&elb.CreateLoadBalancerListenersOutput{}, nil)

	authorizeIngressInput := authorizeSGIngressHelper(requestPorts[0])
	authorizeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(authorizeIngressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	authorizeIngressInput = authorizeSGIngressHelper(requestPorts[1])
	authorizeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(authorizeIngressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Update("lb_name", req); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationLoadBalancerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	tags := models.Tags{
		{
			EntityID:   "lb_name",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name",
		},
		{
			EntityID:   "lb_name",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "lb_name",
			EntityType: "load_balancer",
			Key:        "type",
			Value:      string(models.ApplicationLoadBalancerType),
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	requestPorts := []models.Port{
		models.Port{
			CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name",
			ContainerPort:  8080,
			HostPort:       8088,
			Protocol:       "http",
		},
		models.Port{
			CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name",
			ContainerPort:  4444,
			HostPort:       444,
			Protocol:       "https",
		},
	}

	requestHealthCheck := &models.HealthCheck{
		Target:             "HTTPS:444/path/to/site",
		Path:               "/",
		Interval:           15,
		Timeout:            10,
		HealthyThreshold:   5,
		UnhealthyThreshold: 4,
	}

	req := models.UpdateLoadBalancerRequest{
		Ports:       &requestPorts,
		HealthCheck: requestHealthCheck,
	}

	updateHealthCheckInput := healthCheckTargetGroupHelper(requestHealthCheck)
	updateHealthCheckInput.SetTargetGroupArn("arn:target_name")
	mockAWS.ALB.EXPECT().
		ModifyTargetGroup(updateHealthCheckInput).
		Return(&alb.ModifyTargetGroupOutput{}, nil)

	describeClassicLBInput := &elb.DescribeLoadBalancersInput{}
	describeClassicLBInput.SetLoadBalancerNames([]*string{aws.String("l0-test-lb_name")})
	describeClassicLBInput.SetPageSize(1)

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(describeClassicLBInput).
		Return(nil, awserr.New("LoadBalancerNotFound", "", nil))

	describeAppLBInput := &alb.DescribeLoadBalancersInput{}
	describeAppLBInput.SetNames([]*string{aws.String("l0-test-lb_name")})
	describeAppLBOutput := &alb.DescribeLoadBalancersOutput{
		LoadBalancers: []*alb.LoadBalancer{
			{
				LoadBalancerArn: aws.String("arn:l0-test-lb_id"),
			},
		},
	}

	mockAWS.ALB.EXPECT().
		DescribeLoadBalancers(describeAppLBInput).
		Return(describeAppLBOutput, nil)

	describeTGInput := &alb.DescribeTargetGroupsInput{}
	describeTGInput.SetNames([]*string{aws.String("l0-test-lb_name")})
	describeTGOutput := &alb.DescribeTargetGroupsOutput{}
	describeTGOutput.SetTargetGroups([]*alb.TargetGroup{
		{
			TargetGroupArn: aws.String("arn:target_name"),
		},
	})

	mockAWS.ALB.EXPECT().
		DescribeTargetGroups(describeTGInput).
		Return(describeTGOutput, nil).
		Times(2)

	descListenersInput := &alb.DescribeListenersInput{}
	descListenersInput.LoadBalancerArn = aws.String("arn:l0-test-lb_id")
	descListenersInput.SetPageSize(10)
	fnListListenersPage := func(
		input *alb.DescribeListenersInput,
		fn func(o *alb.DescribeListenersOutput, lastPage bool) bool) error {
		output := &alb.DescribeListenersOutput{}
		fn(output, true)
		return nil
	}

	mockAWS.ALB.EXPECT().
		DescribeListenersPages(descListenersInput, gomock.Any()).
		Do(fnListListenersPage).
		Return(nil)

	for _, p := range requestPorts {
		createListenerInput := &alb.CreateListenerInput{}
		createListenerInput.SetPort(p.HostPort)
		createListenerInput.SetProtocol(p.Protocol)
		createListenerInput.LoadBalancerArn = aws.String("arn:l0-test-lb_id")
		createListenerInput.SetDefaultActions([]*alb.Action{
			{
				TargetGroupArn: aws.String("arn:target_name"),
				Type:           aws.String(alb.ActionTypeEnumForward),
			},
		})

		if p.CertificateARN != "" {
			createListenerInput.SetCertificates([]*alb.Certificate{
				{
					CertificateArn: aws.String(p.CertificateARN),
				},
			})
		}

		mockAWS.ALB.EXPECT().
			CreateListener(createListenerInput).
			Return(&alb.CreateListenerOutput{}, nil)
	}

	filter := &ec2.Filter{}
	filter.SetName("group-name")
	filter.SetValues([]*string{aws.String("l0-test-lb_name-lb")})

	input := &ec2.DescribeSecurityGroupsInput{}
	input.SetFilters([]*ec2.Filter{filter})

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName("l0-test-lb_name-lb")
	securityGroup.SetGroupId("lb_sg")
	securityGroup.IpPermissions = []*ec2.IpPermission{
		{
			FromPort: aws.Int64(config.DefaultLoadBalancerPort().HostPort),
		},
	}

	output := &ec2.DescribeSecurityGroupsOutput{}
	output.SetSecurityGroups([]*ec2.SecurityGroup{securityGroup})

	mockAWS.EC2.EXPECT().
		DescribeSecurityGroups(input).
		Return(output, nil)

	revokeIngressInput := revokeSGIngressHelper(config.DefaultLoadBalancerPort())
	revokeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		RevokeSecurityGroupIngress(revokeIngressInput).
		Return(&ec2.RevokeSecurityGroupIngressOutput{}, nil)

	authorizeIngressInput := authorizeSGIngressHelper(requestPorts[0])
	authorizeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(authorizeIngressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	authorizeIngressInput = authorizeSGIngressHelper(requestPorts[1])
	authorizeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(authorizeIngressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Update("lb_name", req); err != nil {
		t.Fatal(err)
	}
}

package cloudformation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/quintilesims/layer0/common/aws/provider"
)

type Provider interface {
	CreateStack(string, string, []*Parameter) (stackID string, err error)
	UpdateStack(string, string, []*Parameter) (stackID string, err error)
	DeleteStack(stackName string) error
	DescribeStack(stackName string) (*Stack, error)
}

type CloudFormation struct {
	credProvider provider.CredProvider
	region       string
	Connect      func() (cloudFormationInternal, error)
}

type cloudFormationInternal interface {
	CreateStack(input *cloudformation.CreateStackInput) (output *cloudformation.CreateStackOutput, err error)
	UpdateStack(input *cloudformation.UpdateStackInput) (output *cloudformation.UpdateStackOutput, err error)
	DeleteStack(*cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
	DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
}

type Stack struct {
	*cloudformation.Stack
}

func NewStack() *Stack {
	return &Stack{&cloudformation.Stack{}}
}

type Parameter struct {
	*cloudformation.Parameter
}

func NewParameter(key, value string) *Parameter {
	return &Parameter{
		&cloudformation.Parameter{
			ParameterKey:   aws.String(key),
			ParameterValue: aws.String(value),
		},
	}
}

func Connect(credProvider provider.CredProvider, region string) (cloudFormationInternal, error) {
	connection, err := provider.GetCloudFormationConnection(credProvider, region)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func NewCloudFormation(credProvider provider.CredProvider, region string) (Provider, error) {
	cloudformation := CloudFormation{
		credProvider: credProvider,
		region:       region,
		Connect:      func() (cloudFormationInternal, error) { return Connect(credProvider, region) },
	}

	_, err := cloudformation.Connect()
	if err != nil {
		return nil, err
	}

	return &cloudformation, nil
}

func (this *CloudFormation) CreateStack(name, templateBody string, params []*Parameter) (stackID string, err error) {
	parameters := []*cloudformation.Parameter{}
	for _, param := range params {
		parameters = append(parameters, param.Parameter)
	}

	input := &cloudformation.CreateStackInput{
		StackName:    aws.String(name),
		TemplateBody: aws.String(templateBody),
		Parameters:   parameters,
	}

	connection, err := this.Connect()
	if err != nil {
		return "", err
	}

	output, err := connection.CreateStack(input)
	if err != nil {
		return
	}

	stackID = *output.StackId
	return
}

func (this *CloudFormation) UpdateStack(name, templateBody string, params []*Parameter) (stackID string, err error) {
	parameters := []*cloudformation.Parameter{}
	for _, param := range params {
		parameters = append(parameters, param.Parameter)
	}

	input := &cloudformation.UpdateStackInput{
		StackName:    aws.String(name),
		TemplateBody: aws.String(templateBody),
		Parameters:   parameters,
	}

	connection, err := this.Connect()
	if err != nil {
		return "", err
	}

	output, err := connection.UpdateStack(input)
	if err != nil {
		return
	}

	stackID = *output.StackId
	return
}

func (this *CloudFormation) DescribeStack(stackName string) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.DescribeStacks(input)
	if err != nil {
		return nil, err
	}

	var stack *Stack
	if output.Stacks != nil && len(output.Stacks) > 0 {
		stack = &Stack{output.Stacks[0]}
	}

	return stack, nil
}

func (this *CloudFormation) DeleteStack(stackName string) error {
	input := &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteStack(input)

	return err
}

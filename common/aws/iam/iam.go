package iam

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	//	"github.com/quintilesims/layer0/api/config"
	"strings"

	"github.com/quintilesims/layer0/common/aws/provider"
)

type Provider interface {
	UploadServerCertificate(string, string, string, string, *string) (*ServerCertificateMetadata, error)
	ListCertificates() ([]*ServerCertificateMetadata, error)
	GetUser(username *string) (*User, error)
	DeleteServerCertificate(certName string) error
	CreateRole(roleName, servicePrincipal string) (*Role, error)
	GetRole(roleName string) (*Role, error)
	PutRolePolicy(roleName, policy string) error
	GetAccountId() (string, error)
	DeleteRole(roleName string) error
	DeleteRolePolicy(roleName, policyName string) error
	ListRolePolicies(roleName string) ([]*string, error)
	ListRoles() ([]*string, error)
}

type iamInternal interface {
	UploadServerCertificate(input *iam.UploadServerCertificateInput) (output *iam.UploadServerCertificateOutput, err error)
	ListServerCertificates(*iam.ListServerCertificatesInput) (*iam.ListServerCertificatesOutput, error)
	DeleteServerCertificate(input *iam.DeleteServerCertificateInput) (*iam.DeleteServerCertificateOutput, error)
	GetUser(input *iam.GetUserInput) (*iam.GetUserOutput, error)
	PutRolePolicy(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error)
	CreateRole(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error)
	GetRole(*iam.GetRoleInput) (*iam.GetRoleOutput, error)
	DeleteRole(*iam.DeleteRoleInput) (*iam.DeleteRoleOutput, error)
	DeleteRolePolicy(*iam.DeleteRolePolicyInput) (*iam.DeleteRolePolicyOutput, error)
	ListRolePolicies(*iam.ListRolePoliciesInput) (*iam.ListRolePoliciesOutput, error)
	ListRoles(*iam.ListRolesInput) (*iam.ListRolesOutput, error)
}

type IAM struct {
	credProvider provider.CredProvider
	region       string
	Connect      func() (iamInternal, error)
}

type ServerCertificateMetadata struct {
	*iam.ServerCertificateMetadata
}

func NewServerCertificateMetadata(name, arn string) *ServerCertificateMetadata {
	return &ServerCertificateMetadata{
		&iam.ServerCertificateMetadata{
			ServerCertificateName: aws.String(name),
			Arn: aws.String(arn),
		},
	}
}

type User struct {
	*iam.User
}

func NewUser() *User {
	return &User{&iam.User{}}
}

type Role struct {
	*iam.Role
}

func Connect(credProvider provider.CredProvider, region string) (iamInternal, error) {
	connection, err := provider.GetIAMConnection(credProvider, region)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func NewIAM(credProvider provider.CredProvider, region string) (Provider, error) {
	iam := IAM{
		credProvider: credProvider,
		region:       region,
		Connect:      func() (iamInternal, error) { return Connect(credProvider, region) },
	}

	_, err := iam.Connect()
	if err != nil {
		return nil, err
	}

	return &iam, nil
}

func (this *IAM) UploadServerCertificate(name, path, body, pk string, optionalChain *string) (*ServerCertificateMetadata, error) {
	input := &iam.UploadServerCertificateInput{
		ServerCertificateName: aws.String(name),
		CertificateBody:       aws.String(body),
		CertificateChain:      optionalChain,
		PrivateKey:            aws.String(pk),
		Path:                  aws.String(path),
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.UploadServerCertificate(input)
	if err != nil {
		return nil, err
	}

	metadata := output.ServerCertificateMetadata
	return &ServerCertificateMetadata{metadata}, nil
}

func (this *IAM) ListCertificates() ([]*ServerCertificateMetadata, error) {
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.ListServerCertificates(&iam.ListServerCertificatesInput{})
	if err != nil {
		return nil, err
	}

	certs := []*ServerCertificateMetadata{}
	for _, metadata := range output.ServerCertificateMetadataList {
		certs = append(certs, &ServerCertificateMetadata{metadata})
	}

	return certs, nil
}

func (this *IAM) GetUser(username *string) (*User, error) {
	input := &iam.GetUserInput{
		UserName: username,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.GetUser(input)
	if err != nil {
		return nil, err
	}

	var user *User
	if output.User != nil {
		user = &User{output.User}
	}

	return user, nil
}

func (this *IAM) DeleteServerCertificate(certName string) error {
	input := &iam.DeleteServerCertificateInput{
		ServerCertificateName: aws.String(certName),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteServerCertificate(input)
	return err
}

func (this *IAM) CreateRole(roleName, servicePrincipal string) (*Role, error) {
	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(fmt.Sprintf(`{"Version":"2008-10-17","Statement":[{"Sid":"","Effect":"Allow","Principal":{"Service":["%s"]},"Action":["sts:AssumeRole"]}]}`, servicePrincipal)),
		RoleName:                 &roleName,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	out, err := connection.CreateRole(input)
	if err != nil {
		return nil, err
	}

	return &Role{out.Role}, nil
}

func (this *IAM) GetRole(roleName string) (*Role, error) {
	input := &iam.GetRoleInput{
		RoleName: &roleName,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.GetRole(input)
	if err != nil {
		return nil, err
	}

	return &Role{output.Role}, nil
}

func (this *IAM) DeleteRole(roleName string) error {
	input := &iam.DeleteRoleInput{
		RoleName: &roleName,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteRole(input)
	return err
}

func (this *IAM) DeleteRolePolicy(roleName, policyName string) error {
	input := &iam.DeleteRolePolicyInput{
		RoleName:   &roleName,
		PolicyName: &policyName,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteRolePolicy(input)
	return err
}

func (this *IAM) ListRolePolicies(roleName string) ([]*string, error) {
	input := &iam.ListRolePoliciesInput{
		RoleName: &roleName,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.ListRolePolicies(input)
	if err != nil {
		return nil, err
	}
	return output.PolicyNames, nil
}

func (this *IAM) ListRoles() ([]*string, error) {
	input := &iam.ListRolesInput{}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.ListRoles(input)
	if err != nil {
		return nil, err
	}

	roles := []*string{}
	for _, role := range output.Roles {
		roles = append(roles, role.RoleName)
	}

	return roles, nil
}

func (this *IAM) PutRolePolicy(roleName, policy string) error {
	input := &iam.PutRolePolicyInput{
		PolicyName:     &roleName,
		PolicyDocument: &policy,
		RoleName:       &roleName,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.PutRolePolicy(input)
	return err
}

func (this *IAM) GetAccountId() (string, error) {
	connection, err := this.Connect()
	if err != nil {
		return "", err
	}

	out, err := connection.GetUser(nil)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to get current IAM user: %s", err.Error())
	}
	// Sample ARN: "arn:aws:iam::123456789012:user/layer0/l0/bootstrap-user-user-ABCDEFGHIJKL"
	return strings.Split(*out.User.Arn, ":")[4], nil
}

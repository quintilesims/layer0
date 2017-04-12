package context

import (
	"encoding/base64"
	"fmt"
	"github.com/quintilesims/layer0/common/aws/iam"
	"github.com/quintilesims/layer0/common/aws/provider"
	"math/rand"
	"strings"
	"time"
)

type VariableSchema struct {
	Name            string
	DisplayName     string
	Default         string
	Description     string
	IsCalculated    bool
	AlwaysCalculate bool
	Calculate       func(*Context) (string, error)
	Validate        func(*Context, string) error
}

var TerraformVariables = []VariableSchema{
	VariableSchema{
		Name:        "access_key",
		Description: "The aws_access_key_id for your Layer0",
		DisplayName: "AWS Access Key ID",
	},
	VariableSchema{
		Name:        "secret_key",
		Description: "The aws_secret_access_key for your Layer0",
		DisplayName: "AWS Secret Access Key",
	},
	VariableSchema{
		Name:        "region",
		Default:     "us-west-2",
		DisplayName: "AWS Region",
		Description: "The aws_region for your Layer0",
		Validate:    validateRegion,
	},
	VariableSchema{
		Name:         "l0_prefix",
		IsCalculated: true,
		Calculate:    func(c *Context) (string, error) { return c.Instance, nil },
	},
	VariableSchema{
		Name:         "account_id",
		IsCalculated: true,
		Calculate:    getAccountID,
	},
	VariableSchema{
		Name:        "key_pair",
		Description: "The name of an EC2 key pair",
		DisplayName: "Key Pair Name",
	},
	VariableSchema{
		Name:         "db_master_username",
		IsCalculated: true,
		Calculate:    func(c *Context) (string, error) { return "layer0_master", nil },
	},
	VariableSchema{
		Name:         "s3_bucket",
		IsCalculated: true,
		Calculate:    getS3Bucket,
	},
	VariableSchema{
		Name:         "api_auth_token",
		IsCalculated: true,
		Calculate:    getAuthToken,
	},
	VariableSchema{
		Name:        "api_docker_image",
		Description: "Layer0 API image name (like quintilesims/l0-api)",
		Default:     "quintilesims/l0-api",
	},
	VariableSchema{
		Name:         "api_docker_image_tag",
		Description:  "Layer0 API image tag (like v0.6.0)",
		IsCalculated: true,
		Calculate:    getImageTag,
	},
	VariableSchema{
		Name:         "runner_docker_image_tag",
		Description:  "Layer0 Runner image tag (like v0.6.0)",
		IsCalculated: true,
		Calculate:    getImageTag,
	},
}

func GetTerraformVariable(name string) (VariableSchema, bool) {
	for _, variable := range TerraformVariables {
		if variable.Name == name {
			return variable, true
		}
	}

	return VariableSchema{}, false
}

func validateRegion(c *Context, val string) error {
	regions := []string{
		"us-west-2",
		"us-west-1",
		"us-east-1",
		"eu-west-1",
	}

	for _, r := range regions {
		if val == r {
			return nil
		}
	}

	return fmt.Errorf("Region '%s' is not a supported region! Supported regions: %v", val, regions)
}

func getAWSVars(c *Context, requester string) (string, string, string, error) {
	accessKey, ok := c.TerraformVars["access_key"]
	if !ok {
		return "", "", "", fmt.Errorf("Variable 'access_key' not set prior to calculating '%s'", requester)
	}

	secretKey, ok := c.TerraformVars["secret_key"]
	if !ok {
		return "", "", "", fmt.Errorf("Variable 'secret_key' not set prior to calculating '%s'", requester)
	}

	region, ok := c.TerraformVars["region"]
	if !ok {
		// use default region
		region = "us-west-2"
	}

	return accessKey, secretKey, region, nil
}

func getAccountID(c *Context) (string, error) {
	accessKey, secretKey, region, err := getAWSVars(c, "account_id")
	if err != nil {
		return "", err
	}

	creds := provider.NewExplicitCredProvider(accessKey, secretKey)
	conn, err := iam.NewIAM(creds, region)
	if err != nil {
		return "", fmt.Errorf("Failed to connect to IAM: %s", err.Error())
	}

	user, err := conn.GetUser(nil)
	if err != nil {
		return "", fmt.Errorf("Failed to get current IAM user: %s", err.Error())
	}

	// Sample ARN: "arn:aws:iam::123456789012:user/layer0/l0/bootstrap-user-user-ABCDEFGHIJKL"
	return strings.Split(*user.Arn, ":")[4], nil
}

func getRandomPassword(c *Context) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 10)

	for i := 0; i < 10; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result), nil
}

func getS3Bucket(c *Context) (string, error) {
	accountID, err := getAccountID(c)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("layer0-%s-%s", c.Instance, accountID), nil
}

func getAuthToken(c *Context) (string, error) {
	password, err := getRandomPassword(c)
	if err != nil {
		return "", err
	}

	data := []byte(fmt.Sprintf("layer0:%s", password))
	return base64.StdEncoding.EncodeToString(data), nil
}

func getImageTag(c *Context) (string, error) {
	if val, ok := c.TerraformVars["setup_version"]; ok {
		if strings.Contains(val, "unset") {
			fmt.Println("[WARNING] Developer API setup_version unset, using 'latest' image")
			return "latest", nil
		} else {
			return val, nil
		}
	}

	return "", fmt.Errorf("Key 'setup_version' not set in context TerraformVars")
}

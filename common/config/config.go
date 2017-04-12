package config

import (
	"fmt"
	"os"
	"strings"
)

// IMPORTANT!
// The environment variables represented as constants here should
// always line up with the environment variables in setup/container_definitions.json
const (
	AWS_ACCOUNT_ID                  = "LAYER0_AWS_ACCOUNT_ID"
	AWS_ACCESS_KEY_ID               = "LAYER0_AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY           = "LAYER0_AWS_SECRET_ACCESS_KEY"
	AWS_VPC_ID                      = "LAYER0_AWS_VPC_ID"
	AWS_PRIVATE_SUBNETS             = "LAYER0_AWS_PRIVATE_SUBNETS"
	AWS_PUBLIC_SUBNETS              = "LAYER0_AWS_PUBLIC_SUBNETS"
	AWS_ECS_ROLE                    = "LAYER0_AWS_ECS_ROLE"
	AWS_KEY_PAIR                    = "LAYER0_AWS_KEY_PAIR"
	AWS_S3_BUCKET                   = "LAYER0_AWS_S3_BUCKET"
	AWS_ECS_AGENT_SECURITY_GROUP_ID = "LAYER0_AWS_ECS_AGENT_SECURITY_GROUP_ID"
	AWS_ECS_INSTANCE_PROFILE        = "LAYER0_AWS_ECS_INSTANCE_PROFILE"
	AWS_DYNAMO_TAG_TABLE_NAME           = "LAYER0_AWS_DYNAMO_TAG_TABLE_NAME"
	 AWS_DYNAMO_JOB_TABLE_NAME           = "LAYER0_AWS_DYNAMO_JOB_TABLE_NAME"
	JOB_ID                          = "LAYER0_JOB_ID"
	AWS_LINUX_SERVICE_AMI           = "LAYER0_AWS_LINUX_SERVICE_AMI"
	AWS_WINDOWS_SERVICE_AMI         = "LAYER0_AWS_WINDOWS_SERVICE_AMI"
	AWS_REGION                      = "LAYER0_AWS_REGION"
	AUTH_TOKEN                      = "LAYER0_AUTH_TOKEN"
	API_ENDPOINT                    = "LAYER0_API_ENDPOINT"
	API_PORT                        = "LAYER0_API_PORT"
	API_LOG_LEVEL                   = "LAYER0_API_LOG_LEVEL"
	PREFIX                          = "LAYER0_PREFIX"
	RUNNER_LOG_LEVEL                = "LAYER0_RUNNER_LOG_LEVEL"
	RUNNER_VERSION_TAG              = "LAYER0_RUNNER_VERSION_TAG"
	SKIP_SSL_VERIFY                 = "LAYER0_SKIP_SSL_VERIFY"
	SKIP_VERSION_VERIFY             = "LAYER0_SKIP_VERSION_VERIFY"
	DB_CONNECTION                   = "LAYER0_DB_CONNECTION"
	DB_NAME                         = "LAYER0_DB_NAME"
	TEST_AWS_TAG_DYNAMO_TABLE_NAME      = "LAYER0_TEST_AWS_TAG_DYNAMO_TABLE_NAME"
	 TEST_AWS_JOB_DYNAMO_TABLE_NAME      = "LAYER0_TEST_AWS_JOB_DYNAMO_TABLE_NAME"
)

// defaults
// bGF5ZXIwOm5vaGF4cGx6 = layer0:nohaxplz, base64 encoded (basic http auth)
const (
	DEFAULT_DB_CONNECTION = "layer0:nohaxplz@tcp(127.0.0.1:3306)/"
	DEFAULT_AUTH_TOKEN    = "bGF5ZXIwOm5vaGF4cGx6"
	DEFAULT_API_ENDPOINT  = "http://localhost:9090/"
	DEFAULT_API_PORT      = "9090"
	DEFAULT_AWS_REGION    = "us-west-2"
)

// api resource tags
const (
	API_ENVIRONMENT_ID     = "api"
	API_ENVIRONMENT_NAME   = "api"
	API_LOAD_BALANCER_ID   = "api"
	API_LOAD_BALANCER_NAME = "api"
	API_SERVICE_ID         = "api"
	API_SERVICE_NAME       = "api"
)

var RequiredAPIVariables = []string{
	AWS_ACCOUNT_ID,
	AWS_ACCESS_KEY_ID,
	AWS_SECRET_ACCESS_KEY,
	AWS_VPC_ID,
	AWS_PRIVATE_SUBNETS,
	AWS_PUBLIC_SUBNETS,
	AWS_ECS_ROLE,
	AWS_KEY_PAIR,
	AWS_S3_BUCKET,
	AWS_ECS_AGENT_SECURITY_GROUP_ID,
	AWS_ECS_INSTANCE_PROFILE,
	AWS_LINUX_SERVICE_AMI,
	AWS_WINDOWS_SERVICE_AMI,
}

var RequiredCLIVariables = []string{}

var RequiredRunnerVariables = []string{
	AWS_ACCESS_KEY_ID,
	AWS_SECRET_ACCESS_KEY,
	AWS_VPC_ID,
	AWS_PRIVATE_SUBNETS,
	AWS_PUBLIC_SUBNETS,
}

func Validate(required []string) error {
	for _, key := range required {
		if os.Getenv(key) == "" {
			return fmt.Errorf("Required environment variable '%s' not set", key)
		}
	}

	return nil
}

func get(key string) string {
	return os.Getenv(key)
}

func getOr(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func DBName() string {
	return getOr(DB_NAME, fmt.Sprintf("layer0_%s", Prefix()))
}

func DBConnection() string {
	return getOr(DB_CONNECTION, DEFAULT_DB_CONNECTION)
}

var apiVersion string

func SetAPIVersion(version string) {
	apiVersion = version
}

func APIVersion() string {
	return apiVersion
}

var cliVersion string

func CLIVersion() string {
	return cliVersion
}

func SetCLIVersion(version string) {
	cliVersion = version
}

func AWSLogGroupID() string {
	return fmt.Sprintf("l0-%s", Prefix())
}

func AWSLinuxServiceAMI() string {
	return get(AWS_LINUX_SERVICE_AMI)
}

func AWSWindowsServiceAMI() string {
	return get(AWS_WINDOWS_SERVICE_AMI)
}

func AWSAccountID() string {
	return get(AWS_ACCOUNT_ID)
}

func AWSAccessKey() string {
	return get(AWS_ACCESS_KEY_ID)
}

func AWSSecretKey() string {
	return get(AWS_SECRET_ACCESS_KEY)
}

func AWSRegion() string {
	return getOr(AWS_REGION, DEFAULT_AWS_REGION)
}

func AWSVPCID() string {
	return get(AWS_VPC_ID)
}

func AWSPrivateSubnets() string {
	v := get(AWS_PRIVATE_SUBNETS)
	return strings.Replace(v, " ", "", -1)
}

func AWSPublicSubnets() string {
	v := get(AWS_PUBLIC_SUBNETS)
	return strings.Replace(v, " ", "", -1)
}

func AWSECSRole() string {
	return get(AWS_ECS_ROLE)
}

func AWSKeyPair() string {
	return get(AWS_KEY_PAIR)
}

func AWSS3Bucket() string {
	return get(AWS_S3_BUCKET)
}

func AuthToken() string {
	return getOr(AUTH_TOKEN, DEFAULT_AUTH_TOKEN)
}

func APIEndpoint() string {
	return getOr(API_ENDPOINT, DEFAULT_API_ENDPOINT)
}

func APIPort() string {
	return getOr(API_PORT, DEFAULT_API_PORT)
}

func APILogLevel() string {
	return getOr(API_LOG_LEVEL, "1")
}

func DynamoTagTableName() string {
	other := fmt.Sprintf("l0-%s-tags", Prefix())
	return getOr(AWS_DYNAMO_TAG_TABLE_NAME, other)
}

func TestDynamoTagTableName() string {
	return get(TEST_AWS_TAG_DYNAMO_TABLE_NAME)
}

func DynamoJobTableName() string {
        other := fmt.Sprintf("l0-%s-jobs", Prefix())
        return getOr(AWS_DYNAMO_JOB_TABLE_NAME, other)
}

func TestDynamoJobTableName() string {
        return get(TEST_AWS_JOB_DYNAMO_TABLE_NAME)
}

func Prefix() string {
	return getOr(PREFIX, "l0")
}

func RunnerLogLevel() string {
	return getOr(RUNNER_LOG_LEVEL, "1")
}

func RunnerVersionTag() string {
	return getOr(RUNNER_VERSION_TAG, "latest")
}

func AWSAgentGroupID() string {
	return get(AWS_ECS_AGENT_SECURITY_GROUP_ID)
}

func AWSECSInstanceProfile() string {
	return get(AWS_ECS_INSTANCE_PROFILE)
}

func ShouldVerifySSL() bool {
	val := strings.ToLower(getOr(SKIP_SSL_VERIFY, ""))
	if val == "1" || val == "true" {
		return false
	}

	return true
}

func ShouldVerifyVersion() bool {
	val := strings.ToLower(getOr(SKIP_VERSION_VERIFY, ""))
	if val == "1" || val == "true" {
		return false
	}

	return true
}

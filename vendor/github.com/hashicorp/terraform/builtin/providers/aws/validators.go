package aws

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/schema"
)

func validateRdsId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[0-9a-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q", k))
	}
	if !regexp.MustCompile(`^[a-z]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"first character of %q must be a letter", k))
	}
	if regexp.MustCompile(`--`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot contain two consecutive hyphens", k))
	}
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen", k))
	}
	return
}

func validateElastiCacheClusterId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 20) {
		errors = append(errors, fmt.Errorf(
			"%q must contain from 1 to 20 alphanumeric characters or hyphens", k))
	}
	if !regexp.MustCompile(`^[0-9a-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q", k))
	}
	if !regexp.MustCompile(`^[a-z]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"first character of %q must be a letter", k))
	}
	if regexp.MustCompile(`--`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot contain two consecutive hyphens", k))
	}
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen", k))
	}
	return
}

func validateASGScheduleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(awsAutoscalingScheduleTimeLayout, value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as iso8601 Timestamp Format", value))
	}

	return
}

// validateTagFilters confirms the "value" component of a tag filter is one of
// AWS's three allowed types.
func validateTagFilters(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "KEY_ONLY" && value != "VALUE_ONLY" && value != "KEY_AND_VALUE" {
		errors = append(errors, fmt.Errorf(
			"%q must be one of \"KEY_ONLY\", \"VALUE_ONLY\", or \"KEY_AND_VALUE\"", k))
	}
	return
}

func validateDbParamGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[0-9a-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q", k))
	}
	if !regexp.MustCompile(`^[a-z]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"first character of %q must be a letter", k))
	}
	if regexp.MustCompile(`--`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot contain two consecutive hyphens", k))
	}
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen", k))
	}
	if len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be greater than 255 characters", k))
	}
	return

}

func validateStreamViewType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	viewTypes := map[string]bool{
		"KEYS_ONLY":          true,
		"NEW_IMAGE":          true,
		"OLD_IMAGE":          true,
		"NEW_AND_OLD_IMAGES": true,
	}

	if !viewTypes[value] {
		errors = append(errors, fmt.Errorf("%q be a valid DynamoDB StreamViewType", k))
	}
	return
}

func validateElbName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[0-9A-Za-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alphanumeric characters and hyphens allowed in %q: %q",
			k, value))
	}
	if len(value) > 32 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 32 characters: %q", k, value))
	}
	if regexp.MustCompile(`^-`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot begin with a hyphen: %q", k, value))
	}
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen: %q", k, value))
	}
	return

}

func validateEcrRepositoryName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 {
		errors = append(errors, fmt.Errorf(
			"%q must be at least 2 characters long: %q", k, value))
	}
	if len(value) > 256 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 256 characters: %q", k, value))
	}

	// http://docs.aws.amazon.com/AmazonECR/latest/APIReference/API_CreateRepository.html
	pattern := `^(?:[a-z0-9]+(?:[._-][a-z0-9]+)*/)*[a-z0-9]+(?:[._-][a-z0-9]+)*$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

func validateCloudWatchEventRuleName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}

	// http://docs.aws.amazon.com/AmazonCloudWatchEvents/latest/APIReference/API_PutRule.html
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

func validateMaxLength(length int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if len(value) > length {
			errors = append(errors, fmt.Errorf(
				"%q cannot be longer than %d characters: %q", k, length, value))
		}
		return
	}
}

func validateIntegerInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d: %d", k, min, value))
		}
		if value > max {
			errors = append(errors, fmt.Errorf(
				"%q cannot be higher than %d: %d", k, max, value))
		}
		return
	}
}

func validateCloudWatchEventTargetId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}

	// http://docs.aws.amazon.com/AmazonCloudWatchEvents/latest/APIReference/API_Target.html
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

func validateLambdaFunctionName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 140 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 140 characters: %q", k, value))
	}
	// http://docs.aws.amazon.com/lambda/latest/dg/API_AddPermission.html
	pattern := `^(arn:aws:lambda:)?([a-z]{2}-[a-z]+-\d{1}:)?(\d{12}:)?(function:)?([a-zA-Z0-9-_]+)(:(\$LATEST|[a-zA-Z0-9-_]+))?$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

func validateLambdaQualifier(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 128 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 128 characters: %q", k, value))
	}
	// http://docs.aws.amazon.com/lambda/latest/dg/API_AddPermission.html
	pattern := `^[a-zA-Z0-9$_]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

func validateLambdaPermissionAction(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	// http://docs.aws.amazon.com/lambda/latest/dg/API_AddPermission.html
	pattern := `^(lambda:[*]|lambda:[a-zA-Z]+|[*])$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

func validateAwsAccountId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	// http://docs.aws.amazon.com/lambda/latest/dg/API_AddPermission.html
	pattern := `^\d{12}$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't look like AWS Account ID (exactly 12 digits): %q",
			k, value))
	}

	return
}

func validateArn(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return
	}

	// http://docs.aws.amazon.com/lambda/latest/dg/API_AddPermission.html
	pattern := `^arn:aws:([a-zA-Z0-9\-])+:([a-z]{2}-[a-z]+-\d{1})?:(\d{12})?:(.*)$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't look like a valid ARN (%q): %q",
			k, pattern, value))
	}

	return
}

func validatePolicyStatementId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 100 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 100 characters: %q", k, value))
	}

	// http://docs.aws.amazon.com/lambda/latest/dg/API_AddPermission.html
	pattern := `^[a-zA-Z0-9-_]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't look like a valid statement ID (%q): %q",
			k, pattern, value))
	}

	return
}

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}

func validateHTTPMethod(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	validMethods := map[string]bool{
		"ANY":     true,
		"DELETE":  true,
		"GET":     true,
		"HEAD":    true,
		"OPTIONS": true,
		"PATCH":   true,
		"POST":    true,
		"PUT":     true,
	}

	if _, ok := validMethods[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid method %q. Valid methods are either %q, %q, %q, %q, %q, %q, %q, or %q.",
			k, value, "ANY", "DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"))
	}
	return
}

func validateLogMetricFilterName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 512 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 512 characters: %q", k, value))
	}

	// http://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_PutMetricFilter.html
	pattern := `^[^:*]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q isn't a valid log metric name (must not contain colon nor asterisk): %q",
			k, value))
	}

	return
}

func validateLogMetricFilterTransformationName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 255 characters: %q", k, value))
	}

	// http://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_MetricTransformation.html
	pattern := `^[^:*$]*$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q isn't a valid log metric transformation name (must not contain"+
				" colon, asterisk nor dollar sign): %q",
			k, value))
	}

	return
}

func validateLogGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 512 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 512 characters: %q", k, value))
	}

	// http://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_CreateLogGroup.html
	pattern := `^[\.\-_/#A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q isn't a valid log group name (alphanumeric characters, underscores,"+
				" hyphens, slashes, hash signs and dots are allowed): %q",
			k, value))
	}

	return
}

func validateS3BucketLifecycleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as RFC3339 Timestamp Format", value))
	}

	return
}

func validateS3BucketLifecycleStorageClass(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != s3.TransitionStorageClassStandardIa && value != s3.TransitionStorageClassGlacier {
		errors = append(errors, fmt.Errorf(
			"%q must be one of '%q', '%q'", k, s3.TransitionStorageClassStandardIa, s3.TransitionStorageClassGlacier))
	}

	return
}

func validateS3BucketReplicationRuleId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 255 characters: %q", k, value))
	}

	return
}

func validateS3BucketReplicationRulePrefix(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 1024 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 1024 characters: %q", k, value))
	}

	return
}

func validateS3BucketReplicationDestinationStorageClass(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != s3.StorageClassStandard && value != s3.StorageClassStandardIa && value != s3.StorageClassReducedRedundancy {
		errors = append(errors, fmt.Errorf(
			"%q must be one of '%q', '%q' or '%q'", k, s3.StorageClassStandard, s3.StorageClassStandardIa, s3.StorageClassReducedRedundancy))
	}

	return
}

func validateS3BucketReplicationRuleStatus(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != s3.ReplicationRuleStatusEnabled && value != s3.ReplicationRuleStatusDisabled {
		errors = append(errors, fmt.Errorf(
			"%q must be one of '%q' or '%q'", k, s3.ReplicationRuleStatusEnabled, s3.ReplicationRuleStatusDisabled))
	}

	return
}

func validateS3BucketLifecycleRuleId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q cannot exceed 255 characters", k))
	}
	return
}

func validateDbEventSubscriptionName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[0-9A-Za-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alphanumeric characters and hyphens allowed in %q", k))
	}
	if len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 255 characters", k))
	}
	return
}

func validateApiGatewayIntegrationPassthroughBehavior(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "WHEN_NO_MATCH" && value != "WHEN_NO_TEMPLATES" && value != "NEVER" {
		errors = append(errors, fmt.Errorf(
			"%q must be one of 'WHEN_NO_MATCH', 'WHEN_NO_TEMPLATES', 'NEVER'", k))
	}
	return
}

func validateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

func validateApiGatewayIntegrationType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	validTypes := map[string]bool{
		"AWS":        true,
		"AWS_PROXY":  true,
		"HTTP":       true,
		"HTTP_PROXY": true,
		"MOCK":       true,
	}

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid integration type %q. Valid types are either %q, %q, %q, %q, or %q.",
			k, value, "AWS", "AWS_PROXY", "HTTP", "HTTP_PROXY", "MOCK"))
	}
	return
}

func validateApiGatewayIntegrationContentHandling(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	validTypes := map[string]bool{
		"CONVERT_TO_BINARY": true,
		"CONVERT_TO_TEXT":   true,
	}

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid integration type %q. Valid types are either %q or %q.",
			k, value, "CONVERT_TO_BINARY", "CONVERT_TO_TEXT"))
	}
	return
}

func validateSQSQueueName(v interface{}, k string) (errors []error) {
	value := v.(string)
	if len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 80 characters", k))
	}

	if !regexp.MustCompile(`^[0-9A-Za-z-_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("only alphanumeric characters and hyphens allowed in %q", k))
	}
	return
}

func validateSQSFifoQueueName(v interface{}, k string) (errors []error) {
	value := v.(string)

	if len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 80 characters", k))
	}

	if !regexp.MustCompile(`^[0-9A-Za-z-_.]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("only alphanumeric characters and hyphens allowed in %q", k))
	}

	if regexp.MustCompile(`^[^a-zA-Z0-9-_]`).MatchString(value) {
		errors = append(errors, fmt.Errorf("FIFO queue name must start with one of these characters [a-zA-Z0-9-_]: %v", value))
	}

	if !regexp.MustCompile(`\.fifo$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("FIFO queue name should ends with \".fifo\": %v", value))
	}

	return
}

func validateSNSSubscriptionProtocol(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	forbidden := []string{"email", "sms"}
	for _, f := range forbidden {
		if strings.Contains(value, f) {
			errors = append(
				errors,
				fmt.Errorf("Unsupported protocol (%s) for SNS Topic", value),
			)
		}
	}
	return
}

func validateSecurityRuleType(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))

	validTypes := map[string]bool{
		"ingress": true,
		"egress":  true,
	}

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Security Group Rule type %q. Valid types are either %q or %q.",
			k, value, "ingress", "egress"))
	}
	return
}

func validateOnceAWeekWindowFormat(v interface{}, k string) (ws []string, errors []error) {
	// valid time format is "ddd:hh24:mi"
	validTimeFormat := "(sun|mon|tue|wed|thu|fri|sat):([0-1][0-9]|2[0-3]):([0-5][0-9])"

	value := strings.ToLower(v.(string))
	if !regexp.MustCompile(validTimeFormat + "-" + validTimeFormat).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must satisfy the format of \"ddd:hh24:mi-ddd:hh24:mi\".", k))
	}
	return
}

func validateOnceADayWindowFormat(v interface{}, k string) (ws []string, errors []error) {
	// valid time format is "hh24:mi"
	validTimeFormat := "([0-1][0-9]|2[0-3]):([0-5][0-9])"

	value := v.(string)
	if !regexp.MustCompile(validTimeFormat + "-" + validTimeFormat).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must satisfy the format of \"hh24:mi-hh24:mi\".", k))
	}
	return
}

func validateRoute53RecordType(v interface{}, k string) (ws []string, errors []error) {
	// Valid Record types
	// SOA, A, TXT, NS, CNAME, MX, NAPTR, PTR, SRV, SPF, AAAA
	validTypes := map[string]struct{}{
		"SOA":   {},
		"A":     {},
		"TXT":   {},
		"NS":    {},
		"CNAME": {},
		"MX":    {},
		"NAPTR": {},
		"PTR":   {},
		"SRV":   {},
		"SPF":   {},
		"AAAA":  {},
	}

	value := v.(string)
	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q must be one of [SOA, A, TXT, NS, CNAME, MX, NAPTR, PTR, SRV, SPF, AAAA]", k))
	}
	return
}

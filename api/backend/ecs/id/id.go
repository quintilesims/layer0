package id

import (
	"crypto/md5"
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"regexp"
	"strings"
	"time"
)

const (
	// maximum elb name length (the most limiting id) = 32
	// maximum prefix length = 15
	// maximum prefix length = l0-<prefix>- = 19
	// maximum id length = 31 - 19 = 13 (add buffer of 1 just to be safe) = 12
	MAX_ID_LENGTH      = 12
	MIN_ID_LENGTH      = 2
	MIN_ID_HASH_LENGTH = 5
)

var PREFIX = fmt.Sprintf("l0-%s-", config.Prefix())

func isQuote(r rune) bool {
	return r == '"'
}

func addPrefix(id string) string {
	return fmt.Sprintf("%s%s", PREFIX, id)
}

func removePrefix(id string) string {
	return strings.TrimPrefix(id, PREFIX)
}

// used for testing
func StubIDGeneration(id string) func() {
	tmpHashed := GenerateHashedEntityID
	tmpHashless := GenerateHashlessEntityID

	GenerateHashedEntityID = func(string) string { return id }
	GenerateHashlessEntityID = func(string) string { return id }

	return func() {
		GenerateHashedEntityID = tmpHashed
		GenerateHashlessEntityID = tmpHashless
	}
}

// generates an id using only the name while still being
// being safe in regards to length and character limitations (if possible)
var GenerateHashlessEntityID = func(name string) string {
	id := filterUsableName(name)

	if len(id) > MAX_ID_LENGTH {
		id = id[:MAX_ID_LENGTH]
	}

	if len(id) < MIN_ID_LENGTH {
		hashLength := MIN_ID_LENGTH - len(id)
		hash := hashNow()[:hashLength]
		id = id + hash
	}

	return id
}

// generates an id using as much of the name as possible
// with at least MIN_ID_HASH_LENGTH characters of the id randomly hashed
// the id will always be MAX_ID_LENGTH characters in length
var GenerateHashedEntityID = func(name string) string {
	prefix := filterUsableName(name)

	if maxPrefixLength := MAX_ID_LENGTH - MIN_ID_HASH_LENGTH; len(prefix) > maxPrefixLength {
		prefix = prefix[:maxPrefixLength]
	}

	hashLength := MAX_ID_LENGTH - len(prefix)
	hash := hashNow()[:hashLength]

	return prefix + hash
}

func filterUsableName(name string) string {
	// only allow alphanumerics in entity ids
	reg := regexp.MustCompile("[^A-Za-z0-9]+")
	return reg.ReplaceAllString(name, "")
}

func hashNow() string {
	salt := time.Now().Format(time.StampNano)
	return fmt.Sprintf("%x", md5.Sum([]byte(salt)))
}

type ECSLoadBalancerID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id ECSLoadBalancerID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id ECSLoadBalancerID) L0LoadBalancerID() string {
	return removePrefix(id.String())
}

func (id ECSLoadBalancerID) SecurityGroupName() string {
	return fmt.Sprintf("%s-lb", id.String())
}

func (id ECSLoadBalancerID) RoleName() string {
	return fmt.Sprintf("%s-lb", id.String())
}

type L0LoadBalancerID string

func (id L0LoadBalancerID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id L0LoadBalancerID) ECSLoadBalancerID() ECSLoadBalancerID {
	str := addPrefix(id.String())
	return ECSLoadBalancerID(str)
}

type ECSEnvironmentID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id ECSEnvironmentID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id ECSEnvironmentID) L0EnvironmentID() string {
	return removePrefix(id.String())
}

func (id ECSEnvironmentID) SecurityGroupName() string {
	return fmt.Sprintf("%s-env", id.String())
}

func (id ECSEnvironmentID) LaunchConfigurationName() string {
	return id.String()
}

func (id ECSEnvironmentID) AutoScalingGroupName() string {
	return id.String()
}

func ClusterARNToECSEnvironmentID(arn string) ECSEnvironmentID {
	clusterName := strings.SplitN(arn, "/", 2)[1]
	return ECSEnvironmentID(clusterName)
}

type L0EnvironmentID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id L0EnvironmentID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id L0EnvironmentID) ECSEnvironmentID() ECSEnvironmentID {
	str := addPrefix(id.String())
	return ECSEnvironmentID(str)
}

type ECSServiceID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id ECSServiceID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

// service log group format is 'l0-<prefix>-<environment_id>-<service_id>-svc'
func (id ECSServiceID) LogGroupID(environmentID ECSEnvironmentID) string {
	groupName := fmt.Sprintf("%s-%s-svc", environmentID.L0EnvironmentID(), id.L0ServiceID())
	return addPrefix(groupName)
}

func (id ECSServiceID) L0ServiceID() string {
	return removePrefix(id.String())
}

type L0ServiceID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id L0ServiceID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id L0ServiceID) ECSServiceID() ECSServiceID {
	str := addPrefix(id.String())
	return ECSServiceID(str)
}

type ECSDeployID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id ECSDeployID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id ECSDeployID) L0DeployID() string {
	return removePrefix(id.String())
}

func (id ECSDeployID) TaskDefinition() string {
	return strings.Replace(id.String(), ".", ":", -1)
}

func (id ECSDeployID) FamilyName() string {
	return strings.Split(id.L0DeployID(), ".")[0]
}

func (id ECSDeployID) Revision() string {
	if split := strings.Split(id.L0DeployID(), "."); len(split) > 1 {
		return split[1]
	}

	return "1"
}

func TaskDefinitionToECSDeployID(taskDefinition string) ECSDeployID {
	taskDefinitionID := strings.Replace(taskDefinition, ":", ".", -1)
	return ECSDeployID(taskDefinitionID)
}

func TaskDefinitionARNToECSDeployID(arn string) ECSDeployID {
	taskDefinitionID := strings.SplitN(arn, "/", 2)[1]
	taskDefinitionID = strings.Replace(taskDefinitionID, ":", ".", -1)
	return ECSDeployID(taskDefinitionID)
}

// task log group format is 'l0-<prefix>-<environment_id>-<task_id>-tsk'
func (id ECSTaskID) LogGroupID(environmentID ECSEnvironmentID) string {
	groupName := fmt.Sprintf("%s-%s-tsk", environmentID.L0EnvironmentID(), id.L0TaskID())
	return addPrefix(groupName)
}

type L0DeployID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id L0DeployID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id L0DeployID) ECSDeployID() ECSDeployID {
	str := addPrefix(id.String())
	return ECSDeployID(str)
}

type ECSTaskID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id ECSTaskID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id ECSTaskID) L0TaskID() string {
	return removePrefix(id.String())
}

type L0TaskID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id L0TaskID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id L0TaskID) ECSTaskID() ECSTaskID {
	str := addPrefix(id.String())
	return ECSTaskID(str)
}

type ECSCertificateID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id ECSCertificateID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id ECSCertificateID) L0CertificateID() string {
	return removePrefix(id.String())
}

func CertificateARNToECSCertificateID(arn string) ECSCertificateID {
	split := strings.SplitN(arn, "/", -1)
	certificateName := split[len(split)-1]
	return ECSCertificateID(certificateName)
}

type L0CertificateID string

// we need to add a custom .String() function, or else string conversions add quotes
func (id L0CertificateID) String() string {
	return strings.TrimFunc(string(id), isQuote)
}

func (id L0CertificateID) ECSCertificateID() ECSCertificateID {
	str := addPrefix(id.String())
	return ECSCertificateID(str)
}

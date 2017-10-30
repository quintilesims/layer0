package aws

import (
	"crypto/md5"
	"fmt"
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

var entityIDGenerator func(string) string

func init(){
	entityIDGenerator = defaultEntityIDGenerator
}

func prefix(instance string) string {
	return fmt.Sprintf("l0-%s-", instance)
}

func hasLayer0Prefix(instance, v string) bool {
	return strings.HasPrefix(v, prefix(instance))
}

func addLayer0Prefix(instance, v string) string {
	if !hasLayer0Prefix(instance, v) {
		v = fmt.Sprintf("%s%s", prefix(instance), v)
	}

	return v
}

func delLayer0Prefix(instance, v string) string {
	return strings.TrimPrefix(v, prefix(instance))
}

// generates an id using as much of the name as possible
// with at least MIN_ID_HASH_LENGTH characters of the id randomly hashed
// the id will always be MAX_ID_LENGTH characters in length
func defaultEntityIDGenerator(name string) string {
	prefix := filterUsableID(name)

	if maxPrefixLength := MAX_ID_LENGTH - MIN_ID_HASH_LENGTH; len(prefix) > maxPrefixLength {
		prefix = prefix[:maxPrefixLength]
	}

	hashLength := MAX_ID_LENGTH - len(prefix)
	hash := hashNow()[:hashLength]

	return prefix + hash
}

func SetEntityIDGenerator(entityID string) func()  {
	entityIDGenerator = func(string) string{
		return entityID
	}

	return func(){ entityIDGenerator = defaultEntityIDGenerator }
}

// filters out any non-alphanumeric characters
func filterUsableID(name string) string {
	reg := regexp.MustCompile("[^A-Za-z0-9]+")
	return reg.ReplaceAllString(name, "")
}

func hashNow() string {
	salt := time.Now().Format(time.StampNano)
	return fmt.Sprintf("%x", md5.Sum([]byte(salt)))
}

package lock

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type DynamoDBExpiringLock struct {
	table  dynamo.Table
	lockID string
	expiry time.Duration
}

func NewDynamoDBExpiringLock(config *aws.Config, table, lockID string, expiry time.Duration) *DynamoDBExpiringLock {
	session := session.New(config)
	db := dynamo.New(session)

	return &DynamoDBExpiringLock{
		table:  db.Table(table),
		lockID: lockID,
		expiry: expiry,
	}
}

func (d *DynamoDBExpiringLock) Acquire() error {
	// if entry at d.lockID does not exist, create it and return nil
	// if entry exists but entry.Timestamp older than d.expiry, update it return nil
	// if entry exists and entry.Timestamp is not older than d.expiry, return AcquiredError
	return nil
}

func (d *DynamoDBExpiringLock) Release() error {
	// delete entry from table at lockID
	return nil
}

package lock

import (
	"log"
	"time"
)

func NewDaemonFN(dynamoLock *DynamoLock, expiry time.Duration) func() error {
	return func() error {
		locks, err := dynamoLock.List()
		if err != nil {
			return err
		}

		oldestPossibleAcquiredTime := time.Now().Add(-expiry)
		for _, l := range locks {
			if l.Acquired < oldestPossibleAcquiredTime.UnixNano() {
				log.Printf("[DEBUG] [LockDaemon] Deleting expired lock %s", l.LockID)
				if err := dynamoLock.Release(l.LockID); err != nil {
					log.Printf("[ERROR] [LockDaemon] Failed to delete expired lock %s: %v", l.LockID, err)
				}
			}
		}

		return nil
	}
}

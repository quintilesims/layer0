package lock

type Lock interface {
	Acquire(lockID string) (bool, error)
	Release(lockID string) error
}

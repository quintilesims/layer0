package lock

type Lock interface {
	Acquire() (bool, error)
	Release() error
}

package instance

import (
	"os"
)

func (l *LocalInstance) Destroy(force bool) error {
	if err := l.assertExists(); err != nil {
		return err
	}

	// todo: use layer0 client to destroy all resources

	if err := l.Terraform.Destroy(l.Dir, force); err != nil {
		return err
	}

	return os.RemoveAll(l.Dir)
}

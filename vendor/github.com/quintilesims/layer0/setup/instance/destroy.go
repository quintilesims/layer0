package instance

import (
	"os"
)

func (l *LocalInstance) Destroy(force bool) error {
	if err := l.assertExists(); err != nil {
		return err
	}

	if err := l.Terraform.Destroy(l.Dir, force); err != nil {
		return err
	}

	return os.RemoveAll(l.Dir)
}

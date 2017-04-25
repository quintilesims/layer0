package instance

import (
	"os"
)

func (i *Instance) Destroy(force bool) error {
	if err := i.assertExists(); err != nil {
		return err
	}

	// todo: use layer0 client to destroy all resources

	if err := i.Terraform.Destroy(i.Dir, force); err != nil {
		return err
	}

	return os.RemoveAll(i.Dir)
}

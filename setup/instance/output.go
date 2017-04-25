package instance

func (i *Instance) Output(key string) (string, error) {
	if err := i.assertExists(); err != nil {
		return "", err
	}

	return i.Terraform.Output(i.Dir, key)
}

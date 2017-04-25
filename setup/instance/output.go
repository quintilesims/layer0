package instance

func (i *LocalInstance) Output(key string) (string, error) {
	if err := i.assertExists(); err != nil {
		return "", err
	}

	return i.Terraform.Output(i.Dir, key)
}

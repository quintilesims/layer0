package instance

func (l *LocalInstance) Output(key string) (string, error) {
	if err := l.assertExists(); err != nil {
		return "", err
	}

	return l.Terraform.Output(l.Dir, key)
}

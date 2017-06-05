package instance

func (l *LocalInstance) Plan() error {
	if err := l.assertExists(); err != nil {
		return err
	}

	return l.Terraform.Plan(l.Dir)
}

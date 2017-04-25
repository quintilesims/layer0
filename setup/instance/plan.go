package instance

func (i *Instance) Plan() error {
	if err := i.assertExists(); err != nil {
		return err
	}

	return i.Terraform.Plan(i.Dir)
}

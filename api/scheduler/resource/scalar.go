package resource

type EnvironmentScalar interface {
	ScaleTo(int) error
}

type EnvironmentScalarFunc func(int) error

func (e EnvironmentScalarFunc) ScaleTo(scale int) error {
	return e(scale)
}

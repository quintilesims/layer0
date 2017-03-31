package instance

type Factory interface {
	NewInstance(name string) (Instance, error)
}

type Layer0Factory struct{}

func NewLayer0Factory() *Layer0Factory {
	return &Layer0Factory{}
}

func (f *Layer0Factory) NewInstance(name string) (Instance, error) {
	return NewLayer0Instance(name), nil
}

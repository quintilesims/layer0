package scaler

type Scaler interface {
	Scale(environmentID string) error
}

type ScalerFunc func(string) error

func (s ScalerFunc) Scale(environmentID string) error {
	return s(environmentID)
}

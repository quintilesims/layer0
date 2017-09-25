package scaler

type Scaler interface {
	Scale(environmentID string) error
}

type ScalerFunc func(environmentID string) error

func (r ScalerFunc) Scale(environmentID string) error {
	return r(environmentID)
}

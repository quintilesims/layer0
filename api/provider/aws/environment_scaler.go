package aws

import "fmt"

type EnvironmentScaler struct {
}

func NewEnvironmentScaler() *EnvironmentScaler {
	return &EnvironmentScaler{}
}

func (e *EnvironmentScaler) Scale(environmentID string) error {
	return fmt.Errorf("EnvironmentScaler not implemented")
}

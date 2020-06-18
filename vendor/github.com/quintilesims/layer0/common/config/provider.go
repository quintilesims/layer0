package config

type ConfigCredProvider struct{}

func NewConfigCredProvider() *ConfigCredProvider {
	return &ConfigCredProvider{}
}

func (this *ConfigCredProvider) GetAWSAccessKeyID() (string, error) {
	return AWSAccessKey(), nil
}

func (this *ConfigCredProvider) GetAWSSecretAccessKey() (string, error) {
	return AWSSecretKey(), nil
}

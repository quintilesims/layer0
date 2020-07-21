package provider

type CredProvider interface {
	GetAWSAccessKeyID() (string, error)
	GetAWSSecretAccessKey() (string, error)
}

type ExplicitCredProvider struct {
	AccessKey, SecretKey string
}

func NewExplicitCredProvider(access, secret string) *ExplicitCredProvider {
	return &ExplicitCredProvider{access, secret}
}

func (this *ExplicitCredProvider) GetAWSAccessKeyID() (string, error) {
	return this.AccessKey, nil
}

func (this *ExplicitCredProvider) GetAWSSecretAccessKey() (string, error) {
	return this.SecretKey, nil
}

type MockCredProvider struct {
	GetAWSAccessKeyID_fn     func() (string, error)
	GetAWSSecretAccessKey_fn func() (string, error)
}

func (this *MockCredProvider) GetAWSAccessKeyID() (string, error) {
	if this.GetAWSAccessKeyID_fn == nil {
		this.GetAWSAccessKeyID_fn = func() (string, error) { return "", nil }
	}
	return this.GetAWSAccessKeyID_fn()
}

func (this *MockCredProvider) GetAWSSecretAccessKey() (string, error) {
	if this.GetAWSSecretAccessKey_fn == nil {
		this.GetAWSSecretAccessKey_fn = func() (string, error) { return "", nil }
	}
	return this.GetAWSSecretAccessKey_fn()
}

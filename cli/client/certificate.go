package client

import (
	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) CreateCertificate(name string, public, private, chain []byte) (*models.Certificate, error) {
	req := models.CreateCertificateRequest{
		CertificateName:  name,
		PublicCert:       string(public),
		PrivateKey:       string(private),
		IntermediateCert: string(chain),
	}

	var certificate *models.Certificate
	if err := c.Execute(c.Sling("certificate/").Post("").BodyJSON(req), &certificate); err != nil {
		return nil, err
	}

	return certificate, nil
}

func (c *APIClient) DeleteCertificate(id string) error {
	var response *string
	if err := c.Execute(c.Sling("certificate/").Delete(id), &response); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) GetCertificate(id string) (*models.Certificate, error) {
	var certificate *models.Certificate
	if err := c.Execute(c.Sling("certificate/").Get(id), &certificate); err != nil {
		return nil, err
	}

	return certificate, nil
}

func (c *APIClient) ListCertificates() ([]*models.Certificate, error) {
	var certificates []*models.Certificate
	if err := c.Execute(c.Sling("certificate/").Get(""), &certificates); err != nil {
		return nil, err
	}

	return certificates, nil
}

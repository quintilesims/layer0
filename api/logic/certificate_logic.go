package logic

import (
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

type CertificateLogic interface {
	CreateCertificate(models.CreateCertificateRequest) (*models.Certificate, error)
	ListCertificates() ([]*models.Certificate, error)
	GetCertificate(string) (*models.Certificate, error)
	DeleteCertificate(string) error
}

type L0CertificateLogic struct {
	Logic
}

func NewL0CertificateLogic(lgc Logic) *L0CertificateLogic {
	return &L0CertificateLogic{lgc}
}

func (this *L0CertificateLogic) ListCertificates() ([]*models.Certificate, error) {
	certificates, err := this.Backend.ListCertificates()
	if err != nil {
		return nil, errors.New(errors.UnexpectedError, err)
	}

	for _, certificate := range certificates {
		if err := this.populateModel(certificate); err != nil {
			return nil, err
		}
	}

	return certificates, nil
}

func (this *L0CertificateLogic) GetCertificate(certificateID string) (*models.Certificate, error) {
	certificate, err := this.Backend.GetCertificate(certificateID)
	if err != nil {
		return nil, err
	}

	if err := this.populateModel(certificate); err != nil {
		return nil, err
	}

	return certificate, nil
}

func (this *L0CertificateLogic) DeleteCertificate(certificateID string) error {
	if err := this.Backend.DeleteCertificate(certificateID); err != nil {
		return err
	}

	if err := this.deleteEntityTags(certificateID, "certificate"); err != nil {
		return err
	}

	return nil
}

func (this *L0CertificateLogic) CreateCertificate(req models.CreateCertificateRequest) (*models.Certificate, error) {
	if req.CertificateName == "" {
		return nil, errors.Newf(errors.MissingParameter, "EnvironmentName is required")
	}

	certificate, err := this.Backend.CreateCertificate(req.CertificateName, req.PublicCert, req.PrivateKey, req.IntermediateCert)
	if err != nil {
		return nil, err
	}

	if err := this.upsertTagf(certificate.CertificateID, "certificate", "name", req.CertificateName); err != nil {
		return nil, errors.New(errors.FailedTagging, err)
	}

	if err := this.populateModel(certificate); err != nil {
		return nil, err
	}

	return certificate, nil
}

func (this *L0CertificateLogic) populateModel(model *models.Certificate) error {
	filter := map[string]string{
		"type": "certificate",
		"id":   model.CertificateID,
	}

	tags, err := this.TagData.GetTags(filter)
	if err != nil {
		return err
	}

	for _, tag := range rangeTags(tags) {
		if tag.Key == "name" {
			model.CertificateName = tag.Value
			break
		}
	}

	return nil

}

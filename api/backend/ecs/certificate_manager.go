package ecsbackend

import (
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs/id"
	"gitlab.imshealth.com/xfra/layer0/common/aws/iam"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"strings"
)

type ECSCertificateManager struct {
	IAM iam.Provider
}

func NewECSCertificateManager(iam iam.Provider) *ECSCertificateManager {
	return &ECSCertificateManager{
		IAM: iam,
	}
}

func (this *ECSCertificateManager) GetCertificate(certificateID string) (*models.Certificate, error) {
	certificates, err := this.ListCertificates()
	if err != nil {
		return nil, err
	}

	for _, certificate := range certificates {
		if certificate.CertificateID == certificateID {
			return certificate, nil
		}
	}

	err = fmt.Errorf("Certificate with id '%s' does not exist", certificateID)
	return nil, errors.New(errors.InvalidCertificateID, err)
}

func (this *ECSCertificateManager) ListCertificates() ([]*models.Certificate, error) {
	certificates, err := this.IAM.ListCertificates(CertificatePath())
	if err != nil {
		return nil, err
	}

	models := []*models.Certificate{}
	for _, certificate := range certificates {
		if strings.HasPrefix(*certificate.ServerCertificateName, id.PREFIX) {
			model := this.populateModel(certificate)
			models = append(models, model)
		}
	}

	return models, nil
}

func (this *ECSCertificateManager) DeleteCertificate(certificateID string) error {
	ecsCertificateID := id.L0CertificateID(certificateID).ECSCertificateID()

	if err := this.IAM.DeleteServerCertificate(ecsCertificateID.String()); err != nil {
		return err
	}

	return nil
}

func (this *ECSCertificateManager) CreateCertificate(certificateName, public, private, chain string) (*models.Certificate, error) {
	// we don't generate a hashed id for certificates since aws enforces unique certificate names
	certificateID := id.GenerateHashlessEntityID(certificateName)
	ecsCertificateID := id.L0CertificateID(certificateID).ECSCertificateID()

	var optionalChain *string
	if chain != "" {
		optionalChain = &chain
	}

	certificate, err := this.IAM.UploadServerCertificate(
		ecsCertificateID.String(),
		CertificatePath(),
		public,
		private,
		optionalChain)
	if err != nil {
		return nil, err
	}

	return this.populateModel(certificate), nil
}

func (this *ECSCertificateManager) populateModel(cert *iam.ServerCertificateMetadata) *models.Certificate {
	return &models.Certificate{
		CertificateID:  id.ECSCertificateID(*cert.ServerCertificateName).L0CertificateID(),
		CertificateARN: *cert.Arn,
	}
}

func CertificatePath() string {
	return fmt.Sprintf("/l0/l0-%v/", config.Prefix())
}

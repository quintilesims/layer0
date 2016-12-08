package entity

import (
	"github.com/quintilesims/layer0/cli/printer/table"
	"github.com/quintilesims/layer0/common/models"
)

type Certificate models.Certificate

func NewCertificate(model *models.Certificate) *Certificate {
	certificate := Certificate(*model)
	return &certificate
}

func (this *Certificate) Table() table.Table {
	table := []table.Column{
		table.NewSingleRowColumn("CERTIFICATE ID", this.CertificateID),
		table.NewSingleRowColumn("CERTIFICATE NAME", this.CertificateName),
	}

	return table
}

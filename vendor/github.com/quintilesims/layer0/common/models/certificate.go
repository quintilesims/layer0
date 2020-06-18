package models

type Certificate struct {
	CertificateARN  string `json:"certificate_arn"`
	CertificateID   string `json:"certificate_id"`
	CertificateName string `json:"certificate_name"`
}

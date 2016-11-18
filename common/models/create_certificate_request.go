package models

type CreateCertificateRequest struct {
	CertificateName  string `json:"certificate_name"`
	IntermediateCert string `json:"intermediate_cert"`
	PrivateKey       string `json:"private_key"`
	PublicCert       string `json:"public_cert"`
}

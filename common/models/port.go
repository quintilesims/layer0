package models

type Port struct {
	CertificateID string `json:"certificate_id"`
	ContainerPort int64  `json:"container_port"`
	HostPort      int64  `json:"host_port"`
	Protocol      string `json:"protocol"`
}

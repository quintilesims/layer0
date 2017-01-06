package models

type Port struct {
	CertificateName string `json:"certificate_name"`
	ContainerPort   int64  `json:"container_port"`
	HostPort        int64  `json:"host_port"`
	Protocol        string `json:"protocol"`
}

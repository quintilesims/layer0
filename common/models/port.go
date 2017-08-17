package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type Port struct {
	CertificateName string `json:"certificate_name"`
	ContainerPort   int64  `json:"container_port"`
	HostPort        int64  `json:"host_port"`
	Protocol        string `json:"protocol"`
}

func (p Port) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"certificate_name": swagger.NewStringProperty(),
			"container_port":   swagger.NewIntProperty(),
			"host_port":        swagger.NewIntProperty(),
			"protocol":         swagger.NewStringProperty(),
		},
	}
}

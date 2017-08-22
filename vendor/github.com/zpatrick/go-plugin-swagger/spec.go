package swagger

type Spec struct {
	SwaggerVersion      string                        `json:"swagger"`
	Host                string                        `json:"host,omitempty"`
	BasePath            string                        `json:"basePath,omitempty"`
	Info                *Info                         `json:"info,omitempty"`
	Schemes             []string                      `json:"schemes"`
	ExternalDocs        *ExternalDocs                 `json:"externalDocs,omitempty"`
	Tags                []Tag                         `json:"tags"`
	Paths               map[string]Path               `json:"paths"`
	Definitions         map[string]Definition         `json:"definitions"`
	SecurityDefinitions map[string]SecurityDefinition `json:"securityDefinitions"`
}

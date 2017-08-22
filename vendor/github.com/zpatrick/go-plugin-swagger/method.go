package swagger

type Method struct {
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	Responses   map[string]Response   `json:"responses,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
}

func BasicAuthSecurity(key string) []map[string][]string {
	return []map[string][]string{
		{
			key: []string{},
		},
	}
}

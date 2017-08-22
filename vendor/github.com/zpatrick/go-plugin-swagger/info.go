package swagger

type Info struct {
	Title          string   `json:"title,omitempty"`
	Version        string   `json:"version,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

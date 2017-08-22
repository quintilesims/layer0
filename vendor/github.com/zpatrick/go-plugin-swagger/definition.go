package swagger

type Definition struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

package swagger

type Response struct {
	Description string  `json:"description"`
	Schema      *Schema `json:"schema,omitempty"`
}

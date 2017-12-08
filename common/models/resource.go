package models

type Resource struct {
	ID     string  `json:"id"`
	InUse  bool    `json:"in_use"`
	Ports  []int   `json:"ports"`
	Memory float64 `json:"memory"`
	CPU    float64 `json:"cpu"`
}

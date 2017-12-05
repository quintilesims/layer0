package models

type ResourceConsumer struct {
	CPU    string `json:"cpu"`
	ID     string `json:"id"`
	Memory string `json:"memory"`
	Ports  []int  `json:"ports"`
}

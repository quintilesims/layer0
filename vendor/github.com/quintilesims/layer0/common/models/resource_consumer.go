package models

type ResourceConsumer struct {
	ID     string `json:"id"`
	Memory string `json:"memory"`
	Ports  []int  `json:"ports"`
}

package models

type ResourceProvider struct {
	ID              string `json:"id"`
	InUse           bool   `json:"in_use"`
	UsedPorts       []int  `json:"used_ports"`
	AvailableMemory string `json:"available_memory"`
}

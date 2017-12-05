package models

type ResourceProvider struct {
	AgentConnected  bool   `json:"agent_connected"`
	AvailableCPU    string `json:"available_cpu"`
	AvailableMemory string `json:"available_memory"`
	ID              string `json:"id"`
	InUse           bool   `json:"in_use"`
	Status          string `json:"status"`
	UsedPorts       []int  `json:"used_ports"`
}

package models

type Environment struct {
	EnvironmentID   string   `json:"environment_id"`
	EnvironmentName string   `json:"environment_name"`
	ClusterCount    int      `json:"cluster_count"`
	MinCount        int      `json:"min_count"`
	MaxCount        int      `json:"max_count"`
	TargetCapSize   int      `json:"target_cap_size"`
	InstanceSize    string   `json:"instance_size"`
	SecurityGroupID string   `json:"security_group_id"`
	OperatingSystem string   `json:"operating_system"`
	AMIID           string   `json:"ami_id"`
	Links           []string `json:"links"`
}

package models

type CreateEnvironmentRequest struct {
	EnvironmentName  string `json:"environment_name"`
	InstanceSize     string `json:"instance_size"`
	UserDataTemplate []byte `json:"user_data_template"`
	MinClusterCount  int    `json:"min_cluster_count"`
	MaxClusterCount  int    `json:"max_cluster_count"`
	TargetCapSize    int    `json:"target_cap_size"`
	OperatingSystem  string `json:"operating_system"`
	AMIID            string `json:"ami_id"`
}

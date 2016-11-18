package models

type CreateEnvironmentRequest struct {
	EnvironmentName  string `json:"environment_name"`
	InstanceSize     string `json:"instance_size"`
	UserDataTemplate []byte `json:"user_data_template"`
	MinClusterCount  int    `json:"min_cluster_count"`
}

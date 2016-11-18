package models

type LoadBalancer struct {
	EnvironmentID    string   `json:"environment_id"`
	EnvironmentName  string   `json:"environment_name"`
	IsPublic         bool     `json:"is_public"`
	LoadBalancerID   string   `json:"load_balancer_id"`
	LoadBalancerName string   `json:"load_balancer_name"`
	Ports            []Port   `json:"ports"`
	Services         []string `json:"services"`
	URL              string   `json:"url"`
}

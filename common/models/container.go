package models

type Container struct {
	ID         string            `json:"id"`
	Command    string            `json:"command"`
	Created    int64             `json:"created"`
	Image      string            `json:"image"`
	ImageID    string            `json:"image_id"`
	Labels     map[string]string `json:"labels"`
	Names      []string          `json:"names"`
	Ports      []Port            `json:"ports"`
	SizeRootFS int64             `json:"size_root_fs"`
	SizeRW     int64             `json:"size_rw"`
	State      string            `json:"state"`
	Status     string            `json:"status"`
}

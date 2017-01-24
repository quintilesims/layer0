package models

type Health struct {
	TimeCreated string
	Mode        string
}

type SetHealthRequest struct {
	Mode string
}

type Command struct {
	Name   string
	Args   []string
	Output string
}

type CreateCommandRequest struct {
	Name string
	Args []string
}

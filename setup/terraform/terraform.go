package terraform

import (
	"os"
	"os/exec"
)

type Terraform struct{}

func New() *Terraform {
	return &Terraform{}
}

func (t *Terraform) Apply(dir string) error {
	return t.run(dir, "apply")
}

func (t *Terraform) Get(dir string) error {
	return t.run(dir, "get")
}

func (t *Terraform) FMT(dir string) error {
        return t.run(dir, "fmt")
}

func (t *Terraform) run(dir string, args ...string) error {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

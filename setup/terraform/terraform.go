package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type Terraform struct{}

func New() *Terraform {
	return &Terraform{}
}

func (t *Terraform) Apply(dir string) error {
	return t.run(dir, "apply")
}

func (t *Terraform) Destroy(dir string, force bool) error {
	if force {
		return t.run(dir, "destroy", "-force")
	}

	return t.run(dir, "destroy")
}

func (t *Terraform) FMT(dir string) error {
	return t.run(dir, "fmt")
}

func (t *Terraform) Get(dir string) error {
	return t.run(dir, "get", "-update")
}

func (t *Terraform) Output(dir, key string) (string, error) {
	cmd := exec.Command("terraform", "output", "-module", "layer0", key)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		text := fmt.Sprintf("Error running %v from %s: %v\n", cmd.Args, cmd.Dir, err)
		for _, line := range strings.Split(string(output), "\n") {
			text += line + "\n"
		}

		return "", fmt.Errorf(text)
	}

	return strings.Replace(string(output), "\n", "", 1), nil
}

func (t *Terraform) Plan(dir string) error {
	return t.run(dir, "plan")
}

func (t *Terraform) run(dir string, args ...string) error {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go t.handleSIGTERM(cmd)
	return cmd.Run()
}

func (t *Terraform) handleSIGTERM(cmd *exec.Cmd) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	cmd.Process.Kill()
	os.Exit(1)
}

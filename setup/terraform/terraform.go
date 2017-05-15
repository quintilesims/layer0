package terraform

import (
	"fmt"
	"github.com/blang/semver"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

const REQUIRED_TERRAFORM_VERSION = "0.9.4"

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

func (t *Terraform) Validate(dir string) error {
	return t.run(dir, "validate")
}

func (t *Terraform) Get(dir string) error {
	return t.run(dir, "get", "-update")
}

func (t *Terraform) Output(dir, key string) (string, error) {
	if err := t.validateTerraformVersion(); err != nil {
		return "", err
	}

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
	if err := t.validateTerraformVersion(); err != nil {
		return err
	}

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

func (t *Terraform) validateTerraformVersion() error {
	required, err := semver.Make(REQUIRED_TERRAFORM_VERSION)
	if err != nil {
		return fmt.Errorf("Failed to parse required Terraform version: %v", err)
	}

	cmd := exec.Command("terraform", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Could not determine current Terraform version: %v", err)
	}

	// only grab the first line - terraform will add additional messages
	// when terraform is out o fdate
	version := strings.Split(string(output), "\n")[0]
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "Terraform v")

	current, err := semver.Make(version)
	if err != nil {
		return fmt.Errorf("Failed to parse current Terraform version: %v", err)
	}

	if current.LT(required) {
		text := fmt.Sprintf("Current version of Terraform (%s) is less than the ", current)
		text += fmt.Sprintf("minimum required version (%s)", required)
		return fmt.Errorf(text)
	}

	return nil
}

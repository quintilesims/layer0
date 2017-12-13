package tftest

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Context struct {
	Vars   map[string]string
	dir    string
	dryRun bool
}

func NewContext(options ...ContextOption) *Context {
	context := &Context{
		Vars: map[string]string{},
		dir:  ".",
	}

	for _, option := range options {
		option(context)
	}

	return context
}

func (c *Context) DryRun() bool {
	return c.dryRun
}

func (c *Context) Dir() string {
	return c.dir
}

func (c *Context) Init() ([]byte, error) {
	return c.Terraformf("init")
}

func (c *Context) Apply() ([]byte, error) {
	if c.dryRun {
		return c.Terraformf("plan")
	}

	return c.Terraformf("apply")
}

func (c *Context) Destroy() ([]byte, error) {
	if c.dryRun {
		return c.Terraformf("plan", "-destroy")
	}

	return c.Terraformf("destroy", "-force")
}

func (c *Context) Import(resource, id string) ([]byte, error) {
	return c.Terraformf("import", resource, id)
}

func (c *Context) Output(name string) (string, error) {
	output, err := c.Terraformf("output", name)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(output), "\n"), nil
}

func (c *Context) Terraformf(command string, args ...string) ([]byte, error) {
	// configure terraform variables using 'TF_VAR_<name>'
	// see: https://www.terraform.io/docs/configuration/variables.html
	env := []string{}
	for name, val := range c.Vars {
		env = append(env, fmt.Sprintf("TF_VAR_%s=%s", name, val))
	}

	args = append([]string{command}, args...)
	cmd := exec.Command("terraform", args...)
	cmd.Env = env
	cmd.Dir = c.dir

	log.Printf("[DEBUG] Running %v from %s", cmd.Args, cmd.Dir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] Error running %v from %s: %v\n", cmd.Args, cmd.Dir, err)
		return nil, err
	}

	return output, nil
}

package command

import (
	"github.com/urfave/cli"
	"testing"
)

func extractAction(t *testing.T, command cli.Command) func(*cli.Context) error {
	action, ok := command.Action.(func(*cli.Context) error)
	if !ok {
		t.Fatal("Failed to convert command.Action")
	}

	return action
}

package command

import (
	"testing"

	"github.com/urfave/cli"
)

func extractAction(t *testing.T, command cli.Command) func(*cli.Context) error {
	action, ok := command.Action.(func(*cli.Context) error)
	if !ok {
		t.Fatal("Failed to convert command.Action")
	}

	return action
}

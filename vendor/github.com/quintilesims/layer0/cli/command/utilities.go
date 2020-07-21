package command

import (
	"strings"
	"time"

	"github.com/urfave/cli"
)

func assertSingleID(entityType, target string, ids []string) (string, error) {
	if len(ids) == 0 {
		return "", noMatchesError(entityType, target)
	}

	if len(ids) > 1 {
		return "", multipleMatchesError(entityType, target, ids)
	}

	return ids[0], nil
}

func entityTitle(entityType string) string {
	title := strings.Replace(entityType, "_", " ", -1)
	return strings.Title(title)
}

func wrapAction(cm *Command, action func(c *cli.Context) error) func(*cli.Context) {
	return func(c *cli.Context) {
		if err := action(c); err != nil {
			cm.handleError(c, err)
		}
	}
}

func extractArgs(received []string, names ...string) (map[string]string, error) {
	args := map[string]string{}
	for i, name := range names {
		if len(received)-1 < i {
			return nil, NewUsageError("Argument %s is required", name)
		}

		args[name] = received[i]
	}

	return args, nil
}

func getTimeout(c *cli.Context) (time.Duration, error) {
	return time.ParseDuration(c.GlobalString("timeout"))
}

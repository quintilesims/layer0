package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

type NoMatchError struct {
	error
}

func noMatchesError(entityType, target string) error {
	err := fmt.Errorf("%s lookup using '%s' yielded no matches.", entityTitle(entityType), target)
	return NoMatchError{err}
}

func multipleMatchesError(entityType, target string, ids []string) error {
	text := fmt.Sprintf("%s lookup using '%s' yielded multiple matches: \n", entityTitle(entityType), target)
	for _, id := range ids {
		text += fmt.Sprintf("%s \n", id)
	}

	return fmt.Errorf(text)
}

type UsageError struct {
	error
}

func NewUsageError(format string, tokens ...interface{}) *UsageError {
	return &UsageError{
		error: fmt.Errorf(format, tokens...),
	}
}

func handleUsageError(c *cli.Context, err error) {
	fmt.Printf("Incorrect Usage: %s \n", err.Error())
	cli.ShowSubcommandHelp(c)
	os.Exit(1)
}

func errorSuggestion(err error) (string, bool) {
	errorContains := func(text string) bool { return strings.Contains(err.Error(), text) }

	switch {
	case errorContains("Access denied for user 'layer0_api'"):
		text := "It appears your Layer0 API service database hasn't been configured.\n"
		text += "Have you tried running ./l0 admin sql?"
		return text, true
	case errorContains("Layer0 API returned invalid status code: 401 Unauthorized"):
		text := "It appears your Layer0 CLI is using invalid credentials.\n"
		text += "Have you run ./l0-setup endpoint <instance>?"
		return text, true
	case errorContains("Unable to connect to API with error"):
		text := fmt.Sprintf("%s\n", err.Error())
		if errorContains("localhost:9090") {
			text += "\nIt appears you may not have set your LAYER0_API_ENDPOINT environment variable."
		}
		text += "\nHave you run ./l0-setup endpoint <instance>?"
		return text, true
	}

	switch err.(type) {
	case NoMatchError:
		text := err.Error()
		text += "\nDid you forget to add a wildcard ('*') to the end of your query?"
		return text, true
	}

	return "", false
}

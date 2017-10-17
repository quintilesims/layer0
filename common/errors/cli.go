package errors

import "fmt"

func NoMatchesError(entityType, target string) error {
	text := fmt.Sprintf("%s lookup using '%s' yielded no matches.\n", entityType, target)
	text += "Did you forget to add a wildcard ('*') to your target?"
	return fmt.Errorf(text)
}

func MultipleMatchesError(entityType, target string, entityIDs []string) error {
	text := fmt.Sprintf("%s lookup using '%s' yielded multiple matches: \n", entityType, target)
	for _, entityID := range entityIDs {
		text += fmt.Sprintf("%s\n", entityID)
	}

	return fmt.Errorf(text)
}

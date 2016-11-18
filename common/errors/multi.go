package errors

import (
	"fmt"
)

func MultiError(errors []error) error {
	if count := len(errors); count == 0 {
		return nil
	} else if count == 1 {
		return errors[0]
	}

	text := "Multiple Errors: \n"
	for _, err := range errors {
		text += fmt.Sprintf("\t%s\n", err.Error())
	}

	return fmt.Errorf(text)
}

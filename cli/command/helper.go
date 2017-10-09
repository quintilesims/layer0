package command

import (
	"fmt"

	"github.com/quintilesims/layer0/cli/resolver"
)

func extractArgs(received []string, names ...string) (map[string]string, error) {
	args := map[string]string{}
	for i, name := range names {
		if len(received)-1 < i {
			return nil, fmt.Errorf("Argument %s is required", name)
		}

		args[name] = received[i]
	}

	return args, nil
}

func resolveSingleEntityID(resolver resolver.Resolver, entityType, target string) (string, error) {
	entityIDs, err := resolver.Resolve(entityType, target)
	if err != nil {
		return "", err
	}

	switch len(entityIDs) {
	case 0:
		return "", fmt.Errorf("%s lookup using '%s' yielded no matches.", entityType, target)
	case 1:
		return entityIDs[0], nil
	default:
		text := fmt.Sprintf("%s lookup using '%s' yielded multiple matches: \n", entityType, target)
		for _, entityID := range entityIDs {
			text += fmt.Sprintf("%s \n", entityID)
		}

		return "", fmt.Errorf(text)
	}
}

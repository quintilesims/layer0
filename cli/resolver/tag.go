package resolver

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/errors"
)

// todo: support glob matching '*' in api
func NewTagResolver(c client.Client) ResolverFunc {
	return func(entityType, target string) ([]string, error) {
		switch entityType {
		case "environment", "job":
			return resolveUnqualifiedEntity(c, entityType, target)
		case "service", "load_balancer", "task":
			targets := strings.Split(target, ":")
			switch len(targets) {
			case 1:
				return resolveUnqualifiedEntity(c, entityType, target)
			case 2:
				return resolveFullyQualifiedEntity(c, entityType, targets[0], targets[1])
			default:
				return nil, fmt.Errorf("Invalid target format (expected ENVIRONMENT:%s)", strings.ToUpper(entityType))
			}
		case "deploy":
			targets := strings.Split(target, ":")
			switch len(targets) {
			case 1:
				return resolveUnqualifiedEntity(c, entityType, target)
			case 2:
				return resolveFullyQualifedDeploy(c, targets[0], targets[1])
			default:
				return nil, fmt.Errorf("Invalid target format (expected DEPLOY:VERSION)")
			}
		default:
			return nil, fmt.Errorf("Unrecognized entity type '%s'", entityType)
		}
	}
}

func resolveUnqualifiedEntity(c client.Client, entityType, target string) ([]string, error) {
	query := url.Values{}
	query.Set(client.TagQueryParamType, entityType)
	query.Set(client.TagQueryParamFuzz, target)

	return findMatches(c, target, query)
}

func resolveFullyQualifiedEntity(c client.Client, entityType, environmentTarget, entityTarget string) ([]string, error) {
	environmentIDs, err := resolveUnqualifiedEntity(c, "environment", environmentTarget)
	if err != nil {
		return nil, err
	}

	switch len(environmentIDs) {
	case 0:
		return nil, errors.NoMatchesError("environment", environmentTarget)
	case 1:
		query := url.Values{}
		query.Set(client.TagQueryParamType, entityType)
		query.Set(client.TagQueryParamFuzz, entityTarget)
		query.Set(client.TagQueryParamEnvironmentID, environmentIDs[0])

		return findMatches(c, entityTarget, query)
	default:
		return nil, errors.MultipleMatchesError("environment", environmentTarget, environmentIDs)
	}
}

func resolveFullyQualifedDeploy(c client.Client, target, version string) ([]string, error) {
	query := url.Values{}
	query.Set(client.TagQueryParamType, "deploy")
	query.Set(client.TagQueryParamFuzz, target)
	query.Set(client.TagQueryParamVersion, version)

	return findMatches(c, target, query)
}

func findMatches(c client.Client, target string, query url.Values) ([]string, error) {
	tags, err := c.ListTags(query)
	if err != nil {
		return nil, err
	}

	// check if we have an exact ID match
	if tag, ok := tags.WithID(target).First(); ok {
		return []string{tag.EntityID}, nil
	}

	uniqueEntityIDs := map[string]bool{}
	for _, tag := range tags {
		uniqueEntityIDs[tag.EntityID] = true
	}

	entityIDs := make([]string, 0, len(uniqueEntityIDs))
	for entityID := range uniqueEntityIDs {
		entityIDs = append(entityIDs, entityID)
	}

	return entityIDs, nil
}
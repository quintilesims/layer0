package command

import (
	"fmt"
	"strings"

	"github.com/quintilesims/layer0/cli/client"
	"github.com/quintilesims/layer0/common/models"
)

type Resolver interface {
	Resolve(entityType, target string) ([]string, error)
}

type TagResolver struct {
	client client.Client
}

func NewTagResolver(client client.Client) *TagResolver {
	return &TagResolver{
		client: client,
	}
}

func (r *TagResolver) Resolve(entityType, target string) ([]string, error) {
	var resolveFunc func(entityType, target string) ([]string, error)

	switch entityType {
	case "environment", "certificate", "job":
		resolveFunc = r.resolveGlobalScope
	case "service", "load_balancer", "task":
		resolveFunc = r.resolveEnvironmentScope
	case "deploy":
		resolveFunc = r.resolveDeploy
	default:
		return nil, fmt.Errorf("Unrecognized entity type '%s'", entityType)
	}

	ids, err := resolveFunc(entityType, target)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, noMatchesError(entityType, target)
	}

	return ids, nil
}

func (r *TagResolver) resolveGlobalScope(entityType, target string) ([]string, error) {
	return r.query(entityType, target, nil)
}

func (r *TagResolver) resolveEnvironmentScope(entityType string, target string) ([]string, error) {
	targets := strings.Split(target, ":")
	if len(targets) == 1 {
		return r.resolveGlobalScope(entityType, target)
	}

	if len(targets) > 2 {
		return nil, fmt.Errorf("Invalid target format (expected ENVIRONMENT:%s)", strings.ToUpper(entityType))
	}

	environmentIDs, err := r.resolveGlobalScope("environment", targets[0])
	if err != nil {
		return nil, err
	}

	environmentID, err := assertSingleID("environment", targets[0], environmentIDs)
	if err != nil {
		return nil, err
	}

	extraParams := map[string]string{"environment_id": environmentID}
	return r.query(entityType, targets[1], extraParams)
}

func (r *TagResolver) resolveDeploy(entityType, target string) ([]string, error) {
	extraParams := map[string]string{}

	targets := strings.Split(target, ":")
	if len(targets) == 2 {
		extraParams["version"] = targets[1]
	}

	if len(targets) > 2 {
		return nil, fmt.Errorf("Invalid target format (expected DEPLOY[:VERSION])")
	}

	return r.query(entityType, targets[0], extraParams)
}

func (r *TagResolver) query(entityType string, target string, extraParams map[string]string) ([]string, error) {
	requireExactMatches := !strings.Contains(target, "*")
	if err := cleanTarget(&target); err != nil {
		return nil, err
	}

	params := map[string]string{
		"type": entityType,
		"fuzz": target,
	}

	for k, v := range extraParams {
		params[k] = v
	}

	tags, err := r.client.SelectByQuery(params)
	if err != nil {
		return nil, err
	}

	matcher := matchAnything()
	if requireExactMatches {
		matcher = matchExact(target)
	}

	return extractIDs(tags, matcher), nil
}

func cleanTarget(target *string) error {
	*target = strings.TrimSuffix(*target, "*")
	if strings.Contains(*target, "*") {
		return NewUsageError("Wildcard matching ('*') can only be done at the end of a target.")
	}

	return nil
}

func extractIDs(tags []*models.EntityWithTags, shouldInclude func(*models.EntityWithTags) bool) []string {
	ids := []string{}
	for _, tag := range tags {
		if shouldInclude(tag) {
			ids = append(ids, tag.EntityID)
		}
	}

	return ids
}

func matchAnything() func(*models.EntityWithTags) bool {
	return func(*models.EntityWithTags) bool {
		return true
	}
}

func matchExact(target string) func(*models.EntityWithTags) bool {
	return func(ewt *models.EntityWithTags) bool {
		if ewt.EntityID == target {
			return true
		}

		for _, tag := range ewt.Tags {
			if tag.Key == "name" && tag.Value == target {
				return true
			}
		}

		return false
	}
}

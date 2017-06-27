package main

import (
	"fmt"

	"github.com/quintilesims/layer0/cli/client"
)

func resolveTags(client client.Client, target string, params map[string]string) (string, error) {
	params["fuzz"] = target

	taggedEntities, err := client.SelectByQuery(params)
	if err != nil {
		return "", err
	}

	if len(taggedEntities) == 0 {
		return "", fmt.Errorf("No entities found")
	}

	if len(taggedEntities) == 1 {
		return taggedEntities[0].EntityID, nil
	}

	entityIDs := []string{}
	for _, taggedEntity := range taggedEntities {
		entityIDs = append(entityIDs, taggedEntity.EntityID)

		for _, tag := range taggedEntity.Tags {
			if tag.Key == "name" && tag.Value == target {
				return taggedEntity.EntityID, nil
			}
		}
	}

	text := fmt.Sprintf("Multiple entities found: \n")
	for _, id := range entityIDs {
		text += fmt.Sprintf("%s \n", id)
	}

	return "", fmt.Errorf(text)
}

func setResourceData(setter func(string, interface{}) error, values map[string]interface{}) error {
	for key, value := range values {
		err := setter(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

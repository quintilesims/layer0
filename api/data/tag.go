package data

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"regexp"
	"strconv"
)

type TagData interface {
	GetTags(filter map[string]string) ([]models.EntityWithTags, error)
	Find(entityType string, filter map[string]string) ([]models.EntityWithTags, error)
	Make(tag models.EntityTag) error
	Delete(tag models.EntityTag) error
}

type TagLogicLayer struct {
	DataStore TagDataStore
}

func NewTagLogicLayer(dataStore TagDataStore) *TagLogicLayer {
	return &TagLogicLayer{
		DataStore: dataStore,
	}
}

func (this *TagLogicLayer) GetTags(filter map[string]string) ([]models.EntityWithTags, error) {
	// with no additional filters, list all the tags at this entity
	var filtered_set map[string][]models.EntityTag = nil

	if len(filter) == 0 {
		entities, err := this.DataStore.Select()
		if err != nil {
			return nil, err
		}
		filtered_set = toMap(entities)
	} else {
		for key, value := range filter {
			entities, err := this.selectByKey(key, value)
			if err != nil {
				return nil, err
			}

			entity_map := toMap(entities)
			if filtered_set == nil {
				filtered_set = entity_map
			} else {
				filtered_set = intersect(filtered_set, entity_map)
			}

			// early exit if not tags were found
			if len(filtered_set) == 0 {
				break
			}
		}

		// special case for version.
		if filter["version"] == "latest" {
			filtered_set = filterLatest(filtered_set)
		}
	}

	return toArray(filtered_set), nil
}

func validateType(entityType string) error {
	if _, ok := AllowedTagMap()[entityType]; !ok {
		log.Infof("Unexpected entityType: %s", entityType)
		err := fmt.Errorf("Entity type must be one of %v", AllowedTagTypes)
		return errors.New(errors.InvalidEntityType, err)
	}

	return nil
}

func (this *TagLogicLayer) selectByKey(key, value string) ([]models.EntityTag, error) {
	switch key {
	case "type":
		if err := validateType(value); err != nil {
			return nil, err
		}
		return this.DataStore.SelectByType(value)
	case "name_prefix":
		return this.DataStore.SelectByTagPrefix("name", value)
	case "id_prefix":
		return this.DataStore.SelectByIdPrefix(value)
	case "id":
		return this.DataStore.SelectById(value)
	case "version":
		if value == "latest" {
			return this.DataStore.SelectByTagKey(key)
		} else {
			return this.DataStore.SelectByTag(key, value)
		}
	case "fuzz":
		return this.fuzzyMatch(value)
	default:
		return this.DataStore.SelectByTag(key, value)
	}
}

// match first by name, if none fit - match by id.
func (this *TagLogicLayer) fuzzyMatch(value string) ([]models.EntityTag, error) {
	matchedNames, err := this.selectByKey("name_prefix", value)
	if err != nil {
		return nil, err
	}

	matchedIDs, err := this.selectByKey("id_prefix", value)
	if err != nil {
		return nil, err
	}

	return append(matchedNames, matchedIDs...), nil
}

// this hook exists for internal callers
func (this *TagLogicLayer) Find(entityType string, filter map[string]string) ([]models.EntityWithTags, error) {
	if filter == nil {
		filter = map[string]string{}
	}

	filter["type"] = entityType
	return this.GetTags(filter)
}

func (this *TagLogicLayer) Make(tag models.EntityTag) error {
	err := this.validateTag(tag)
	if err != nil {
		return err
	}

	// todo - check for duplicate first.
	return this.DataStore.Insert(tag)
}

var reservedTags = map[string]bool{
	"name_prefix": true,
	"id_prefix":   true,
	"fuzz":        true,
}

func (this *TagLogicLayer) validateTag(tag models.EntityTag) error {
	if _, ok := reservedTags[tag.Key]; ok {
		return errors.Newf(errors.InvalidTagKey, "%s is a reserved key, and cannot be created", tag.Key)
	}

	match, err := regexp.MatchString("^[a-zA-Z0-9_.-]*$", tag.Value)
	if err != nil {
		return err
	}

	if !match {
		return errors.Newf(errors.InvalidTagValue, "Tag values may only be letters, numbers, '.' '_' or '-'")
	}

	if len(tag.Key) > 64 {
		return errors.Newf(errors.InvalidTagKey, "Tag Key may be at most 64 characters")
	}

	if len(tag.Value) > 64 {
		return errors.Newf(errors.InvalidTagValue, "Tag Value may be at most 64 characters")
	}

	if tag.Key == "version" {
		if _, err := strconv.ParseInt(tag.Value, 10, 64); err != nil {
			return errors.Newf(errors.InvalidTagValue, "Version must be an integer")
		}
	}

	return nil
}

func (this *TagLogicLayer) Delete(tag models.EntityTag) error {
	return this.DataStore.Delete(tag)
}

func toMap(tags []models.EntityTag) map[string][]models.EntityTag {
	result := make(map[string][]models.EntityTag)

	for _, t := range tags {
		key := fmt.Sprintf("%s~%s", t.EntityType, t.EntityID)
		array, ok := result[key]
		if !ok {
			array = make([]models.EntityTag, 0, 1)
		}

		result[key] = append(array, t)
	}

	return result
}

func intersect(tags_left, tags_right map[string][]models.EntityTag) map[string][]models.EntityTag {
	result := make(map[string][]models.EntityTag)

	for key, value := range tags_left {
		if right_val, ok := tags_right[key]; ok {
			// todo(jett) prevent duplicates
			result[key] = append(value, right_val...)
		}
	}

	return result
}

func toArray(tags map[string][]models.EntityTag) []models.EntityWithTags {
	result := make([]models.EntityWithTags, len(tags))
	index := 0
	for _, value := range tags {
		if len(value) < 1 {
			log.Warnf("Empty tag array in toEntityWithTag")
			continue
		}
		first := value[0]
		result[index] = models.EntityWithTags{
			EntityID:   first.EntityID,
			EntityType: first.EntityType,
			Tags:       removeDuplicates(value),
		}
		index++
	}

	return result
}

func removeDuplicates(tags []models.EntityTag) []models.EntityTag {
	// assume that the serviceID and serviceTag already match
	temp_map := make(map[string]models.EntityTag)
	for _, t := range tags {
		key := fmt.Sprintf("%s~%s", t.Key, t.Value)
		temp_map[key] = t
	}

	result := make([]models.EntityTag, len(temp_map), len(temp_map))
	index := 0
	for _, val := range temp_map {
		result[index] = val
		index++
	}

	return result
}

func filterLatest(tags map[string][]models.EntityTag) map[string][]models.EntityTag {
	// find the latest version number
	var maxVersion int64 = 0
	for _, tagSet := range tags {
		for _, t := range tagSet {
			if t.Key == "version" {
				val, err := strconv.ParseInt(t.Value, 10, 64)
				if err != nil {
					// non integer version shouldn't happen, but skip for now
					log.Warnf("Unexpected version tag: %v", t)
					continue
				}

				if val > maxVersion {
					maxVersion = val
				}
			}
		}
	}

	// keep only values with that maxVersion
	result := make(map[string][]models.EntityTag)
	for key, tagSet := range tags {
		for _, t := range tagSet {
			if t.Key == "version" {
				val, err := strconv.ParseInt(t.Value, 10, 64)
				if err != nil {
					continue
				}

				if val == maxVersion {
					result[key] = tagSet
				}
			}
		}
	}

	return result
}

var AllowedTagTypes = []string{"service", "deploy", "environment", "certificate", "addon", "load_balancer", "task", "job"}

func AllowedTagMap() map[string]string {
	allowedTagMap := map[string]string{}
	for _, t := range AllowedTagTypes {
		allowedTagMap[t] = t
	}

	return allowedTagMap
}

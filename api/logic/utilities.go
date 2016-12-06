package logic

import (
	"github.com/quintilesims/layer0/common/models"
)

func rangeTags(entityTags []models.EntityWithTags) []models.EntityTag {
	tags := []models.EntityTag{}

	for _, et := range entityTags {
		for _, tag := range et.Tags {
			tags = append(tags, tag)
		}
	}

	return tags
}

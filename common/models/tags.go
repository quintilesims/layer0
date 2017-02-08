package models

import (
	"fmt"
)

type Tags []*Tag

type filter func(Tag) bool

func (t Tags) RemoveIf(f filter) Tags {
	cp := make(Tags, len(t))
	copy(cp, t)

	for i := 0; i < len(cp); i++ {
		if f(*cp[i]) {
			cp = append(cp[:i], cp[i+1:]...)
			i--
		}
	}

	return cp
}

func (t Tags) WithKey(key string) Tags {
	return t.RemoveIf(func(t Tag) bool {
		return t.Key != key
	})
}

func (t Tags) WithValue(value string) Tags {
	return t.RemoveIf(func(t Tag) bool {
		return t.Value != value
	})
}

func (t Tags) First() *Tag {
	if len(t) > 0 {
		return t[0]
	}

	return nil
}

func (t Tags) Any(f filter) bool {
	for _, tag := range t {
		if f(*tag) {
			return true
		}
	}

	return false
}

func (t Tags) GroupByEntity() EntitiesWithTags {
	catalog := map[string]Tags{}

	for i, tag := range t {
		key := fmt.Sprintf("%s%s", tag.EntityID, tag.EntityType)
		if _, ok := catalog[key]; !ok {
			catalog[key] = Tags{}
		}

		catalog[key] = append(catalog[key], t[i])
	}

	ewts := EntitiesWithTags{}
	for key, tags := range catalog {
		ewt := &EntityWithTags{
			EntityID:   tags[0].EntityID,
			EntityType: tags[0].EntityType,
			Tags:       catalog[key],
		}

		ewts = append(ewts, ewt)
	}

	return ewts
}

// for testing
func (t Tags) String() string {
	pt := make([]Tag, len(t))
	for i := 0; i < len(t); i++ {
		pt[i] = *t[i]
	}

	return fmt.Sprintf("%s", pt)
}

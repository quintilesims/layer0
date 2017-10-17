package models

type Tags []Tag

type filter func(Tag) bool

func (t Tags) RemoveIf(f filter) Tags {
	cp := make(Tags, len(t))
	copy(cp, t)

	for i := 0; i < len(cp); i++ {
		if f(cp[i]) {
			cp = append(cp[:i], cp[i+1:]...)
			i--
		}
	}

	return cp
}

func (t Tags) WithType(entityType string) Tags {
	return t.RemoveIf(func(t Tag) bool {
		return t.EntityType != entityType
	})
}

func (t Tags) WithID(entityID string) Tags {
	return t.RemoveIf(func(t Tag) bool {
		return t.EntityID != entityID
	})
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

func (t Tags) First() (Tag, bool) {
	if len(t) > 0 {
		return t[0], true
	}

	return Tag{}, false
}

func (t Tags) Any(f filter) bool {
	for _, tag := range t {
		if f(tag) {
			return true
		}
	}

	return false
}

func (t Tags) GroupByID() map[string]Tags {
	entityTags := map[string]Tags{}
	for _, tag := range t {
		tags, ok := entityTags[tag.EntityID]
		if !ok {
			tags = Tags{}
		}

		entityTags[tag.EntityID] = append(tags, tag)
	}

	return entityTags
}

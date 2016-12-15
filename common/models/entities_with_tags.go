package models

type EntitiesWithTags []*EntityWithTags

type ewtFilter func(e EntityWithTags) bool

func (e EntitiesWithTags) RemoveIf(f ewtFilter) EntitiesWithTags {
	cp := make(EntitiesWithTags, len(e))
	copy(cp, e)

	for i := 0; i < len(cp); i++ {
		if f(*cp[i]) {
			cp = append(cp[:i], cp[i+1:]...)
			i--
		}
	}

	return cp
}

// removes each EntityWithTags object from e if e.Tags does
// not contain at least one tag with the specified key
func (e EntitiesWithTags) WithKey(key string) EntitiesWithTags {
	return e.RemoveIf(func(ewt EntityWithTags) bool {
		hasKey := ewt.Tags.Any(func(t Tag) bool {
			return t.Key == key
		})

		return !hasKey
	})
}

// removes each EntityWithTags object from e if e.Tags does
// not contain at least one tag with the specified value
func (e EntitiesWithTags) WithValue(value string) EntitiesWithTags {
	return e.RemoveIf(func(ewt EntityWithTags) bool {
		hasKey := ewt.Tags.Any(func(t Tag) bool {
			return t.Value == value
		})

		return !hasKey
	})
}

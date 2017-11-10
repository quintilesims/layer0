package models

type LinkTags []string

func (l LinkTags) Contains(value string) bool {
	for _, v := range l {
		if v == value {
			return true
		}
	}

	return false
}

func (l LinkTags) Distinct() []string {
	m := map[string]int{}

	for idx, k := range l {
		m[k] = idx
	}

	d := []string{}
	for k, _ := range m {
		d = append(d, k)
	}

	return d
}

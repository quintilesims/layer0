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

	d := make([]string, len(m))
	i := 0
	for k, _ := range m {
		d[i] = k
		i++
	}

	return d
}

package command

import (
	"fmt"
	"net/url"
	"strconv"
)

func extractArgs(received []string, names ...string) (map[string]string, error) {
	args := map[string]string{}
	for i, name := range names {
		if len(received)-1 < i {
			return nil, fmt.Errorf("Argument %s is required", name)
		}

		args[name] = received[i]
	}

	return args, nil
}

func buildQueryHelper(id, start, end string, tail int) url.Values {
	query := url.Values{}

	if tail > 0 {
		query.Set("tail", strconv.Itoa(tail))
	}

	if start != "" {
		query.Set("start", start)
	}

	if end != "" {
		query.Set("end", end)
	}

	return query
}

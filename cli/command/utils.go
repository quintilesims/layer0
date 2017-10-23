package command

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/quintilesims/layer0/client"
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

func buildLogQueryHelper(start, end string, tail int) url.Values {
	query := url.Values{}

	if tail > 0 {
		query.Set(client.LogQueryParamTail, strconv.Itoa(tail))
	}

	if start != "" {
		query.Set(client.LogQueryParamStart, start)
	}

	if end != "" {
		query.Set(client.LogQueryParamEnd, end)
	}

	return query
}

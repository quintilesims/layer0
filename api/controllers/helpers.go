package controllers

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// 'YYYY-MM-DD HH:MM' time layout as described by https://golang.org/src/time/format.go
const TimeLayout = "2006-01-02 15:04"

func ParseLoggingQuery(query url.Values) (int, time.Time, time.Time, error) {
	var tail int
	if v := query.Get("tail"); v != "" {
		t, err := strconv.Atoi(v)
		if err != nil {
			return 0, time.Time{}, time.Time{}, fmt.Errorf("Tail must be an integer")
		}

		tail = t
	}

	parseTime := func(v string) (time.Time, error) {
		if v == "" {
			return time.Time{}, nil
		}

		return time.Parse(TimeLayout, v)
	}

	start, err := parseTime(query.Get("start"))
	if err != nil {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("Invalid time: start must be in format YYYY-MM-DD HH:MM")
	}

	end, err := parseTime(query.Get("end"))
	if err != nil {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("Invalid time: end must be in format YYYY-MM-DD HH:MM")
	}

	return tail, start, end, nil
}

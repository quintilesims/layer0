package controllers

import (
	"net/url"
	"strconv"
	"time"

	"github.com/quintilesims/layer0/common/errors"
)

func parseLoggingQuery(query url.Values) (int, time.Time, time.Time, error) {
	var tail int
	if v := query.Get("tail"); v != "" {
		t, err := strconv.Atoi(v)
		if err != nil {
			return 0, time.Time{}, time.Time{}, errors.Newf(errors.InvalidRequest, "Tail must be an integer")
		}

		tail = t
	}

	parseTime := func(v string) (time.Time, error) {
		if v == "" {
			return time.Time{}, nil
		}

		return time.Parse(TIME_LAYOUT, v)
	}

	start, err := parseTime(query.Get("start"))
	if err != nil {
		return 0, time.Time{}, time.Time{}, errors.Newf(errors.InvalidRequest, "Invalid time: start must be in format YYYY-MM-DD HH:MM")
	}

	end, err := parseTime(query.Get("end"))
	if err != nil {
		return 0, time.Time{}, time.Time{}, errors.Newf(errors.InvalidRequest, "Invalid time: end must be in format YYYY-MM-DD HH:MM")
	}

	return tail, start, end, nil
}

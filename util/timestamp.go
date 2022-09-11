package util

import "time"

func ParseTimestamp(str string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

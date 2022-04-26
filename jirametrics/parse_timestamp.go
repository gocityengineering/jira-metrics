package jirametrics

import (
	"time"
)

func parseTimestamp(timestamp string) (int64, error) {
	// Jira timestamp format in Go notation
	t, err := time.Parse("2006-01-02T15:04:05.000-0700", timestamp)
	if err != nil {
		return time.Now().Unix(), err
	}

	return t.Unix(), nil
}

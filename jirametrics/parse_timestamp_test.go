package jirametrics

import (
	"testing"
)

func TestParseTimestamp(t *testing.T) {
	var tests = []struct {
		description   string
		jiraTimestamp string
		expected      int64
	}{
		{"timestamp01", "2001-09-28T07:09:05.601-0400", 1001675345},
		{"timestamp02", "2022-03-04T07:11:11.087-0500", 1646395871},
		{"timestamp03", "1999-10-04T05:14:51.487-0400", 939028491},
		{"timestamp04", "2018-11-22T04:27:03.503-0500", 1542878823},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			unix, _ := parseTimestamp(test.jiraTimestamp)
			if unix != test.expected {
				t.Errorf("Unexpected result %d for input %s; expected %d: %s", unix, test.jiraTimestamp, test.expected, test.description)
			}
		})
	}
}

package jirametrics

import (
	"testing"
)

func TestDateStampOufOfHours(t *testing.T) {
	var tests = []struct {
		description string
		unixTime    int64
		timeZone    string
		expected    bool
	}{
		{"9:30 am UTC", 9.5 * 3600, "UTC", false},
		{"7:30 am UTC", 7.5 * 3600, "UTC", true},
		{"11:00 pm UTC", 23 * 3600, "UTC", true},
		{"9:00 pm UTC in New York", 21 * 3600, "America/New_York", false},
		{"6:30 am UTC in London", 6.5 * 3600, "Europe/London", true},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			b, _ := dateStampOutOfHours(test.unixTime, test.timeZone)
			if b != test.expected {
				t.Errorf("Unexpected result '%t' for Unix time %ds in time zone %s: %s", b, test.unixTime, test.timeZone, test.description)
			}
		})
	}
}

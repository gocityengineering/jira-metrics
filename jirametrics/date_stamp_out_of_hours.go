package jirametrics

import "time"

func dateStampOutOfHours(unixTime int64, timeZone string) (bool, error) {
	t := time.Unix(unixTime, 0)
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return false, err
	}

	localTime := t.In(location)

	// weekend always out of hours
	if localTime.Weekday() == time.Saturday ||
		localTime.Weekday() == time.Sunday ||
		localTime.Hour() < 8 ||
		localTime.Hour() > 16 {
		return true, nil
	}

	return false, nil
}

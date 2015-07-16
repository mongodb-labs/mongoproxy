package bi

import (
	"fmt"
	"time"
)

// GetSuffix returns the proper suffix for the collection to store metrics
// with the given time granularity. Returns an error if the granularity provided
// is not a valid time granularity.
func GetSuffix(granularity string) (string, error) {
	switch granularity {
	case Monthly:
		return "-month", nil
	case Daily:
		return "-day", nil
	case Hourly:
		return "-hour", nil
	case Minutely:
		return "-minute", nil
	case Secondly:
		return "-second", nil
	default:
		return "", fmt.Errorf("Not a valid time granularity")
	}
}

// GetStartTime takes a time and a granularity, and rounds it up to the next
// highest time granularity (or the proper start time for the metric document). Returns
// an error if the granularity provided is not a valid time granularity.
func GetStartTime(t time.Time, granularity string) (time.Time, error) {
	var start time.Time
	switch granularity {
	case Monthly:
		start = time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
	case Daily:
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case Hourly:
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0,
			t.Location())
	case Minutely:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0,
			t.Location())
	case Secondly:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
			t.Minute(), 0, 0, t.Location())
	default:
		return start, fmt.Errorf("Not a valid time granularity")
	}

	return start, nil
}

// GetRoundedTime takes a time and a granularity, and returns the time rounded down
// to that particular granularity. Returns an error if the granularity provided
// is not a valid time granularity.
func GetRoundedTime(t time.Time, granularity string) (time.Time, error) {
	var start time.Time
	switch granularity {
	case Monthly:
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case Daily:
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0,
			0, t.Location())
	case Hourly:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
			0, 0, 0, t.Location())
	case Minutely:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
			t.Minute(), 0, 0, t.Location())
	case Secondly:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
			t.Minute(), t.Second(), 0, t.Location())
	default:
		return start, fmt.Errorf("Not a valid time granularity")
	}

	return start, nil
}

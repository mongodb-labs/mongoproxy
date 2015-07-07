package bi

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

const (
	Monthly  string = "M"
	Daily    string = "D"
	Hourly   string = "h"
	Minutely string = "m"
	Secondly string = "s"
)

type Rule struct {
	OriginDatabase    string
	OriginCollection  string
	PrefixDatabase    string
	PrefixCollection  string
	TimeGranularities []string
	ValueField        string
	TimeField         *time.Time
}

func getSuffix(granularity string) (string, error) {
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

func createSelector(t time.Time, granularity string, valueField string) (bson.D, error) {
	var start time.Time
	switch granularity {
	case Monthly:
		start = time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
	case Daily:
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case Hourly:
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case Minutely:
		start = t.Round(time.Hour)
	case Secondly:
		start = t.Round(time.Minute)
	default:
		return nil, fmt.Errorf("Not a valid time granularity")
	}

	doc := bson.D{{"start", start}, {"valueField", valueField}}

	return doc, nil
}

func createUpdate(t time.Time, granularity string, value float64) (bson.D, error) {

	var M int
	var granularityField string

	switch granularity {
	case Monthly:
		M = int(t.Month())
		granularityField = "month"
	case Daily:
		M = t.Day()
		granularityField = "day"
	case Hourly:
		M = t.Hour()
		granularityField = "hour"
	case Minutely:
		M = t.Minute()
		granularityField = "minute"
	case Secondly:
		M = t.Second()
		granularityField = "second"
	default:
		return nil, fmt.Errorf("Not a valid time granularity")
	}

	timeField := strconv.Itoa(M)

	totalUpdate := bson.DocElem{"$inc", bson.D{{"total", value}}}
	fieldUpdate := bson.DocElem{"$inc", bson.D{{granularityField + "." + timeField, value}}}

	doc := bson.D{
		totalUpdate, fieldUpdate,
	}

	return doc, nil

}

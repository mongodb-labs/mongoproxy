package bi

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
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

func createSelector(t time.Time, granularity string, valueField string) (bson.D, error) {
	var start time.Time
	switch granularity {
	case "M":
		start = time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
	case "D":
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case "h":
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case "m":
		start = t.Round(time.Hour)
	case "s":
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
	case "M":
		M = int(t.Month())
		granularityField = "month"
	case "D":
		M = t.Day()
		granularityField = "day"
	case "h":
		M = t.Hour()
		granularityField = "hour"
	case "m":
		M = t.Minute()
		granularityField = "minute"
	case "s":
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

package bi

import (
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/bsonutil"
	"github.com/mongodbinc-interns/mongoproxy/convert"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

// Constants for the time granularities.
const (
	Monthly  string = "M"
	Daily    string = "D"
	Hourly   string = "h"
	Minutely string = "m"
	Secondly string = "s"
)

// A Rule determines how the BI module processes each request.
type Rule struct {
	// The origin of requests processed by this rule.
	OriginDatabase   string
	OriginCollection string

	// The prefix of the location that the rule will store generated metrics.
	PrefixDatabase   string
	PrefixCollection string

	// A list of the time granularities for metrics that are saved.
	TimeGranularities []string

	// A field in the request that is analyzed and saved in metric documents.
	ValueField string

	// An optional time for when the request was saved.
	TimeField *time.Time
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
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	case Secondly:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
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

func createSingleUpdate(doc bson.D, time time.Time, granularity string,
	rule Rule) (*messages.SingleUpdate, error) {

	selectorRaw, err := createSelector(time, granularity, rule.ValueField)

	if err != nil {
		return nil, fmt.Errorf("Error creating selector: %v", err)
	}

	valueRaw := bsonutil.FindValueByKey(rule.ValueField, doc)
	if valueRaw == nil {
		return nil, fmt.Errorf("No value found: %v", err)
	}

	value := convert.ToFloat64(valueRaw)
	if value == 0 {
		return nil, fmt.Errorf("No need to continue, value is 0: %v", err)
	}

	// TODO: actually grab the value
	updateRaw, err := createUpdate(time, granularity, value)

	if err != nil {
		return nil, fmt.Errorf("Error creating update: %v", err)
	}

	single := messages.SingleUpdate{
		Upsert:   true,
		Selector: selectorRaw,
		Update:   updateRaw,
	}

	return &single, nil
}

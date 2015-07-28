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
	TimeField *string
}

func createSelector(t time.Time, granularity string, valueField string) (bson.D, error) {
	start, err := GetStartTime(t, granularity)
	if err != nil {
		return nil, err
	}

	doc := bson.D{{"start", start}, {"valueField", valueField}, {"value", bson.D{{"$exists", false}}}}

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

func createSelectorString(t time.Time, granularity string, valueField string,
	value string) (bson.D, error) {
	start, err := GetStartTime(t, granularity)
	if err != nil {
		return nil, err
	}

	doc := bson.D{{"start", start}, {"valueField", valueField}, {"value", value}, {"single", true}}

	return doc, nil
}

// helper function to create the upsert command for a single document matching a single
// rule at a single time granularity.
func createSingleUpdate(doc bson.D, time time.Time, granularity string,
	rule Rule) (*messages.SingleUpdate, *messages.SingleUpdate, error) {

	docMap := doc.Map()
	valueRaw := bsonutil.FindDeepValueInMap(rule.ValueField, docMap)
	if valueRaw == nil {
		return nil, nil, fmt.Errorf("No value found")
	}

	valueStr, ok := valueRaw.(string)
	if ok {
		// it's a string
		selectorStr, err := createSelectorString(time, granularity, rule.ValueField, valueStr)
		if err != nil {
			return nil, nil, fmt.Errorf("Error creating selector: %v", err)
		}
		updateStr, err := createUpdate(time, granularity, 1)
		if err != nil {
			return nil, nil, fmt.Errorf("Error creating update: %v", err)
		}

		single := messages.SingleUpdate{
			Upsert:   true,
			Selector: selectorStr,
			Update:   updateStr,
		}
		meta := saveMetadataForValue(rule, granularity, valueStr)

		return &single, &meta, nil
	}

	// otherwise, it's a float.
	valueFloat := convert.ToFloat64(valueRaw, 0)
	if valueFloat == 0 {
		return nil, nil, fmt.Errorf("No need to continue, value is 0")
	}

	selectorRaw, err := createSelector(time, granularity, rule.ValueField)

	if err != nil {
		return nil, nil, fmt.Errorf("Error creating selector: %v", err)
	}

	// TODO: actually grab the value
	updateRaw, err := createUpdate(time, granularity, valueFloat)

	if err != nil {
		return nil, nil, fmt.Errorf("Error creating update: %v", err)
	}

	single := messages.SingleUpdate{
		Upsert:   true,
		Selector: selectorRaw,
		Update:   updateRaw,
	}

	return &single, nil, nil
}

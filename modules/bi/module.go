package bi

import (
	"encoding/json"
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/bsonutil"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/modules"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
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

type BIModule struct {
	Rules []Rule
}

func createSelector(t time.Time, granularity string, valueField string) (bson.D, error) {
	var start time.Time
	switch granularity {
	case "M":
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case "D":
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case "h":
		start = t.Round(time.Hour)
	case "m":
		start = t.Round(time.Minute)
	case "s":
		start = t.Round(time.Second)
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
		M = t.Month()
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

	totalUpdate := bson.DocElem{"total", bson.D{{"$inc", value}}}
	fieldUpdate := bson.DocElem{granularityField, bson.D{{timeField, bson.D{"$inc", value}}}}
	doc := bson.D{
		totalUpdate, fieldUpdate,
	}

	return doc, nil

}
func (b BIModule) Process(req messages.Requester, res messages.Responder,
	next modules.PipelineFunc) {

	resNext := messages.ModuleResponse{}
	next(req, &resNext)
	res.Write(resNext.Writer)

	if resNext.CommandError != nil {
		return // we're done. An error occured, so we shouldn't do any aggregating
	}

	updates := make([]messages.Update, 0)

	if req.Type() == messages.InsertType {
		// create metrics
		opi := req.(messages.Insert)

		for i := 0; i < len(b.Rules); i++ {
			rule := b.Rules[i]

			time := time.Now()

			// use the time field instead if it exists
			if rule.TimeField != nil {
				time = *rule.TimeField
			}

			// if the message matches the aggregation, create an upsert
			// and pass it on to mongod
			if opi.Collection != rule.OriginCollection ||
				opi.Database != rule.OriginDatabase {
				continue
			}

			for j := 0; j < len(rule.TimeGranularities); j++ {
				granularity := rule.TimeGranularities[j]
				var suffix string
				switch granularity {
				case "M":
					suffix = "-month"
				case "D":
					suffix = "-day"
				case "h":
					suffix = "-hour"
				case "m":
					suffix = "-minute"
				case "s":
					suffix = "-second"
				default:
					continue
				}

				update := messages.Update{
					Database:   rule.PrefixDatabase,
					Collection: rule.PrefixCollection + suffix,
					Ordered:    false,
				}

				for k := 0; k < len(opi.Documents); k++ {

					doc := opi.Documents[k]

					selectorRaw, err := createSelector(time, granularity, rule.ValueField)

					valueRaw, ok := doc[rule.ValueField]
					if !ok {
						continue // no value to actually add an update for
					}

					value := convert.ToFloat64(valueRaw)
					if value == 0 {
						continue // no need for an update if the value is 0
					}

					// TODO: actually grab the value
					updateRaw, err := createUpdate(time, granularity, value)

					single := messages.SingleUpdate{
						Upsert:   true,
						Selector: selectorRaw,
						Update:   updateRaw,
					}

					update.Updates = append(update.Updates, single)

				}
				updates = append(updates, update)
			}
		}

		// TODO: convert those all to wire protocol messages and send
		// to mongod
		Log(INFO, "%#v\n", updates)

	}

}

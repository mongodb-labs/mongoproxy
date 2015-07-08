package controllers

import (
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func getDataOverRange(session *mgo.Session, rule bi.Rule, granularity string, valueField string, start time.Time, end time.Time) ([]bson.M, error) {
	// get a query for the data over that range

	// for the time range
	// the first time should be start rounded down
	// the second time should be end
	db := session.DB(rule.PrefixDatabase)

	var startRange time.Time
	collectionName := rule.PrefixCollection

	switch granularity {
	case bi.Monthly:
		startRange = time.Date(start.Year(), time.January, 1, 0, 0, 0, 0, start.Location())
		collectionName += "-month"
	case bi.Daily:
		startRange = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
		collectionName += "-day"
	case bi.Hourly:
		startRange = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		collectionName += "-hour"
	case bi.Minutely:
		startRange = start.Round(time.Hour)
		collectionName += "-minute"
	case bi.Secondly:
		startRange = start.Round(time.Minute)
		collectionName += "-second"
	default:
		return nil, fmt.Errorf("%v is not a valid time granularity", granularity)
	}

	c := db.C(collectionName)
	query := bson.M{
		"valueField": valueField,
		"start":      bson.M{"$gte": startRange, "$lte": end},
	}

	iter := c.Find(query).Iter()

	var results []bson.M

	err := iter.All(&results)

	if err != nil {
		return nil, err
	}

	return results, nil
}

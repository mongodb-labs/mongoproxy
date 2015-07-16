package controllers

import (
	"github.com/mongodbinc-interns/mongoproxy/modules/bi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// getDataOverRange is a helper function that queries the MongoDB database for
// metric documents matching the given rule and granularity between the start
// and end times.
func getDataOverRange(session *mgo.Session, rule bi.Rule, granularity string,
	start time.Time, end time.Time) ([]bson.M, error) {

	db := session.DB(rule.PrefixDatabase)

	startRange, err := bi.GetStartTime(start, granularity)
	if err != nil {
		return nil, err
	}

	collectionSuffix, err := bi.GetSuffix(granularity)
	if err != nil {
		return nil, err
	}
	collectionName := rule.PrefixCollection + collectionSuffix
	c := db.C(collectionName)

	query := bson.M{
		"valueField": rule.ValueField,
		"start":      bson.M{"$gte": startRange, "$lte": end},
	}

	// make sure the documents are in sorted order.
	iter := c.Find(query).Sort("start").Iter()

	var results []bson.M

	err = iter.All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

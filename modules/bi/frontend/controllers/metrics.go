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
	start time.Time, end time.Time, value *string) ([]bson.M, error) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	db := sessionCopy.DB(rule.PrefixDatabase)

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

	if value == nil {
		query["value"] = bson.M{"$exists": false}
	} else {
		query["value"] = value
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

// getMetadataForRule is a helper function that queries the MongoDB database for
// a metadata document matching the given rule and granularity.
func getMetadataForRule(session *mgo.Session, rule bi.Rule, granularity string) (bson.M, error) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	db := sessionCopy.DB(rule.PrefixDatabase)
	collectionSuffix, err := bi.GetSuffix(granularity)
	if err != nil {
		return nil, err
	}
	collectionName := rule.PrefixCollection + collectionSuffix
	c := db.C(collectionName)

	query := bson.M{
		"_id": "metadata",
	}
	var result bson.M
	err = c.Find(query).One(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

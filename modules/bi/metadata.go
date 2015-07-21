package bi

import (
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"gopkg.in/mgo.v2/bson"
)

func saveMetadataForValue(rule Rule, granularity string,
	value string) messages.SingleUpdate {

	selector := bson.D{{"_id", "metadata"}}
	update := bson.D{{rule.ValueField, bson.D{{value, true}}}}

	single := messages.SingleUpdate{
		Selector: selector,
		Update:   update,
		Upsert:   true,
	}
	return single
}

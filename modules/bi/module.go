package bi

import (
	"github.com/mongodbinc-interns/mongoproxy/bsonutil"
	"github.com/mongodbinc-interns/mongoproxy/convert"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/server"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type BIModule struct {
	Rules []Rule
}

var mongoSession *mgo.Session
var mongoDBDialInfo = &mgo.DialInfo{
	// TODO: Allow configurable connection info
	Addrs:    []string{"localhost:27017"},
	Timeout:  60 * time.Second,
	Database: "test",
}

func init() {
	var err error
	mongoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		Log(ERROR, "%#v\n", err)
		return
	}

	mongoSession.SetPrefetch(0)
}

func (b BIModule) Process(req messages.Requester, res messages.Responder,
	next server.PipelineFunc) {

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
				suffix, err := getSuffix(granularity)
				if err != nil {
					Log(ERROR, "%v is not a time granularity", granularity)
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

					if err != nil {
						continue
					}

					valueRaw := bsonutil.FindValueByKey(rule.ValueField, doc)
					if valueRaw == nil {
						continue // no value to actually add an update for
					}

					value := convert.ToFloat64(valueRaw)
					if value == 0 {
						Log(ERROR, "Value was 0\n")
						continue // no need for an update if the value is 0
					}

					// TODO: actually grab the value
					updateRaw, err := createUpdate(time, granularity, value)

					if err != nil {
						continue
					}

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
		Log(NOTICE, "%#v\n", updates)

		for i := 0; i < len(updates); i++ {
			u := updates[i]
			b := u.ToBSON()

			reply := bson.D{}
			mongoSession.DB(u.Database).Run(b, &reply)
			Log(NOTICE, "%#v\n", reply)
		}

	}

}

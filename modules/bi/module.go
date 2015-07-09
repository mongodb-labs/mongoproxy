// Package bi contains a real-time reporting and analytics module for the Mongo Proxy.
// It receives requests from a mongo client and creates time series data based on user-defined criteria.
package bi

import (
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/server"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// BIModule calls the next module immediately, and then collects and aggregates
// data from inserts that successfully traveled the pipeline. The requests it analyzes
// and the metrics it aggregates is based upon its rules.
type BIModule struct {
	Rules []Rule
}

// Temporary code to set up a connection with mongod. Should eventually be replaced
// by user-set configuration.
var mongoSession *mgo.Session
var mongoDBDialInfo = &mgo.DialInfo{
	// TODO: Allow configurable connection info
	Addrs:    []string{"localhost:27017"},
	Timeout:  60 * time.Second,
	Database: "test",
}

// TODO: have a specific function for configuring modules.
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
					Log(INFO, "%v is not a time granularity", granularity)
					continue
				}

				update := messages.Update{
					Database:   rule.PrefixDatabase,
					Collection: rule.PrefixCollection + suffix,
					Ordered:    false,
				}

				for k := 0; k < len(opi.Documents); k++ {

					doc := opi.Documents[k]
					single, err := createSingleUpdate(doc, time, granularity, rule)
					if err != nil {
						continue
					}

					update.Updates = append(update.Updates, *single)

				}
				updates = append(updates, update)
			}
		}

		for i := 0; i < len(updates); i++ {
			u := updates[i]
			b := u.ToBSON()

			reply := bson.D{}
			mongoSession.DB(u.Database).Run(b, &reply)
		}

	}

}

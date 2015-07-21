// Package bi contains a real-time reporting and analytics module for the Mongo Proxy.
// It receives requests from a mongo client and creates time series data based on user-defined criteria.
package bi

import (
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/convert"
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
	Rules      []Rule
	Connection mgo.DialInfo
}

var mongoSession *mgo.Session

func init() {
	server.Publish(&BIModule{})
}

func (b *BIModule) New() server.Module {
	return &BIModule{}
}

func (b *BIModule) Name() string {
	return "bi"
}

/*
Configuration structure:
{
	connection: {
		addresses: []string,
		direct: boolean,
		timeout: integer,
		auth: {
			username: string,
			password: string,
			database: string
		}
	}
	rules: [
		{
			origin: string,
			prefix: string,
			timeGranularity: []string,
			valueField: string,
			timeField: string
		}
	]
}
*/
func (b *BIModule) Configure(conf bson.M) error {

	conn := convert.ToBSONMap(conf["connection"])
	if conn == nil {
		return fmt.Errorf("No connection data")
	}
	addrs, err := convert.ConvertToStringSlice(conn["addresses"])
	if err != nil {
		return fmt.Errorf("Invalid addresses: %v", err)
	}

	timeout := time.Duration(convert.ToInt64(conn["timeout"], -1))
	if timeout == -1 {
		timeout = time.Second * 10
	}

	dialInfo := mgo.DialInfo{
		Addrs:   addrs,
		Direct:  convert.ToBool(conn["direct"]),
		Timeout: timeout,
	}

	auth := convert.ToBSONMap(conn["auth"])
	if auth != nil {
		username, ok := auth["username"].(string)
		if ok {
			dialInfo.Username = username
		}
		password, ok := auth["password"].(string)
		if ok {
			dialInfo.Password = password
		}
		database, ok := auth["database"].(string)
		if ok {
			dialInfo.Database = database
		}
	}

	b.Connection = dialInfo

	// Rules
	b.Rules = make([]Rule, 0)
	rules, err := convert.ConvertToBSONMapSlice(conf["rules"])
	if err != nil {
		return fmt.Errorf("Error parsing rules: %v", err)
	}

	for i := 0; i < len(rules); i++ {
		r := rules[i]
		originD, originC, err := messages.ParseNamespace(convert.ToString(r["origin"]))
		if err != nil {
			return fmt.Errorf("Error parsing origin namespace: %v", err)
		}
		prefixD, prefixC, err := messages.ParseNamespace(convert.ToString(r["prefix"]))
		if err != nil {
			return fmt.Errorf("Error parsing prefix namespace: %v", err)
		}
		timeGranularities, err := convert.ConvertToStringSlice(r["timeGranularity"])
		if err != nil {
			return fmt.Errorf("Error parsing time granularities: %v", err)
		}
		valueField, ok := r["valueField"].(string)
		if !ok {
			return fmt.Errorf("Invalid valueField.")
		}
		rule := Rule{
			OriginDatabase:    originD,
			OriginCollection:  originC,
			PrefixDatabase:    prefixD,
			PrefixCollection:  prefixC,
			TimeGranularities: timeGranularities,
			ValueField:        valueField,
		}
		timeField, ok := r["timeField"].(time.Time)
		if ok {
			rule.TimeField = &timeField
		} else {
			timeFieldRaw, ok := r["timeField"].(string)
			if ok {
				if len(timeFieldRaw) > 0 {
					err := timeField.UnmarshalText([]byte(timeFieldRaw))
					if err != nil {
						rule.TimeField = &timeField
					}
				}

			}
		}

		b.Rules = append(b.Rules, rule)
	}

	return nil
}

func (b *BIModule) Process(req messages.Requester, res messages.Responder,
	next server.PipelineFunc) {

	resNext := messages.ModuleResponse{}
	next(req, &resNext)

	res.Write(resNext.Writer)

	if resNext.CommandError != nil {
		res.Error(resNext.CommandError.ErrorCode, resNext.CommandError.Message)
		return // we're done. An error occured, so we shouldn't do any aggregating
	}

	// spin up the session if it doesn't exist
	if mongoSession == nil {
		var err error
		mongoSession, err = mgo.DialWithInfo(&b.Connection)
		if err != nil {
			Log(ERROR, "%#v\n", err)
			return
		}
		mongoSession.SetPrefetch(0)
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
				Log(DEBUG, "Didn't match database %v.%v. Was %v.%v", rule.OriginDatabase,
					rule.OriginCollection, opi.Database, opi.Collection)
				continue
			}

			for j := 0; j < len(rule.TimeGranularities); j++ {
				granularity := rule.TimeGranularities[j]
				suffix, err := GetSuffix(granularity)
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
					single, meta, err := createSingleUpdate(doc, time, granularity, rule)
					if err != nil {
						continue
					}

					update.Updates = append(update.Updates, *single)
					if meta != nil {
						update.Updates = append(update.Updates, *meta)
					}

				}
				updates = append(updates, update)
			}
		}

		for i := 0; i < len(updates); i++ {
			u := updates[i]
			if len(updates[i].Updates) == 0 {
				continue
			}
			b := u.ToBSON()

			reply := bson.D{}
			Log(NOTICE, "%#v", b)
			err := mongoSession.DB(u.Database).Run(b, &reply)
			if err != nil {
				Log(ERROR, "Error updating database: %v", err)
			} else {
				Log(INFO, "Successfully updated database!")
			}
		}

	}

}

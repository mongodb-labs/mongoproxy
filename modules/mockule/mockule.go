// Package mockule contains a module that can be used as a mock backend for
// proxy core, which stores inserts and queries finds in memory.
package mockule

import (
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/server"
	"gopkg.in/mgo.v2/bson"
)

var maxWireVer = 2

// a 'database' in memory. The string keys are the collections, which
// have an array of bson documents.
var database = make(map[string][]bson.D)

// The Mockule is a mock module used for testing. It currently
// logs requests and sends valid but generally nonsense responses back to
// the client, without touching mongod.
type Mockule struct {
}

func (m Mockule) Process(req messages.Requester, res messages.Responder,
	next server.PipelineFunc) {

	Log(DEBUG, req.Type())
	switch req.Type() {
	case messages.FindType:
		opq := req.(messages.Find)
		Log(INFO, "%#v\n", opq)

		// TODO: actually do something with the query

		// grab the documents from the 'database'. We don't care about
		// the queries at the moment
		r := messages.FindResponse{}
		docs, ok := database[opq.Collection]
		r.Documents = docs
		if !ok {
			r.Documents = make([]bson.D, 0)
		}
		r.Database = opq.Database
		r.Collection = opq.Collection
		res.Write(r)
	case messages.GetMoreType:
		opg := req.(messages.GetMore)
		Log(INFO, "%#v\n", opg)
		r := messages.GetMoreResponse{}
		if opg.CursorID == int64(100) {
			Log(NOTICE, "Retrieved valid getMore\n")
		}
		r.Database = opg.Database
		r.Collection = opg.Collection
		r.CursorID = opg.CursorID
		r.Documents = make([]bson.D, 0)
		res.Write(r)
	case messages.InsertType:
		opi := req.(messages.Insert)
		Log(INFO, "%#v\n", opi)

		// insert documents into the 'database'
		for doc := range opi.Documents {
			_, ok := database[opi.Collection]
			if !ok {
				database[opi.Collection] = make([]bson.D, 0)
			}
			database[opi.Collection] = append(database[opi.Collection], opi.Documents[doc])
		}

		r := messages.InsertResponse{}
		r.N = int32(len(opi.Documents))

		res.Write(r)
	case messages.UpdateType:
		opu := req.(messages.Update)
		r := messages.UpdateResponse{}
		Log(INFO, "%#v\n", opu)
		r.N = 5
		r.NModified = 4

		// res.Write(r)
		res.Error(0, "not supported")
	case messages.DeleteType:
		opd := req.(messages.Delete)
		Log(INFO, "%#v\n", opd)
		r := messages.DeleteResponse{}
		r.N = 1

		res.Write(r)
		res.Error(0, "not supported")
	case messages.CommandType:
		command := req.(messages.Command)
		Log(INFO, "%#v\n", command)

		switch command.CommandName {
		case "isMaster":
			r := bson.M{}
			r["ismaster"] = true
			r["secondary"] = false
			r["localTime"] = bson.Now()
			r["maxWireVersion"] = maxWireVer
			r["minWireVersion"] = 0
			reply := messages.CommandResponse{}
			reply.Reply = r
			res.Write(reply)
			return
		case "whatsmyuri":
			r := bson.M{}
			r["ok"] = 1
			r["you"] = "localhost"
			reply := messages.CommandResponse{}
			reply.Reply = r
			res.Write(reply)
			return
		case "getLog":
			r := bson.M{}

			normalLog := make([]string, 2)
			normalLog[0] = "hello world"
			normalLog[1] = "this is strange"
			warningsLog := make([]string, 0)
			if maxWireVer < 2 {
				warningsLog = append(warningsLog, "Using the various OpCodes rather than commands.")
			}

			t := command.GetArg("getLog")
			if t == "startupWarnings" {
				r["log"] = warningsLog
			} else {
				r["log"] = normalLog
			}
			reply := messages.CommandResponse{}
			reply.Reply = r
			res.Write(reply)
			return
		case "replSetGetStatus":
			r := bson.M{}
			r["set"] = "repl"
			r["date"] = bson.Now()
			r["myState"] = 1
			members := make([]bson.M, 0)

			member := bson.M{}
			member["_id"] = 0
			member["name"] = "m1.example.net:27017"
			member["health"] = 1
			member["state"] = 1
			member["stateStr"] = "PRIMARY"
			member["self"] = true

			members = append(members, member)
			r["members"] = members

			reply := messages.CommandResponse{}
			reply.Reply = r
			res.Write(reply)
			return
		}
		reply := messages.CommandResponse{}
		reply.Reply = bson.M{"ok": 1}
		res.Write(reply)
	}
	next(req, res)
}
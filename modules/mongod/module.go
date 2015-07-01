package mongod

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

var port = 8000

// A MongodModule takes the request, sends it to a mongod instance, and then
// writes the response from mongod into the ResponseWriter before calling
// the next module. It passes on requests unchanged.
type MongodModule struct{}

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
}

func (m MongodModule) Process(req messages.Requester, res messages.Responder,
	next server.PipelineFunc) {

	switch req.Type() {
	case messages.CommandType:
		command, err := messages.ToCommandRequest(req)
		if err != nil {
			Log(ERROR, "%#v\n", err)
			next(req, res)
			return
		}

		b := commandToBSONDoc(command)

		reply := bson.M{}
		mongoSession.DB(command.Database).Run(b, reply)

		response := messages.CommandResponse{
			Reply: reply,
		}

		if convert.ToInt(reply["ok"]) == 0 {
			// we have a command error.
			res.Error(convert.ToInt32(reply["code"]), convert.ToString(reply["errmsg"]))
		}

		res.Write(response)

	case messages.FindType:
		f, err := messages.ToFindRequest(req)
		if err != nil {
			Log(ERROR, "%#v\n", err)
			next(req, res)
			return
		}

		c := mongoSession.DB(f.Database).C(f.Collection)
		query := c.Find(f.Filter).Limit(int(f.Limit)).Skip(int(f.Skip))

		if f.Projection != nil {
			query = query.Select(f.Projection)
		}

		var results []bson.D

		// TODO: implement getMore properly, since this behavior is
		// different than the default shell's. Will just dump all of the
		// documents instead of using the cursor and batches.
		err = query.All(&results)

		if err != nil {
			Log(ERROR, "Error on Find Command: %#v\n", err)

			// log an error if we can
			qErr, ok := err.(*mgo.QueryError)
			if ok {
				res.Error(int32(qErr.Code), qErr.Message)
			}
			next(req, res)
			return
		}

		response := messages.FindResponse{
			Database:   f.Database,
			Collection: f.Collection,
			Documents:  results,
			// TODO: retrieve CursorID
		}

		res.Write(response)

	case messages.InsertType:
		insert, err := messages.ToInsertRequest(req)
		if err != nil {
			Log(ERROR, "%#v\n", err)
			next(req, res)
			return
		}

		b := insertToBSONDoc(insert)

		reply := bson.M{}
		mongoSession.DB(insert.Database).Run(b, reply)

		response := messages.InsertResponse{
			// default to -1 if n doesn't exist to hide the field on export
			N: convert.ToInt32(reply["n"], -1),
		}
		writeErrors, err := convert.ConvertToBSONMapSlice(reply["writeErrors"])
		if err == nil {
			// we have write errors
			response.WriteErrors = writeErrors
		}

		if convert.ToInt(reply["ok"]) == 0 {
			// we have a command error.
			res.Error(convert.ToInt32(reply["code"]), convert.ToString(reply["errmsg"]))
		}

		res.Write(response)

	case messages.UpdateType:
		u, err := messages.ToUpdateRequest(req)
		if err != nil {
			Log(ERROR, "%#v\n", err)
			next(req, res)
			return
		}

		b := updateToBSONDoc(u)

		reply := bson.D{}
		mongoSession.DB(u.Database).Run(b, &reply)

		response := messages.UpdateResponse{
			N:         convert.ToInt32(bsonutil.FindValueByKey("n", reply), -1),
			NModified: convert.ToInt32(bsonutil.FindValueByKey("nModified", reply), -1),
		}

		writeErrors, err := convert.ConvertToBSONMapSlice(
			bsonutil.FindValueByKey("writeErrors", reply))
		if err == nil {
			// we have write errors
			response.WriteErrors = writeErrors
		}

		rawUpserted := bsonutil.FindValueByKey("upserted", reply)
		upserted, err := convert.ConvertToBSONDocSlice(rawUpserted)
		if err == nil {
			// we have upserts
			response.Upserted = upserted
		}

		if convert.ToInt(bsonutil.FindValueByKey("ok", reply)) == 0 {
			// we have a command error.
			res.Error(convert.ToInt32(bsonutil.FindValueByKey("code", reply)),
				convert.ToString(bsonutil.FindValueByKey("errmsg", reply)))
		}

		res.Write(response)

	case messages.DeleteType:
		d, err := messages.ToDeleteRequest(req)
		if err != nil {
			Log(ERROR, "%#v\n", err)
			next(req, res)
			return
		}

		b := deleteToBSONDoc(d)

		reply := bson.M{}
		mongoSession.DB(d.Database).Run(b, reply)

		response := messages.DeleteResponse{
			N: convert.ToInt32(reply["n"], -1),
		}
		writeErrors, err := convert.ConvertToBSONMapSlice(reply["writeErrors"])
		if err == nil {
			// we have write errors
			response.WriteErrors = writeErrors
		}

		if convert.ToInt(reply["ok"]) == 0 {
			// we have a command error.
			res.Error(convert.ToInt32(reply["code"]), convert.ToString(reply["errmsg"]))
		}

		Log(NOTICE, "Reply: %#v\n", reply)

		res.Write(response)

	case messages.GetMoreType:
		g, err := messages.ToGetMoreRequest(req)
		if err != nil {
			Log(ERROR, "%#v\n", err)
			next(req, res)
			return
		}

		Log(DEBUG, "%#v\n", g)

		// TODO: actually do something. Convert into an OP_GET_MORE, as mgo
		// abstracts it away in the Iter object
	default:
		Log(ERROR, "Unsupported operation")
	}
	// mongoSession.Close()
	next(req, res)

}

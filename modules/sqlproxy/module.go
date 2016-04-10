// Package sqlproxy contains a module that acts as a backend for SQLProxy,
package sqlproxy

import (
	"fmt"
	"github.com/10gen/sqlproxy"
	"github.com/10gen/sqlproxy/schema"
	"github.com/mongodbinc-interns/mongoproxy/bsonutil"
	"github.com/mongodbinc-interns/mongoproxy/convert"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/server"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// A SQLProxyModule takes the request, sends it to a SQLProxy instance, and then
// writes the response from SQLProxy into the ResponseWriter before calling
// the next module. It passes on requests unchanged.
type SQLProxyModule struct {
	evaluator *sqlproxy.Evaluator
	schema    *schema.Schema
	addr      string
}

func init() {
	server.Publish(&SQLProxyModule{})
}

func (_ *SQLProxyModule) New() server.Module {
	return &SQLProxyModule{}
}

func (_ *SQLProxyModule) Name() string {
	return "sqlproxy"
}

func (s *SQLProxyModule) Configure(conf bson.M) error {
	addr, ok := conf["address"].(string)
	if !ok {
		return fmt.Errorf("Invalid addresses: address is not a string")
	}
	s.addr = addr

	schemaStr, ok := conf["schema"].(string)
	if !ok {
		return fmt.Errorf("Invalid schema: schema is not a string")
	}

	cfg := &schema.Schema{}
	err := cfg.LoadFile(schemaStr)
	if err != nil {
		return fmt.Errorf("Error parsing schema file: %v", err)
	}

	opts := sqlproxy.Options{Addr: addr}

	evaluator, err := sqlproxy.NewEvaluator(cfg, opts)
	if err != nil {
		return fmt.Errorf("error creating mongoproxy evaluator: %v", err)
	}

	s.evaluator = evaluator

	s.schema = cfg

	return nil
}

type connCtx struct {
	db      string
	session *mgo.Session
}

func (c *connCtx) LastInsertId() int64 {
	return int64(0)
}

func (c *connCtx) RowCount() int64 {
	return int64(0)
}

func (c *connCtx) ConnectionId() uint32 {
	return uint32(0)
}

func (c *connCtx) DB() string {
	return c.db
}
func (c *connCtx) Session() *mgo.Session {
	return c.session
}

func (s *SQLProxyModule) Process(req messages.Requester, res messages.Responder,
	next server.PipelineFunc) {

	session := s.evaluator.Session()

	switch req.Type() {
	case messages.CommandType:

		command, err := messages.ToCommandRequest(req)
		if err != nil {
			Log(WARNING, "Error converting to command: %#v", err)
			res.Error(9001, fmt.Sprintf("error converting to command: %v", err))
			next(req, res)
			return
		}

		b := command.ToBSON()

		switch command.CommandName {
		case "sql":
			query, ok := command.Args["query"].(string)
			if !ok {
				Log(WARNING, "SQL query must have string query")
				res.Error(9002, "SQL query must have string query")
				next(req, res)
				return
			}

			conn := &connCtx{db: command.Database, session: session}

			headers, resultSet, err := s.evaluator.EvalSelect(command.Database, query, nil, conn)
			if err != nil {
				Log(WARNING, "error running SQL query: %v", err)
				res.Error(9003, fmt.Sprintf("error running SQL query: %v", err))
				next(req, res)
				return
			}

			var results []bson.M

			for _, fields := range resultSet {
				result := bson.M{}
				for j, field := range fields {
					result[headers[j]] = field
				}
				results = append(results, result)
			}

			reply := bson.M{
				"result": results,
				"ok":     1,
			}

			response := messages.CommandResponse{
				Reply: reply,
			}

			res.Write(response)

		default:

			reply := bson.M{}
			err = session.DB(command.Database).Run(b, reply)
			if err != nil {
				// log an error if we can
				qErr, ok := err.(*mgo.QueryError)
				Log(WARNING, "Error running command %v: %v", command.CommandName, err)
				if ok {
					res.Error(int32(qErr.Code), qErr.Message)
				} else {
					res.Error(-1, "Unknown error")
				}
				next(req, res)
				return
			}

			response := messages.CommandResponse{
				Reply: reply,
			}

			if convert.ToInt(reply["ok"]) == 0 {
				// we have a command error.
				res.Error(convert.ToInt32(reply["code"]), convert.ToString(reply["errmsg"]))
				next(req, res)
				return
			}

			res.Write(response)
		}
	case messages.InsertType:
		insert, err := messages.ToInsertRequest(req)
		if err != nil {
			Log(WARNING, "Error converting to Insert command: %#v", err)
			next(req, res)
			return
		}

		b := insert.ToBSON()

		reply := bson.M{}
		err = session.DB(insert.Database).Run(b, reply)
		if err != nil {
			// log an error if we can
			qErr, ok := err.(*mgo.QueryError)
			if ok {
				res.Error(int32(qErr.Code), qErr.Message)
			}
			next(req, res)
			return
		}

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
			next(req, res)
			return
		}

		res.Write(response)

	case messages.UpdateType:
		u, err := messages.ToUpdateRequest(req)
		if err != nil {
			Log(WARNING, "Error converting to Update command: %v", err)
			next(req, res)
			return
		}

		b := u.ToBSON()

		reply := bson.D{}
		err = session.DB(u.Database).Run(b, &reply)
		if err != nil {
			// log an error if we can
			qErr, ok := err.(*mgo.QueryError)
			if ok {
				res.Error(int32(qErr.Code), qErr.Message)
			}
			next(req, res)
			return
		}

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
			next(req, res)
			return
		}

		res.Write(response)

	case messages.DeleteType:
		d, err := messages.ToDeleteRequest(req)
		if err != nil {
			Log(WARNING, "Error converting to Delete command: %v", err)
			next(req, res)
			return
		}

		b := d.ToBSON()

		reply := bson.M{}
		err = session.DB(d.Database).Run(b, reply)
		if err != nil {
			// log an error if we can
			qErr, ok := err.(*mgo.QueryError)
			if ok {
				res.Error(int32(qErr.Code), qErr.Message)
			}
			next(req, res)
			return
		}

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
			next(req, res)
			return
		}

		Log(INFO, "Reply: %#v", reply)

		res.Write(response)

	case messages.FindType:
		f, err := messages.ToFindRequest(req)
		if err != nil {
			Log(WARNING, "Error converting to a Find command: %#v", err)
			next(req, res)
			return
		}

		c := session.DB(f.Database).C(f.Collection)
		query := c.Find(f.Filter).Batch(int(f.Limit)).Skip(int(f.Skip)).Prefetch(0)

		if f.Projection != nil {
			query = query.Select(f.Projection)
		}

		var iter = query.Iter()
		var results []bson.D

		cursorID := int64(0)

		if f.Limit > 0 {
			// only store the amount specified by the limit
			for i := 0; i < int(f.Limit); i++ {
				var result bson.D
				ok := iter.Next(&result)
				if !ok {
					err = iter.Err()
					if err != nil {
						Log(WARNING, "Error on Find Command: %#v", err)

						// log an error if we can
						qErr, ok := err.(*mgo.QueryError)
						if ok {
							res.Error(int32(qErr.Code), qErr.Message)
						}
						iter.Close()
						next(req, res)
						return
					}
					// we ran out of documents, but didn't have an error
					break
				}
				if cursorID == 0 {
					cursorID = iter.CursorID()
				}
				results = append(results, result)
			}
		} else {
			// dump all of them
			err = iter.All(&results)
			if err != nil {
				Log(WARNING, "Error on Find Command: %#v", err)

				// log an error if we can
				qErr, ok := err.(*mgo.QueryError)
				if ok {
					res.Error(int32(qErr.Code), qErr.Message)
				}
				next(req, res)
				return
			}
		}

		response := messages.FindResponse{
			Database:   f.Database,
			Collection: f.Collection,
			Documents:  results,
			CursorID:   cursorID,
		}

		res.Write(response)

	case messages.GetMoreType:
		g, err := messages.ToGetMoreRequest(req)
		if err != nil {
			Log(WARNING, "Error converting to GetMore command: %#v", err)
			next(req, res)
			return
		}
		Log(DEBUG, "%#v", g)

		// make an iterable to get more
		c := session.DB(g.Database).C(g.Collection)
		batch := make([]bson.Raw, 0)
		iter := c.NewIter(session, batch, g.CursorID, nil)
		iter.SetBatch(int(g.BatchSize))

		var results []bson.D
		cursorID := int64(0)

		for i := 0; i < int(g.BatchSize); i++ {
			var result bson.D
			ok := iter.Next(&result)
			if !ok {
				err = iter.Err()
				if err != nil {
					Log(WARNING, "Error on GetMore Command: %#v", err)

					if err == mgo.ErrCursor {
						// we return an empty getMore with an errored out
						// cursor
						response := messages.GetMoreResponse{
							CursorID:      cursorID,
							Database:      g.Database,
							Collection:    g.Collection,
							InvalidCursor: true,
						}
						res.Write(response)
						next(req, res)
						return
					}

					// log an error if we can
					qErr, ok := err.(*mgo.QueryError)
					if ok {
						res.Error(int32(qErr.Code), qErr.Message)
					}
					iter.Close()
					next(req, res)
					return
				}
				break
			}
			if cursorID == 0 {
				cursorID = iter.CursorID()
			}
			results = append(results, result)
		}

		response := messages.GetMoreResponse{
			CursorID:   cursorID,
			Database:   g.Database,
			Collection: g.Collection,
			Documents:  results,
		}

		res.Write(response)
	default:
		Log(WARNING, "Unsupported operation: %v", req.Type())
	}

	next(req, res)

}

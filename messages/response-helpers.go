package messages

import (
	"bytes"
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/buffer"
	"github.com/mongodbinc-interns/mongoproxy/convert"
	"gopkg.in/mgo.v2/bson"
)

// A ResponseWriter is the interface that is used to convert module responses
// to wire protocol messages.
type ResponseWriter interface {
	// ToBytes encodes a ResponseWriter into a valid wire protocol message
	// corresponding to the input response header.
	ToBytes(MsgHeader) ([]byte, error)

	// ToBSON encodes a ResponseWriter into a BSON document that can be examined
	// by other modules.
	ToBSON() bson.M
}

// A struct that represents a response to a generic command.
type CommandResponse struct {
	Metadata  bson.M
	Reply     bson.M
	Documents []bson.D
}

func (c CommandResponse) ToBytes(header MsgHeader) ([]byte, error) {
	resHeader := createResponseHeader(header)
	startingFrom := int32(0)

	flags := int32(8)

	buf := bytes.NewBuffer([]byte{})

	// write all documents
	err := buffer.WriteToBuf(buf, resHeader, int32(flags), int64(0), int32(startingFrom),
		int32(1+len(c.Documents)))
	if err != nil {
		return nil, fmt.Errorf("error writing prepared response %v\n", err)
	}
	reply := c.Reply
	reply["ok"] = 1
	docBytes, err := marshalReplyDocs(reply, c.Documents)
	if err != nil {
		return nil, fmt.Errorf("error marshaling documents")
	}

	resp := append(buf.Bytes(), docBytes...)

	resp = setMessageSize(resp)

	return resp, nil
}

func (c CommandResponse) ToBSON() bson.M {
	return c.Reply
}

// A struct that represents a response to a find command.
type FindResponse struct {
	CursorID   int64
	Database   string
	Collection string
	Documents  []bson.D

	// A QueryFailure is an object with an $err field to show that a query
	// has failed. An empty QueryFailure assumes a successful query.
	QueryFailure bson.M
}

func (f FindResponse) ToBytes(header MsgHeader) ([]byte, error) {
	resHeader := createResponseHeader(header)
	startingFrom := int32(0)

	flags := int32(8)

	buf := bytes.NewBuffer([]byte{})

	// override reply if we had a query failure.
	_, ok := f.QueryFailure["$err"]
	if ok {
		flags = int32(8)
		flags = convert.WriteBit32LE(flags, 1, true)

		err := buffer.WriteToBuf(buf, resHeader, int32(flags), f.CursorID, int32(startingFrom),
			int32(1))
		if err != nil {
			return nil, fmt.Errorf("error writing prepared response %v\n", err)
		}
		docBytes, err := marshalReplyDocs(f.QueryFailure, nil)
		if err != nil {
			return nil, fmt.Errorf("error marshaling documents")
		}
		resp := append(buf.Bytes(), docBytes...)

		resp = setMessageSize(resp)
		return resp, nil
	}

	// write all documents
	err := buffer.WriteToBuf(buf, resHeader, int32(flags), f.CursorID, int32(startingFrom),
		int32(len(f.Documents)))
	if err != nil {
		return nil, fmt.Errorf("error writing prepared response %v\n", err)
	}
	docBytes, err := marshalReplyDocs(nil, f.Documents)
	if err != nil {
		return nil, fmt.Errorf("error marshaling documents")
	}

	resp := append(buf.Bytes(), docBytes...)

	resp = setMessageSize(resp)

	return resp, nil
}

// ToBSON converts a FindResponse to a BSON format that is compatible with
// the command response spec
func (f FindResponse) ToBSON() bson.M {
	r := bson.M{}
	cursor := bson.M{
		"id":         f.CursorID,
		"ns":         f.Database + "." + f.Collection,
		"firstBatch": f.Documents,
	}

	r["cursor"] = cursor
	return r
}

// A struct that represents a response to a getMore command.
type GetMoreResponse struct {
	CursorID   int64
	Database   string
	Collection string
	Documents  []bson.D
}

func (g GetMoreResponse) ToBytes(header MsgHeader) ([]byte, error) {
	b := g.ToBSON()
	b["ok"] = 1
	return EncodeBSON(header, b)
}

func (g GetMoreResponse) ToBSON() bson.M {
	r := bson.M{}
	cursor := bson.M{
		"id":         g.CursorID,
		"ns":         g.Database + "." + g.Collection,
		"firstBatch": g.Documents,
	}

	r["cursor"] = cursor
	return r
}

// A struct that represents a response to an insert command.
type InsertResponse struct {

	// the number of documents inserted
	N int32

	// a list of write errors
	// TODO: create a WriteError struct
	WriteErrors []bson.M
}

func (i InsertResponse) ToBytes(header MsgHeader) ([]byte, error) {
	b := i.ToBSON()
	b["ok"] = 1
	return EncodeBSON(header, b)
}

func (i InsertResponse) ToBSON() bson.M {
	r := bson.M{
		"n": i.N,
	}
	if i.WriteErrors != nil && len(i.WriteErrors) > 0 {
		r["writeErrors"] = i.WriteErrors
	}

	return r
}

// A struct that represents a response to an update command.
type UpdateResponse struct {

	// the number of documents that matched the selector
	N int32

	// the number of documents actually updated
	NModified int32

	// the documents that didn't already exist and were upserted
	Upserted []bson.D

	// a list of write errors that occurred while updating
	WriteErrors []bson.M
}

func (u UpdateResponse) ToBytes(header MsgHeader) ([]byte, error) {
	b := u.ToBSON()
	b["ok"] = 1
	return EncodeBSON(header, b)
}

func (u UpdateResponse) ToBSON() bson.M {
	r := bson.M{
		"n":         u.N,
		"nModified": u.NModified,
	}
	if u.Upserted != nil && len(u.Upserted) > 0 {
		r["upserted"] = u.Upserted
	}
	if u.WriteErrors != nil && len(u.WriteErrors) > 0 {
		r["writeErrors"] = u.WriteErrors
	}

	return r
}

// A struct that represents a response to a delete command.
type DeleteResponse struct {
	// the number of documents deleted
	N int32

	// a list of write errors that occurred while deleting
	WriteErrors []bson.M
}

func (d DeleteResponse) ToBytes(header MsgHeader) ([]byte, error) {
	b := d.ToBSON()
	b["ok"] = 1
	return EncodeBSON(header, b)
}

func (d DeleteResponse) ToBSON() bson.M {
	r := bson.M{
		"n": d.N,
	}
	if d.WriteErrors != nil && len(d.WriteErrors) > 0 {
		r["writeErrors"] = d.WriteErrors
	}

	return r
}

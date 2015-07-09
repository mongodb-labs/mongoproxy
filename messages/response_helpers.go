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
		return nil, fmt.Errorf("error writing prepared response: %v", err)
	}
	reply := c.Reply
	reply["ok"] = 1
	docBytes, err := marshalReplyDocs(reply, c.Documents)
	if err != nil {
		return nil, fmt.Errorf("error marshaling documents: %v", err)
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
			return nil, fmt.Errorf("error writing prepared response: %v", err)
		}
		docBytes, err := marshalReplyDocs(f.QueryFailure, nil)
		if err != nil {
			return nil, fmt.Errorf("error marshaling documents: %v", err)
		}
		resp := append(buf.Bytes(), docBytes...)

		resp = setMessageSize(resp)
		return resp, nil
	}

	// write all documents
	err := buffer.WriteToBuf(buf, resHeader, int32(flags), f.CursorID, int32(startingFrom),
		int32(len(f.Documents)))
	if err != nil {
		return nil, fmt.Errorf("error writing prepared response: %v", err)
	}
	docBytes, err := marshalReplyDocs(nil, f.Documents)
	if err != nil {
		return nil, fmt.Errorf("error marshaling documents: %v", err)
	}

	resp := append(buf.Bytes(), docBytes...)

	resp = setMessageSize(resp)

	return resp, nil
}

// ToBSON converts a FindResponse to a BSON format that is compatible with
// the command response spec
func (f FindResponse) ToBSON() bson.M {
	return bson.M{
		"cursor": bson.M{
			"id":         f.CursorID,
			"ns":         f.Database + "." + f.Collection,
			"firstBatch": f.Documents,
		},
	}
}

// A struct that represents a response to a getMore command.
type GetMoreResponse struct {
	CursorID      int64
	Database      string
	Collection    string
	Documents     []bson.D
	InvalidCursor bool // true if the cursor wasn't valid at the server.
}

func (g GetMoreResponse) ToBytes(header MsgHeader) ([]byte, error) {
	resHeader := createResponseHeader(header)
	startingFrom := int32(0)

	flags := int32(8)

	if g.InvalidCursor {
		// invalid cursor. Return with no documents.
		flags = convert.WriteBit32LE(flags, 0, true)
		buf := bytes.NewBuffer([]byte{})

		// write all documents
		err := buffer.WriteToBuf(buf, resHeader, int32(flags), g.CursorID, int32(startingFrom),
			int32(0))
		if err != nil {
			return nil, fmt.Errorf("error writing prepared response: %v", err)
		}

		resp := setMessageSize(buf.Bytes())

		return resp, nil
	}

	buf := bytes.NewBuffer([]byte{})

	// write all documents
	err := buffer.WriteToBuf(buf, resHeader, int32(flags), g.CursorID, int32(startingFrom),
		int32(len(g.Documents)))
	if err != nil {
		return nil, fmt.Errorf("error writing prepared response: %v", err)
	}
	docBytes, err := marshalReplyDocs(nil, g.Documents)
	if err != nil {
		return nil, fmt.Errorf("error marshaling documents: %v", err)
	}

	resp := append(buf.Bytes(), docBytes...)

	resp = setMessageSize(resp)

	return resp, nil
}

func (g GetMoreResponse) ToBSON() bson.M {

	return bson.M{
		"cursor": bson.M{
			"id":        g.CursorID,
			"ns":        g.Database + "." + g.Collection,
			"nextBatch": g.Documents,
		},
	}
}

// A struct that represents a response to an insert command.
type InsertResponse struct {

	// the number of documents inserted. If the number is negative, the field
	// will not be exported when writing to the wire protocol.
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
	r := bson.M{}

	// some replies require the n field to not exist. The behavior thus is that if
	// n is negative, it won't be sent as BSON.
	if i.N >= 0 {
		r["n"] = i.N
	}
	if i.WriteErrors != nil && len(i.WriteErrors) > 0 {
		r["writeErrors"] = i.WriteErrors
	}

	return r
}

// A struct that represents a response to an update command.
type UpdateResponse struct {

	// the number of documents that matched the selector. If the number is negative,
	// N will not be exported when writing to the wire protocol.
	N int32

	// the number of documents actually updated. Not exported on a negative value.
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
	r := bson.M{}
	if u.N >= 0 {
		r["n"] = u.N
	}
	if u.NModified >= 0 {
		r["nModified"] = u.NModified
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
	// the number of documents deleted. Not exported on a negative value.
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
	r := bson.M{}
	if d.N >= 0 {
		r["n"] = d.N
	}
	if d.WriteErrors != nil && len(d.WriteErrors) > 0 {
		r["writeErrors"] = d.WriteErrors
	}

	return r
}

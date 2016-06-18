package messages

import (
	"bytes"
	"fmt"
	"github.com/mongodb-labs/mongoproxy/buffer"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
)

func shouldHaveSameContents(actual interface{}, expected ...interface{}) string {
	// quick hack since bson.M returns things out of order. This checks to make
	// sure two slices have the same contents

	actualMap := make(map[interface{}]int)
	expectedMap := make(map[interface{}]int)

	aType := reflect.TypeOf(actual)
	eType := reflect.TypeOf(expected[0])

	if aType != eType {
		return "They two values should be of the same type"
	}

	if len(actualMap) != len(expectedMap) {
		return "The two values should be the same length."
	}

	actualSlice, ok := actual.([]byte)
	if !ok {
		return "Only supports byte structs"
	}
	expectedSlice, ok := expected[0].([]byte)
	if !ok {
		return "Only supports byte structs"
	}

	for i := 0; i < len(actualMap); i++ {
		actualMap[actualSlice[i]] += 1
		expectedMap[expectedSlice[i]] += 1
	}

	for doc := range actualMap {
		if actualMap[doc] != expectedMap[doc] {
			return "Value mismatch."
		}
	}
	return ""
}

func createWireProtocolMessage(responseTo int32, flags int32, cursorID int64,
	startingFrom int32, docs []interface{}) ([]byte, error) {

	resHeader := MsgHeader{
		ResponseTo: responseTo, // requestID from the original request
		OpCode:     1,
	}

	numberReturned := len(docs)
	buf := bytes.NewBuffer([]byte{})
	err := buffer.WriteToBuf(buf, resHeader, flags, cursorID, startingFrom,
		int32(numberReturned))
	if err != nil {
		return nil, fmt.Errorf("error writing prepared response %v\n", err)
	}
	resp := buf.Bytes()

	for _, doc := range docs {
		docBytes, err := bson.Marshal(doc)
		if err != nil {
			return nil, fmt.Errorf("error marshaling response document: %v\n", err)
		}
		resp = append(resp, docBytes...)
	}

	resp = setMessageSize(resp)

	return resp, nil
}

func TestEncodeFindResponse(t *testing.T) {
	Convey("Encode a FindResponse struct to send over the wire protocol", t, func() {
		Convey("that has no errors", func() {
			r := FindResponse{}
			r.Database = "db"
			r.Collection = "foo"
			r.CursorID = int64(0)
			docs := make([]bson.D, 2)
			docs[0] = mockQuery
			docs[1] = mockCommand

			r.Documents = docs

			res := ModuleResponse{}
			res.Write(r)

			reqHeader := MsgHeader{
				RequestID: int32(5),
				OpCode:    int32(2004),
			}

			docsInt := make([]interface{}, 2)
			docsInt[0] = mockQuery
			docsInt[1] = mockCommand

			expected, err := createWireProtocolMessage(reqHeader.RequestID, int32(8), int64(0), int32(0), docsInt)
			So(err, ShouldBeNil)
			actual, err := Encode(reqHeader, res)
			So(err, ShouldBeNil)
			So(actual, shouldHaveSameContents, expected)
		})
		Convey("that has a query error", func() {
			r := FindResponse{}
			r.Database = "db"
			r.Collection = "foo"
			r.CursorID = int64(0)

			qErr := bson.M{}
			qErr["errmsg"] = "This is an error"
			qErr["code"] = 100

			qErrMessage := bson.M{}
			qErrMessage["$err"] = qErr

			r.QueryFailure = qErrMessage

			docs := make([]bson.D, 2)
			docs[0] = mockQuery
			docs[1] = mockCommand

			r.Documents = docs

			res := ModuleResponse{}
			res.Write(r)

			reqHeader := MsgHeader{
				RequestID: int32(5),
				OpCode:    int32(2004),
			}

			docsInt := make([]interface{}, 1)
			qErrMessage["ok"] = 0
			docsInt[0] = qErrMessage

			expected, err := createWireProtocolMessage(reqHeader.RequestID, int32(10), int64(0), int32(0), docsInt)
			So(err, ShouldBeNil)
			actual, err := Encode(reqHeader, res)
			So(err, ShouldBeNil)
			So(actual, shouldHaveSameContents, expected)
		})

		Convey("that has a command error", func() {
			r := FindResponse{}
			r.Database = "db"
			r.Collection = "foo"
			r.CursorID = int64(0)
			docs := make([]bson.D, 2)
			docs[0] = mockQuery
			docs[1] = mockCommand

			r.Documents = docs

			res := ModuleResponse{}
			res.Write(r)
			res.Error(0, "this is an error")

			reqHeader := MsgHeader{
				RequestID: int32(5),
				OpCode:    int32(2004),
			}

			docsInt := make([]interface{}, 1)
			commandErr := bson.M{}
			commandErr["ok"] = 0
			commandErr["errmsg"] = "this is an error"
			commandErr["code"] = 0
			docsInt[0] = commandErr

			expected, err := createWireProtocolMessage(reqHeader.RequestID, int32(8), int64(0), int32(0), docsInt)
			So(err, ShouldBeNil)
			actual, err := Encode(reqHeader, res)
			So(err, ShouldBeNil)
			So(actual, shouldHaveSameContents, expected)
		})
	})
}

func TestEncodeCommandResponse(t *testing.T) {
	Convey("Encode a default CommandResponse to send over the wire protocol", t, func() {
		r := bson.M{}
		r["foo"] = "bar"

		reqHeader := MsgHeader{
			RequestID: int32(5),
			OpCode:    int32(2004),
		}

		res := ModuleResponse{}

		reply := CommandResponse{}
		reply.Reply = r
		res.Write(reply)

		r["ok"] = 1
		resSlice := make([]interface{}, 1)
		resSlice[0] = r

		expected, err := createWireProtocolMessage(reqHeader.RequestID, int32(8), int64(0), int32(0), resSlice)
		So(err, ShouldBeNil)
		actual, err := Encode(reqHeader, res)
		So(err, ShouldBeNil)
		So(actual, shouldHaveSameContents, expected)
	})
}

func TestEncodeInsertResponse(t *testing.T) {
	Convey("Encode an InsertResponse to send over the wire protocol", t, func() {

		reqHeader := MsgHeader{
			RequestID: int32(5),
			OpCode:    int32(2004),
		}

		r := InsertResponse{}
		r.N = 5

		res := ModuleResponse{}
		res.Write(r)

		resSlice := make([]interface{}, 1)
		resSlice[0] = bson.M{"n": 5, "ok": 1}

		expected, err := createWireProtocolMessage(reqHeader.RequestID, int32(8), int64(0), int32(0), resSlice)
		So(err, ShouldBeNil)
		actual, err := Encode(reqHeader, res)
		So(err, ShouldBeNil)
		So(actual, shouldHaveSameContents, expected)
	})
}

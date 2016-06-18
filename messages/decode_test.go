package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mongodb-labs/mongoproxy/buffer"
	. "github.com/mongodb-labs/mongoproxy/log"
	"github.com/mongodb-labs/mongoproxy/mock"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

var mockQuery = bson.D{{"hello", 1}}
var mockCommand = bson.D{{"isMaster", 1}}

// creates a valid OP_QUERY find
func createMockQuery(id int32, flags int32, namespace string,
	skip int32, limit int32, query interface{}) []byte {
	responseTo := int32(0)
	opCode := int32(2004)
	queryBytes, err := bson.Marshal(query)

	if err != nil {
		fmt.Println("Error encoding BSON")
	}

	fullCollectionBytes := []byte(namespace)

	fullCollectionBytes = append(fullCollectionBytes, byte('\x00'))
	fmt.Printf("%v\n", fullCollectionBytes)

	buf := new(bytes.Buffer)

	buffer.WriteToBuf(buf, int32(0), id, responseTo, opCode, flags, fullCollectionBytes,
		skip, limit, queryBytes)

	input := buf.Bytes()
	respSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(respSize, uint32(len(input)))
	input[0] = respSize[0]
	input[1] = respSize[1]
	input[2] = respSize[2]
	input[3] = respSize[3]

	return input
}

func createMockQueryNoNullTerm(id int32, flags int32, namespace string,
	skip int32, limit int32, query interface{}) []byte {
	responseTo := int32(0)
	opCode := int32(2004)
	queryBytes, err := bson.Marshal(query)

	q2 := bson.M{}
	bson.Unmarshal(queryBytes, q2)
	fmt.Printf("%v\n", queryBytes)
	fmt.Printf("%v\n", q2)
	if err != nil {
		fmt.Println("Error encoding BSON")
	}

	fullCollectionBytes := []byte(namespace)

	fmt.Printf("%v\n", fullCollectionBytes)

	buf := new(bytes.Buffer)

	buffer.WriteToBuf(buf, int32(0), id, responseTo, opCode, flags, fullCollectionBytes,
		skip, limit, queryBytes)

	input := buf.Bytes()
	respSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(respSize, uint32(len(input)))
	input[0] = respSize[0]
	input[1] = respSize[1]
	input[2] = respSize[2]
	input[3] = respSize[3]

	return input
}

func createMockInsert(id int32, flags int32, namespace string, docs []interface{}) []byte {
	responseTo := int32(0)
	opCode := int32(2002)

	fullCollectionBytes := []byte(namespace)

	fullCollectionBytes = append(fullCollectionBytes, byte('\x00'))

	queryBytes := make([]byte, 0)
	for i := 0; i < len(docs); i++ {
		q, err := bson.Marshal(docs[i])
		if err != nil {
			fmt.Println("Error encoding BSON")
		}
		queryBytes = append(queryBytes, q...)
	}

	buf := new(bytes.Buffer)

	buffer.WriteToBuf(buf, int32(0), id, responseTo, opCode, flags, fullCollectionBytes,
		queryBytes)

	input := buf.Bytes()
	respSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(respSize, uint32(len(input)))
	input[0] = respSize[0]
	input[1] = respSize[1]
	input[2] = respSize[2]
	input[3] = respSize[3]

	return input
}

func createMockUpdate(id int32, flags int32, namespace string, selector interface{}, update interface{}) []byte {
	responseTo := int32(0)
	opCode := int32(2001)

	fullCollectionBytes := []byte(namespace)

	fullCollectionBytes = append(fullCollectionBytes, byte('\x00'))

	queryBytes, err := bson.Marshal(selector)

	if err != nil {
		fmt.Println("Error encoding BSON")
	}

	updateBytes, err := bson.Marshal(update)

	if err != nil {
		fmt.Println("Error encoding BSON")
	}

	buf := new(bytes.Buffer)

	buffer.WriteToBuf(buf, int32(0), id, responseTo, opCode, int32(0), fullCollectionBytes,
		flags, queryBytes, updateBytes)

	input := buf.Bytes()
	respSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(respSize, uint32(len(input)))
	input[0] = respSize[0]
	input[1] = respSize[1]
	input[2] = respSize[2]
	input[3] = respSize[3]

	return input
}

func createMockDelete(id int32, flags int32, namespace string, selector interface{}) []byte {
	responseTo := int32(0)
	opCode := int32(2006)

	fullCollectionBytes := []byte(namespace)

	fullCollectionBytes = append(fullCollectionBytes, byte('\x00'))

	queryBytes, err := bson.Marshal(selector)

	if err != nil {
		fmt.Println("Error encoding BSON")
	}

	buf := new(bytes.Buffer)

	buffer.WriteToBuf(buf, int32(0), id, responseTo, opCode, int32(0), fullCollectionBytes,
		flags, queryBytes)

	input := buf.Bytes()
	respSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(respSize, uint32(len(input)))
	input[0] = respSize[0]
	input[1] = respSize[1]
	input[2] = respSize[2]
	input[3] = respSize[3]

	return input
}

func createMockGetMore(id int32, namespace string, numToReturn int32, cursorID int64) []byte {
	responseTo := int32(0)
	opCode := int32(2005)

	fullCollectionBytes := []byte(namespace)

	fullCollectionBytes = append(fullCollectionBytes, byte('\x00'))

	buf := new(bytes.Buffer)

	buffer.WriteToBuf(buf, int32(0), id, responseTo, opCode, int32(0), fullCollectionBytes,
		numToReturn, cursorID)

	input := buf.Bytes()
	respSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(respSize, uint32(len(input)))
	input[0] = respSize[0]
	input[1] = respSize[1]
	input[2] = respSize[2]
	input[3] = respSize[3]

	return input
}

func TestProcessHeader(t *testing.T) {
	Convey("Decode a header", t, func() {
		Convey("which reads 0 bytes", func() {
			m := mock.MockIO{
				Input:  make([]byte, 0),
				Output: make([]byte, 0)}
			m.Reset()

			_, err := processHeader(&m)

			So(err, ShouldNotBeNil)
		})

		Convey("handle a read error from the reader", func() {
			m := mock.MockIO{
				Input:  make([]byte, 12),
				Output: make([]byte, 0)}
			m.Reset()

			_, err := processHeader(&m)

			So(err, ShouldNotBeNil)

			_, err = processHeader(&m)

			So(err, ShouldNotBeNil)
		})
	})
}

func TestCreateFind(t *testing.T) {
	Convey("Process invalid find", t, func() {
		_, err := createFind(MsgHeader{}, "test", bson.M{})
		So(err, ShouldNotBeNil)
	})
}

func TestDecodeOpQuery(t *testing.T) {
	SetLogLevel(DEBUG)
	Convey("Decode a wire protocol OP_QUERY message", t, func() {
		Convey("that is a valid find command", func() {
			// create the mock connection

			Convey("with all defaults", func() {
				input := createMockQuery(int32(0), int32(0), "db.foo", int32(0), int32(0), mockQuery)
				m := mock.MockIO{
					Input:  input,
					Output: make([]byte, 0)}
				m.Reset()

				request, _, err := Decode(&m)
				So(err, ShouldBeNil)

				t := request.Type()
				So(t, ShouldEqual, "find")

				opq, err := ToFindRequest(request)
				So(err, ShouldBeNil)
				So(opq, ShouldNotBeNil)

			})

			Convey("that uses some of the flags", func() {
				input := createMockQuery(int32(0), int32(2+8+16), "db.foo", int32(0), int32(0), mockQuery)
				m := mock.MockIO{
					Input:  input,
					Output: make([]byte, 0)}
				m.Reset()

				request, _, err := Decode(&m)
				So(err, ShouldBeNil)

				t := request.Type()
				So(t, ShouldEqual, "find")

				opq, err := ToFindRequest(request)
				So(err, ShouldBeNil)
				So(opq, ShouldNotBeNil)

				So(opq.Tailable, ShouldEqual, true)
				So(opq.OplogReplay, ShouldEqual, true)
				So(opq.NoCursorTimeout, ShouldEqual, true)

			})
		})
		Convey("that is an invalid find command", func() {
			Convey("because it has no length", func() {
				input := createMockQuery(int32(0), int32(0), "db.foo", int32(0), int32(0), mockQuery)

				// mess up the length
				respSize := make([]byte, 4)
				binary.LittleEndian.PutUint32(respSize, uint32(0))
				input[0] = respSize[0]
				input[1] = respSize[1]
				input[2] = respSize[2]
				input[3] = respSize[3]

				m := mock.MockIO{
					Input:  input,
					Output: make([]byte, 0)}
				m.Reset()

				Log(ERROR, "%#v", m)

				request, _, err := Decode(&m)
				So(err, ShouldNotBeNil)
				So(request, ShouldBeNil)
			})
			Convey("because the namespace has no collection", func() {
				input := createMockQuery(int32(0), int32(0), "db.", int32(0), int32(0), mockQuery)
				m := mock.MockIO{
					Input:  input,
					Output: make([]byte, 0)}
				m.Reset()

				request, _, err := Decode(&m)
				fmt.Printf("%v\n", request)
				So(err, ShouldNotBeNil)
			})

			Convey("because the namespace doesn't have a null terminator", func() {
				// create the mock connection

				input := createMockQueryNoNullTerm(int32(0), int32(0), "db.foo", int32(0), int32(0), mockQuery)
				m := mock.MockIO{
					Input:  input,
					Output: make([]byte, 0)}
				m.Reset()

				_, _, err := Decode(&m)
				So(err, ShouldNotBeNil)
			})

			Convey("because the namespace has no database", func() {
				input := createMockQuery(int32(0), int32(0), ".foo", int32(0), int32(0), mockQuery)
				m := mock.MockIO{
					Input:  input,
					Output: make([]byte, 0)}
				m.Reset()

				request, _, err := Decode(&m)
				fmt.Printf("%v\n", request)
				So(err, ShouldNotBeNil)
			})

			Convey("because the namespace is just weird", func() {
				input := createMockQuery(int32(0), int32(0), "hello-there~", int32(0), int32(0), mockQuery)
				m := mock.MockIO{
					Input:  input,
					Output: make([]byte, 0)}
				m.Reset()

				request, _, err := Decode(&m)
				fmt.Printf("%v\n", request)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("that is a valid non-specialized command", func() {
			input := createMockQuery(int32(0), int32(0), "test.$cmd", int32(0), int32(0), mockCommand)
			m := mock.MockIO{
				Input:  input,
				Output: make([]byte, 0)}
			m.Reset()

			request, _, err := Decode(&m)
			So(err, ShouldBeNil)

			t := request.Type()
			So(t, ShouldEqual, "command")

			command, err := ToCommandRequest(request)
			So(err, ShouldBeNil)

			name := command.GetArg("isMaster")
			So(name, ShouldEqual, 1)
		})

		Convey("that is a valid insert command", func() {
			docs := make([]bson.D, 2)
			docs[0] = mockCommand
			docs[1] = mockQuery

			mockInsert := bson.D{{"insert", "foo"},
				{"documents", docs}}

			Log(DEBUG, "%#v\n", mockInsert)

			input := createMockQuery(int32(0), int32(0), "db.$cmd", int32(0), int32(0), mockInsert)
			m := mock.MockIO{
				Input:  input,
				Output: make([]byte, 0)}
			m.Reset()

			request, _, err := Decode(&m)
			So(err, ShouldBeNil)

			t := request.Type()
			So(t, ShouldEqual, "insert")

			opq, err := ToInsertRequest(request)
			So(err, ShouldBeNil)
			So(opq, ShouldNotBeNil)

			So(opq.Collection, ShouldEqual, "foo")
			So(opq.Database, ShouldEqual, "db")
			So(opq.Documents, ShouldResemble, docs)
			So(opq.Ordered, ShouldEqual, true)
		})

		Convey("that is a valid update command", func() {
			updates := make([]bson.M, 2)
			updates[0] = bson.M{"q": mockQuery, "u": mockCommand, "upsert": true}
			updates[1] = bson.M{"q": mockQuery, "u": mockCommand, "multi": true}
			mockUpdate := bson.D{{"update", "foo"},
				{"updates", updates},
				{"ordered", false}}

			input := createMockQuery(int32(0), int32(0), "db.$cmd", int32(0), int32(0), mockUpdate)
			m := mock.MockIO{
				Input:  input,
				Output: make([]byte, 0)}
			m.Reset()

			request, _, err := Decode(&m)
			So(err, ShouldBeNil)

			t := request.Type()
			So(t, ShouldEqual, "update")

			opq, err := ToUpdateRequest(request)
			So(err, ShouldBeNil)
			So(opq, ShouldNotBeNil)

			So(opq.Collection, ShouldEqual, "foo")
			So(opq.Database, ShouldEqual, "db")
			So(opq.Ordered, ShouldEqual, false)
			So(len(opq.Updates), ShouldEqual, 2)
			So(opq.Updates[1].Multi, ShouldEqual, true)
		})
	})

}

func TestDecodeOpInsert(t *testing.T) {
	Convey("Decode a wire protocol OP_INSERT message", t, func() {
		Convey("that is a valid insert command", func() {
			input := createMockInsert(int32(0), int32(0), "db.foo", []interface{}{mockQuery})
			m := mock.MockIO{
				Input:  input,
				Output: make([]byte, 0)}
			m.Reset()

			request, _, err := Decode(&m)
			So(err, ShouldBeNil)

			t := request.Type()
			So(t, ShouldEqual, "insert")

			opq, err := ToInsertRequest(request)
			So(err, ShouldBeNil)
			So(opq, ShouldNotBeNil)

			So(opq.Collection, ShouldEqual, "foo")
			So(opq.Database, ShouldEqual, "db")
			So(opq.Documents, ShouldResemble, []bson.D{mockQuery})
			So(opq.Ordered, ShouldEqual, true)
		})
	})

}

func TestDecodeOpUpdate(t *testing.T) {
	Convey("Decode a wire protocol OP_UPDATE message", t, func() {
		Convey("that is a valid update command", func() {
			input := createMockUpdate(int32(0), int32(0), "db.foo", mockQuery, mockCommand)
			m := mock.MockIO{
				Input:  input,
				Output: make([]byte, 0)}
			m.Reset()

			request, _, err := Decode(&m)
			So(err, ShouldBeNil)

			t := request.Type()
			So(t, ShouldEqual, "update")

			opq, err := ToUpdateRequest(request)
			So(err, ShouldBeNil)
			So(opq, ShouldNotBeNil)

			So(opq.Updates[0], ShouldNotBeNil)
			So(len(opq.Updates), ShouldEqual, 1)
			So(opq.Updates[0].Upsert, ShouldEqual, false)
		})
		Convey("that is a valid update command with some flags", func() {
			input := createMockUpdate(int32(0), int32(1), "db.foo", mockQuery, mockCommand)
			m := mock.MockIO{
				Input:  input,
				Output: make([]byte, 0)}
			m.Reset()

			request, _, err := Decode(&m)
			So(err, ShouldBeNil)

			t := request.Type()
			So(t, ShouldEqual, "update")

			opq, err := ToUpdateRequest(request)
			So(err, ShouldBeNil)
			So(opq, ShouldNotBeNil)

			So(opq.Updates[0], ShouldNotBeNil)
			So(len(opq.Updates), ShouldEqual, 1)
			So(opq.Updates[0].Upsert, ShouldEqual, true)
		})
	})

}

func TestDecodeOpDelete(t *testing.T) {
	Convey("Decode a wire protocol OP_DELETE message", t, func() {
		Convey("that is a valid delete command", func() {
			input := createMockDelete(int32(0), int32(0), "db.foo", mockQuery)
			m := mock.MockIO{
				Input:  input,
				Output: make([]byte, 0)}
			m.Reset()

			request, _, err := Decode(&m)
			So(err, ShouldBeNil)

			t := request.Type()
			So(t, ShouldEqual, "delete")

			opq, err := ToDeleteRequest(request)
			So(err, ShouldBeNil)
			So(opq, ShouldNotBeNil)

			So(opq.Database, ShouldEqual, "db")
			So(opq.Collection, ShouldEqual, "foo")
			So(len(opq.Deletes), ShouldEqual, 1)
			So(opq.Deletes[0].Selector, ShouldResemble, mockQuery)
			So(opq.Deletes[0].Limit, ShouldEqual, 0)
		})
	})

}

func TestDecodeOpGetMore(t *testing.T) {
	Convey("Decode a wire protocol OP_GET_MORE message", t, func() {
		Convey("that is a valid delete command", func() {
			input := createMockGetMore(int32(0), "db.foo", int32(20), int64(125))
			m := mock.MockIO{
				Input:  input,
				Output: make([]byte, 0)}
			m.Reset()

			request, _, err := Decode(&m)
			So(err, ShouldBeNil)

			t := request.Type()
			So(t, ShouldEqual, "getMore")

			opq, err := ToGetMoreRequest(request)
			So(err, ShouldBeNil)
			So(opq, ShouldNotBeNil)

			So(opq.Database, ShouldEqual, "db")
			So(opq.Collection, ShouldEqual, "foo")
			So(opq.CursorID, ShouldEqual, int64(125))
			So(opq.BatchSize, ShouldEqual, 20)
		})
	})
}

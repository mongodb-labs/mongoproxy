package buffer

import (
	"encoding/binary"
	"github.com/mongodb-labs/mongoproxy/mock"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestReadDocument(t *testing.T) {
	Convey("Test reading in a BSON document", t, func() {
		Convey("that should succeed on a valid document", func() {
			doc := bson.D{{"ok", 1}, {"two", 2}}
			bytes, err := bson.Marshal(doc)
			So(err, ShouldBeNil)

			m := mock.MockIO{
				Input:  bytes,
				Output: make([]byte, 0)}
			m.Reset()

			docSize, d, err := ReadDocument(&m)
			So(docSize, ShouldEqual, len(bytes))
			So(d, ShouldResemble, doc)
			So(err, ShouldBeNil)
		})
		Convey("that should succeed on a valid empty document", func() {
			bytes := make([]byte, 4)
			doc := bson.D{}
			binary.LittleEndian.PutUint32(bytes, 4)

			m := mock.MockIO{
				Input:  bytes,
				Output: make([]byte, 0)}
			m.Reset()

			docSize, d, err := ReadDocument(&m)
			So(docSize, ShouldEqual, 4)
			So(d, ShouldResemble, doc)
			So(err, ShouldBeNil)
		})

		Convey("that should fail if the docSize is too small", func() {
			bytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bytes, 1)

			m := mock.MockIO{
				Input:  bytes,
				Output: make([]byte, 0)}
			m.Reset()

			docSize, d, err := ReadDocument(&m)
			So(docSize, ShouldEqual, 0)
			So(d, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestRead32BitLE(t *testing.T) {
	Convey("Test a bunch of values", t, func() {
		values := []uint32{31415926, 0, 1, 10, 12, 15, 20, 99}
		for i := 0; i < len(values); i++ {
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, values[i])

			m := mock.MockIO{
				Input:  bs,
				Output: make([]byte, 0)}
			m.Reset()

			val, err := ReadInt32LE(&m)
			So(val, ShouldEqual, values[i])
			So(err, ShouldBeNil)
		}
	})
}

func TestRead64BitLE(t *testing.T) {
	Convey("Test a bunch of values", t, func() {
		values := []uint64{31415926, 0, 1, 10, 12, 15, 20, 99}
		for i := 0; i < len(values); i++ {
			bs := make([]byte, 8)
			binary.LittleEndian.PutUint64(bs, values[i])

			m := mock.MockIO{
				Input:  bs,
				Output: make([]byte, 0)}
			m.Reset()

			val, err := ReadInt64LE(&m)
			So(val, ShouldEqual, values[i])
			So(err, ShouldBeNil)
		}
	})
}

func TestReadNullTerminatedString(t *testing.T) {
	Convey("Test reading a null terminated string", t, func() {
		Convey("that should succeed if valid", func() {
			s := "This is a string"
			sNull := []byte(s)
			sNull = append(sNull, byte('\x00'))

			m := mock.MockIO{
				Input:  sNull,
				Output: make([]byte, 0)}
			m.Reset()

			n, str, err := ReadNullTerminatedString(&m, 999)
			So(n, ShouldEqual, len(sNull))
			So(str, ShouldEqual, s)
			So(err, ShouldBeNil)

			m.Reset()

			n, str, err = ReadNullTerminatedString(&m, int32(len(sNull)))
			So(n, ShouldEqual, len(sNull))
			So(str, ShouldEqual, s)
			So(err, ShouldBeNil)
		})
		Convey("should return an error if the length of the string is too long", func() {
			s := "This is a string"
			sNull := []byte(s)
			sNull = append(sNull, byte('\x00'))

			m := mock.MockIO{
				Input:  sNull,
				Output: make([]byte, 0)}
			m.Reset()

			n, str, err := ReadNullTerminatedString(&m, 1)
			So(n, ShouldEqual, 0)
			So(str, ShouldEqual, "")
			So(err, ShouldNotBeNil)

			m.Reset()

			n, str, err = ReadNullTerminatedString(&m, int32(len(sNull))-1)
			So(n, ShouldEqual, 0)
			So(str, ShouldEqual, "")
			So(err, ShouldNotBeNil)

		})
		Convey("should return an error if the string is not null-terminated", func() {
			s := "This is a string"
			sNull := []byte(s)

			m := mock.MockIO{
				Input:  sNull,
				Output: make([]byte, 0)}
			m.Reset()

			n, str, err := ReadNullTerminatedString(&m, 999)
			So(n, ShouldEqual, 0)
			So(str, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})
	})
}

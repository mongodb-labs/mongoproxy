package bsonutil

import (
	. "github.com/mongodb-labs/mongoproxy/log"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestFindByValueD(t *testing.T) {
	SetLogLevel(DEBUG)
	d := bson.D{{"field", "ok"}, {"field2", "ok"}, {"field3", "ok"}, {"field4", "ok"}}

	Convey("Find the first field in a doc", t, func() {
		raw := FindValueByKey("field", d)
		val, ok := raw.(string)
		So(ok, ShouldEqual, true)
		So(val, ShouldEqual, "ok")
	})

	Convey("Find the last field in a doc", t, func() {
		raw := FindValueByKey("field4", d)
		val, ok := raw.(string)
		So(ok, ShouldEqual, true)
		So(val, ShouldEqual, "ok")
	})

	Convey("Find a non-existent field", t, func() {
		raw := FindValueByKey("empty", d)
		_, ok := raw.(string)
		So(ok, ShouldEqual, false)
		So(raw, ShouldEqual, nil)
	})
}

func TestFindDeepValue(t *testing.T) {
	SetLogLevel(DEBUG)
	r := bson.M{
		"field": "ok",
		"multi": bson.M{
			"level": "ok",
		},
		"multiD": bson.D{{"level", "ok"}},
	}

	Convey("Find a top-level field in a map", t, func() {

		Convey("that exists", func() {
			raw := FindDeepValueInMap("field", r)
			val, ok := raw.(string)
			So(ok, ShouldEqual, true)
			So(val, ShouldEqual, "ok")
			So(raw, ShouldEqual, r["field"])
		})

		Convey("that doesn't exist", func() {
			raw := FindDeepValueInMap("empty", r)
			_, ok := raw.(string)
			So(ok, ShouldEqual, false)
			So(raw, ShouldEqual, nil)
		})
	})

	Convey("Find a deep field in a map", t, func() {
		Convey("that exists", func() {
			raw := FindDeepValueInMap("multi.level", r)
			val, ok := raw.(string)
			So(ok, ShouldEqual, true)
			So(val, ShouldEqual, "ok")
		})

		Convey("that converts a bson.D correctly", func() {
			raw := FindDeepValueInMap("multiD.level", r)
			val, ok := raw.(string)
			So(ok, ShouldEqual, true)
			So(val, ShouldEqual, "ok")
		})

		Convey("that doesn't exist", func() {
			raw := FindDeepValueInMap("multi.empty", r)
			_, ok := raw.(string)
			So(ok, ShouldEqual, false)
			So(raw, ShouldEqual, nil)
		})
	})
}

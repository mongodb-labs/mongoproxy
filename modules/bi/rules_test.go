package bi

import (
	"github.com/mongodbinc-interns/mongoproxy/convert"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestCreateSelector(t *testing.T) {
	SetLogLevel(DEBUG)

	Convey("Create a selector", t, func() {
		Convey("for a monthly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			expectedStart := time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
			valueField := "field"

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField}}

			actual, err := createSelector(t1, Monthly, valueField)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a daily metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			expectedStart := time.Date(2015, time.March, 1, 0, 0, 0, 0, time.UTC)
			valueField := "field"

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField}}

			actual, err := createSelector(t1, Daily, valueField)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for an hourly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			expectedStart := time.Date(2015, time.March, 20, 0, 0, 0, 0, time.UTC)
			valueField := "field"

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField}}

			actual, err := createSelector(t1, Hourly, valueField)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a minutely metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			expectedStart := time.Date(2015, time.March, 20, 14, 0, 0, 0, time.UTC)
			valueField := "field"

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField}}

			actual, err := createSelector(t1, Minutely, valueField)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a secondly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			expectedStart := time.Date(2015, time.March, 20, 14, 35, 0, 0, time.UTC)
			valueField := "field"

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField}}

			actual, err := createSelector(t1, Secondly, valueField)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})
	})

	Convey("Fail to create a selector", t, func() {
		Convey("for non-existent time granularities", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			valueField := "field"

			_, err := createSelector(t1, "err", valueField)

			So(err, ShouldNotBeNil)
		})
	})
}

func TestCreateUpdate(t *testing.T) {
	SetLogLevel(DEBUG)

	// NOTE: all of these fail for some reason, even though the actual / expected are identical
	// in every way. Goconvey complains of type mismatch between bson.D and bson.D (?!)
	Convey("Create an update", t, func() {
		Convey("for a monthly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)
			expected := bson.D{{"$inc", bson.D{{"total", 20}}}, {"$inc", bson.D{{"month.3", 20}}}}

			actual, err := createUpdate(t1, Monthly, valueFloat)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a daily metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)
			expected := bson.D{{"$inc", bson.D{{"total", 20}}}, {"$inc", bson.D{{"day.20", 20}}}}

			actual, err := createUpdate(t1, Daily, valueFloat)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a hourly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)
			expected := bson.D{{"$inc", bson.D{{"total", 20}}}, {"$inc", bson.D{{"hour.14", 20}}}}

			actual, err := createUpdate(t1, Hourly, valueFloat)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a minutely metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)
			expected := bson.D{{"$inc", bson.D{{"total", 20}}}, {"$inc", bson.D{{"minute.35", 20}}}}

			actual, err := createUpdate(t1, Minutely, valueFloat)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a secondly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)
			expected := bson.D{{"$inc", bson.D{{"total", 20}}}, {"$inc", bson.D{{"second.2", 20}}}}

			actual, err := createUpdate(t1, Secondly, valueFloat)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})
	})

	Convey("Fail to create an update", t, func() {
		Convey("for a non-existent time granularity", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)

			_, err := createUpdate(t1, "err", valueFloat)

			So(err, ShouldNotBeNil)
		})
	})
}

func TestCreateSingleUpdate(t *testing.T) {
	Convey("Create a single update request object", t, func() {
		inputDoc := bson.D{{"price", 5}}
		t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
		granularity := Daily
		rule := Rule{
			ValueField: "price",
		}

		expectedSelector, err := createSelector(t1, granularity, rule.ValueField)
		So(err, ShouldBeNil)
		expectedUpdate, err := createUpdate(t1, granularity, convert.ToFloat64(5))
		So(err, ShouldBeNil)

		expected := &messages.SingleUpdate{
			Selector: expectedSelector,
			Update:   expectedUpdate,
			Upsert:   true,
			Multi:    false,
		}

		actual, err := createSingleUpdate(inputDoc, t1, granularity, rule)

		So(err, ShouldBeNil)
		So(actual, ShouldResemble, expected)
	})
	Convey("Fail to create a single update request object", t, func() {
		Convey("for a non-existent time granularity", func() {
			inputDoc := bson.D{{"price", 5}}
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			granularity := "err"
			rule := Rule{
				ValueField: "price",
			}

			_, err := createSingleUpdate(inputDoc, t1, granularity, rule)

			So(err, ShouldNotBeNil)
		})
		Convey("for a doc without a value field", func() {
			inputDoc := bson.D{{"hello", "hi"}}
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			granularity := "err"
			rule := Rule{
				ValueField: "price",
			}

			_, err := createSingleUpdate(inputDoc, t1, granularity, rule)

			So(err, ShouldNotBeNil)

		})
		Convey("for a doc whose value field is not a number", func() {
			inputDoc := bson.D{{"price", "hi"}}
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			granularity := "err"
			rule := Rule{
				ValueField: "price",
			}

			_, err := createSingleUpdate(inputDoc, t1, granularity, rule)

			So(err, ShouldNotBeNil)
		})
	})
}

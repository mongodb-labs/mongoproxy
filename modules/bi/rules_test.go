package bi

import (
	"github.com/mongodb-labs/mongoproxy/convert"
	. "github.com/mongodb-labs/mongoproxy/log"
	"github.com/mongodb-labs/mongoproxy/messages"
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

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField},
				{"value", bson.D{{"$exists", false}}}}
			expectedString := bson.D{{"start", expectedStart}, {"valueField", valueField},
				{"value", "value"}, {"single", true}}

			actual, err := createSelector(t1, Monthly, valueField)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)

			// string value
			actualString, err := createSelectorString(t1, Monthly, valueField, "value")
			So(actualString, ShouldResemble, expectedString)
			So(err, ShouldBeNil)
		})

		Convey("for a daily metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			expectedStart := time.Date(2015, time.March, 1, 0, 0, 0, 0, time.UTC)
			valueField := "field"

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField},
				{"value", bson.D{{"$exists", false}}}}
			expectedString := bson.D{{"start", expectedStart}, {"valueField", valueField},
				{"value", "value"}, {"single", true}}

			actual, err := createSelector(t1, Daily, valueField)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)

			// string value
			actualString, err := createSelectorString(t1, Daily, valueField, "value")
			So(actualString, ShouldResemble, expectedString)
			So(err, ShouldBeNil)
		})

		Convey("for an hourly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			expectedStart := time.Date(2015, time.March, 20, 0, 0, 0, 0, time.UTC)
			valueField := "field"

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField},
				{"value", bson.D{{"$exists", false}}}}

			actual, err := createSelector(t1, Hourly, valueField)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a minutely metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			expectedStart := time.Date(2015, time.March, 20, 14, 0, 0, 0, time.UTC)
			valueField := "field"

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField},
				{"value", bson.D{{"$exists", false}}}}

			actual, err := createSelector(t1, Minutely, valueField)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a secondly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			expectedStart := time.Date(2015, time.March, 20, 14, 35, 0, 0, time.UTC)
			valueField := "field"

			expected := bson.D{{"start", expectedStart}, {"valueField", valueField},
				{"value", bson.D{{"$exists", false}}}}

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
			expected := bson.D{{"$inc", bson.D{{"total", valueFloat}}},
				{"$inc", bson.D{{"month.3", valueFloat}}}}

			actual, err := createUpdate(t1, Monthly, valueFloat)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a daily metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)
			expected := bson.D{{"$inc", bson.D{{"total", valueFloat}}},
				{"$inc", bson.D{{"day.20", valueFloat}}}}

			actual, err := createUpdate(t1, Daily, valueFloat)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a hourly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)
			expected := bson.D{{"$inc", bson.D{{"total", valueFloat}}},
				{"$inc", bson.D{{"hour.14", valueFloat}}}}

			actual, err := createUpdate(t1, Hourly, valueFloat)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a minutely metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)
			expected := bson.D{{"$inc", bson.D{{"total", valueFloat}}},
				{"$inc", bson.D{{"minute.35", valueFloat}}}}

			actual, err := createUpdate(t1, Minutely, valueFloat)

			So(actual, ShouldResemble, expected)
			So(err, ShouldBeNil)
		})

		Convey("for a secondly metric", func() {
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			value := 20
			valueFloat := convert.ToFloat64(value)
			expected := bson.D{{"$inc", bson.D{{"total", valueFloat}}},
				{"$inc", bson.D{{"second.2", valueFloat}}}}

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

		actual, meta, err := createSingleUpdate(inputDoc, t1, granularity, rule)

		So(err, ShouldBeNil)
		So(meta, ShouldBeNil)
		So(actual, ShouldResemble, expected)
	})
	Convey("Create a single update request object for a string value", t, func() {
		inputDoc := bson.D{{"price", "foo"}}
		t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
		granularity := Daily
		rule := Rule{
			ValueField: "price",
		}

		expectedSelector, err := createSelectorString(t1, granularity, rule.ValueField, "foo")
		So(err, ShouldBeNil)
		expectedUpdate, err := createUpdate(t1, granularity, convert.ToFloat64(1))
		So(err, ShouldBeNil)
		expectedMetaUpdate := &messages.SingleUpdate{
			Selector: bson.D{{"_id", "metadata"}},
			Update:   bson.D{{"$set", bson.D{{rule.ValueField + ".foo", true}}}},
			Upsert:   true,
			Multi:    false,
		}

		expected := &messages.SingleUpdate{
			Selector: expectedSelector,
			Update:   expectedUpdate,
			Upsert:   true,
			Multi:    false,
		}

		actual, meta, err := createSingleUpdate(inputDoc, t1, granularity, rule)

		So(err, ShouldBeNil)
		So(meta, ShouldResemble, expectedMetaUpdate)
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

			_, _, err := createSingleUpdate(inputDoc, t1, granularity, rule)

			So(err, ShouldNotBeNil)
		})
		Convey("for a doc without a value field", func() {
			inputDoc := bson.D{{"hello", "hi"}}
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			granularity := "err"
			rule := Rule{
				ValueField: "price",
			}

			_, _, err := createSingleUpdate(inputDoc, t1, granularity, rule)

			So(err, ShouldNotBeNil)

		})
		Convey("for a doc whose value field is not a number", func() {
			inputDoc := bson.D{{"price", "hi"}}
			t1 := time.Date(2015, time.March, 20, 14, 35, 2, 144, time.UTC)
			granularity := "err"
			rule := Rule{
				ValueField: "price",
			}

			_, _, err := createSingleUpdate(inputDoc, t1, granularity, rule)

			So(err, ShouldNotBeNil)
		})
	})
}

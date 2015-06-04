package convert

import (
	. "github.com/mongodbinc-interns/mongoproxy/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBitmask(t *testing.T) {
	SetLogLevel(DEBUG)

	Convey("Read a bitmask", t, func() {

		Convey("an empty bitmask", func() {
			bitMask := int32(0)

			expected := make([]bool, 32)
			match := true
			for i := 0; i < 32; i++ {
				if ReadBit32LE(bitMask, uint(i)) != expected[i] {
					match = false
				}
			}
			So(match, ShouldEqual, true)
		})
		Convey("an initialized bitmask", func() {
			bitMask := int32(8)

			expected := make([]bool, 32)
			match := true
			expected[3] = true
			for i := 0; i < 32; i++ {
				if ReadBit32LE(bitMask, uint(i)) != expected[i] {
					match = false
				}
			}
			So(match, ShouldEqual, true)
		})

	})
	Convey("Set a bitmask", t, func() {
		bitMask := int32(0)

		expected := make([]bool, 32)
		expected[3] = true
		bitMask = WriteBit32LE(bitMask, 3, true)
		match := true
		for i := 0; i < 32; i++ {
			if ReadBit32LE(bitMask, uint(i)) != expected[i] {
				match = false
			}
		}
		So(match, ShouldEqual, true)

		match = true
		bitMask = WriteBit32LE(bitMask, 3, true)
		for i := 0; i < 32; i++ {
			if ReadBit32LE(bitMask, uint(i)) != expected[i] {
				match = false
			}
		}
		So(match, ShouldEqual, true)

		match = true
		expected[3] = false
		bitMask = WriteBit32LE(bitMask, 3, false)
		for i := 0; i < 32; i++ {
			if ReadBit32LE(bitMask, uint(i)) != expected[i] {
				match = false
			}
		}
		So(match, ShouldEqual, true)
	})
}

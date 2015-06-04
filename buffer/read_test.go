package buffer

import (
	"encoding/binary"
	"github.com/mongodbinc-interns/mongoproxy/mock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// Test a 64 bit LE
func TestRead64BitLE(t *testing.T) {
	Convey("Test a bunch of values", t, func() {
		values := []uint64{31415926, 0, 1, 10, 12, 15, 20, 99}
		for i := 0; i < len(values); i++ {
			bs := make([]byte, 8)
			binary.LittleEndian.PutUint64(bs, values[i])

			m := mock.MockIO{
				Input:  bs,
				Output: make([]byte, 0)}
			m.New()

			val, _ := ReadInt64LE(&m)
			So(val, ShouldEqual, values[i])
		}

	})
}

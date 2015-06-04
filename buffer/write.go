package buffer

import (
	"bytes"
	"encoding/binary"
)

// WriteToBuf takes in a buffer and writes the data as bytes to the buffer
// in the order provided in the arguments. Returns an error if writing to
// the buffer fails for any reason.
func WriteToBuf(buf *bytes.Buffer, data ...interface{}) error {
	for _, d := range data {
		if err := binary.Write(buf, binary.LittleEndian, d); err != nil {
			return err
		}
	}
	return nil
}

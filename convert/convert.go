// Package convert contains functions to convert variables of one type into another.
package convert

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

// ConvertToInt32LE converts the first four bytes of a slice to a 32-bit little endian integer.
func ConvertToInt32LE(byteSlice []byte) int32 {
	return int32(
		(uint32(byteSlice[0]) << 0) |
			(uint32(byteSlice[1]) << 8) |
			(uint32(byteSlice[2]) << 16) |
			(uint32(byteSlice[3]) << 24))
}

// ConvertToInt64LE converts the first eight bytes of a slice to a 64-bit little endian integer.
func ConvertToInt64LE(byteSlice []byte) int64 {
	return int64(
		(uint64(byteSlice[0]) << 0) |
			(uint64(byteSlice[1]) << 8) |
			(uint64(byteSlice[2]) << 16) |
			(uint64(byteSlice[3]) << 24) |
			(uint64(byteSlice[4]) << 32) |
			(uint64(byteSlice[5]) << 40) |
			(uint64(byteSlice[6]) << 48) |
			(uint64(byteSlice[7]) << 56))
}

// ReadBit32LE reads a boolean off of a 32-bit integer bitmask in little endian format
// by reading the bit at position n.
func ReadBit32LE(bitMask int32, n uint) bool {
	if n > 31 {
		return false
	}
	t := (bitMask) & (1 << n)
	if t > 0 {
		return true
	}
	return false
}

// WriteBit32LE writes a boolean as 0 or 1 to the given position n of an integer bitmask,
// and returns the new bitmask with the value set.
func WriteBit32LE(bitMask int32, n uint, value bool) int32 {
	if n > 31 {
		return bitMask
	}

	newBitMask := bitMask
	if value {
		newBitMask = bitMask | (1 << n)
	} else {
		newBitMask = bitMask &^ (1 << n)
	}

	return newBitMask
}

// ToInt converts an interface{} to an int. Will also convert float32 and float64
// to integers. A default value can be provided
// if the conversion fails, otherwise 0 will be returned. Any argument after
// the 2nd one will be ignored.
func ToInt(in interface{}, def ...int) int {
	n, ok := in.(int)
	if ok {
		return n
	}
	n2, ok := in.(float32)
	if ok {
		return int(n2)
	}
	n3, ok := in.(float64)
	if ok {
		return int(n3)
	}

	if len(def) == 0 {
		return 0
	}
	return def[0]

}

// ToInt32 converts an interface{} to an int32. A default value can be provided
// if the conversion fails, otherwise 0 will be returned. Any argument after
// the 2nd one will be ignored.
func ToInt32(in interface{}, def ...int32) int32 {
	n, ok := in.(int32)
	if !ok {

		// check to see if it is an int, and cast to int32 as needed
		if len(def) > 0 {
			return int32(ToInt(in, int(def[0])))
		}
		return int32(ToInt(in))

	}
	return n
}

// ToInt64 converts an interface{} to an int64. A default value can be provided
// if the conversion fails, otherwise 0 will be returned. Any argument after
// the 2nd one will be ignored.
func ToInt64(in interface{}, def ...int64) int64 {
	n, ok := in.(int64)
	if !ok {
		// check to see if it is an int, and cast to int64 as needed
		if len(def) > 0 {
			return int64(ToInt(in, int(def[0])))
		}
		return int64(ToInt(in))
	}
	return n
}

// ToBool converts an interface{} to a bool. A default value can be provided
// if the conversion fails. Any argument after the 2nd one will be ignored.
func ToBool(in interface{}, def ...bool) bool {
	b, ok := in.(bool)
	if !ok {
		if len(def) == 0 {
			return false
		}
		return def[0]
	}
	return b
}

// ToString converts an interface{} to a string. A default value can be provided
// if the conversion fails. Any argument after the 2nd one will be ignored.
func ToString(in interface{}, def ...string) string {
	b, ok := in.(string)
	if !ok {
		if len(def) == 0 {
			return ""
		}
		return def[0]
	}
	return b
}

// ToBSONDoc converts an interface{} to a bson.D. Nil is returned if the
// conversion fails.
func ToBSONDoc(in interface{}) bson.D {
	d, ok := in.(bson.D)
	if !ok {
		return nil
	}
	return d
}

// ToBSONMap converts an interface{} to a bson.M. Nil is returned if the
// conversion fails.
func ToBSONMap(in interface{}) bson.M {
	m, ok := in.(bson.M)
	if !ok {
		return nil
	}
	return m
}

// ConvertToBSONMapSlice converts an []interface{}, []bson.D, or []bson.M slice to a []bson.M
// slice (assuming that all contents are either bson.M or bson.D objects)
func ConvertToBSONMapSlice(input interface{}) ([]bson.M, error) {

	inputBSONM, ok := input.([]bson.M)
	if ok {
		return inputBSONM, nil
	}

	inputBSOND, ok := input.([]bson.D)
	if ok {
		// just convert all of the bson.D documents to bson.M
		d := make([]bson.M, len(inputBSOND))
		for i := 0; i < len(inputBSOND); i++ {
			doc := inputBSOND[i]
			d[i] = doc.Map()
		}
		return d, nil
	}

	inputInterface, ok := input.([]interface{})
	if ok {
		d := make([]bson.M, len(inputInterface))
		for i := 0; i < len(inputInterface); i++ {
			doc := inputInterface[i]
			docM, ok2 := doc.(bson.M)
			if !ok2 {
				// check if it's a bson.D
				docD, ok3 := doc.(bson.D)
				if ok3 {
					docM = docD.Map()
				} else {
					// error
					return nil, fmt.Errorf("Slice contents aren't BSON objects")
				}
			}

			d[i] = docM
		}
		return d, nil
	}

	return nil, fmt.Errorf("Unsupported input")
}

// ConvertToBSONDocSlice converts an []interface{} to a []bson.D slice
// assuming contents are bson.D objects
func ConvertToBSONDocSlice(input interface{}) ([]bson.D, error) {
	inputBSOND, ok := input.([]bson.D)
	if ok {
		return inputBSOND, nil
	}

	inputInterface, ok := input.([]interface{})
	if ok {
		d := make([]bson.D, len(inputInterface))
		for i := 0; i < len(inputInterface); i++ {
			doc := inputInterface[i]
			docD, ok2 := doc.(bson.D)
			if !ok2 {
				return nil, fmt.Errorf("Slice contents aren't BSON objects")
			}
			d[i] = docD
		}
		return d, nil
	}

	return nil, fmt.Errorf("Unsupported input")
}

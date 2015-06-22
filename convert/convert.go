// Package convert contains functions to convert variables of one type into another.
package convert

import (
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

func ToInt32(in interface{}) int32 {
	n, ok := in.(int32)
	if !ok {
		return 0
	}
	return n
}

func ToInt64(in interface{}) int64 {
	n, ok := in.(int64)
	if !ok {
		return 0
	}
	return n
}

// Converts an interface{} to a bool. A default value can be provided
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

func ToBSONDoc(in interface{}) bson.D {
	d, ok := in.(bson.D)
	if !ok {
		return nil
	}
	return d
}

func ToBSONMap(in interface{}) bson.M {
	m, ok := in.(bson.M)
	if !ok {
		return nil
	}
	return m
}

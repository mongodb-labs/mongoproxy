// Package buffer provides utilities for working with and reading data from
// a connection.
package buffer

import (
	"encoding/binary"
	"fmt"
	. "github.com/mongodbinc-interns/mongoproxy/convert"
	"gopkg.in/mgo.v2/bson"
	"io"
)

// ReadDocument reads a BSON ordered document from a reader, and returns the
// number of bytes in the document and the document itself in bson.D format.
func ReadDocument(reader io.Reader) (docSize int32, document bson.D, err error) {
	// Read the first 4 bytes from the connection
	docSize, err = ReadInt32LE(reader)
	if err != nil {
		return 0, nil, fmt.Errorf("error reading docsize: %v", err)
	}
	if docSize < 4 {
		return 0, nil, fmt.Errorf("docSize too small")
	}
	documentBuffer := make([]byte, docSize-4)
	n, err := io.ReadAtLeast(reader, documentBuffer, int(docSize-4))
	if err != nil && err != io.EOF {
		return 0, nil, fmt.Errorf("error reading document: %v", err)
	}
	if int32(n) != docSize-4 {
		return 0, nil, fmt.Errorf("insufficient bytes read: %v instead of %v", n, docSize-4)
	}
	if n == 0 {
		// if the document size was only the size of the four header bytes for whatever reason,
		// we have an empty document. We still read the 4 bytes, though.
		if docSize == 4 {
			return docSize, bson.D{}, nil
		}
		// erroneously read an empty document
		return 0, nil, io.EOF
	}
	docSizeBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(docSizeBuffer, uint32(docSize))

	// append buffers
	document = bson.D{}
	fullByteSlice := append(docSizeBuffer, documentBuffer...)
	err = bson.Unmarshal(fullByteSlice, &document)
	if err != nil {
		return 0, nil, fmt.Errorf("error unmarshalling query: %v", err)
	}
	return
}

// ReadInt32LE reads a 32-bit integer from a reader with little endian encoding.
func ReadInt32LE(reader io.Reader) (int32, error) {
	// Read the first 4 bytes from the connection
	buffer := make([]byte, 4)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("error reading from connection: %v", err)
	}
	if n != 4 {
		return 0, fmt.Errorf("insufficient data")
	}
	return ConvertToInt32LE(buffer), nil
}

// ReadInt64LE reads a 64-bit long from a reader with little endian encoding.
func ReadInt64LE(reader io.Reader) (int64, error) {
	// Read the first 4 bytes from the connection
	buffer := make([]byte, 8)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("error reading from connection: %v", err)
	}
	if n != 8 {
		return 0, fmt.Errorf("insufficient data")
	}
	return ConvertToInt64LE(buffer), nil
}

// ReadNullTerminatedString continuously reads bytes from a reader until
// it hits the null byte, or it reads in maxSize bytes. An error is returned if
// the string read is longer than maxSize bytes, or if there were any errors
// reading from the buffer. If there are no errors, it returns the number of
// characters (bytes) read and the string that it read.
func ReadNullTerminatedString(reader io.Reader, maxSize int32) (int32, string, error) {
	buffer := make([]byte, 1)
	numRead := int32(0)
	stringBuffer := []byte{}
	for {
		// return the current string if we reach the maximum size
		if numRead >= maxSize {
			return 0, "", fmt.Errorf("read too many bytes")
		}
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return 0, "", fmt.Errorf("error reading null string from connection: %v", err)
		}
		if n != 1 {
			return 0, "", fmt.Errorf("insufficient bytes read")
		}
		if buffer[0] == '\x00' {
			break
		}
		numRead += int32(n)
		stringBuffer = append(stringBuffer, buffer[0])
	}
	return numRead + 1, string(stringBuffer), nil
}

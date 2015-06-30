package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/buffer"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"gopkg.in/mgo.v2/bson"
)

func marshalReplyDocs(reply interface{}, docs []bson.D) ([]byte, error) {

	var replyBytes []byte
	var err error
	if reply != nil {
		replyBytes, err = bson.Marshal(reply)
		if err != nil {
			return nil, fmt.Errorf("error marshaling response document: %v", err)
		}
	} else {
		replyBytes = make([]byte, 0)
	}

	if docs != nil {
		for _, doc := range docs {
			docBytes, err := bson.Marshal(doc)
			if err != nil {
				return nil, fmt.Errorf("error marshaling response document: %v", err)
			}
			replyBytes = append(replyBytes, docBytes...)
		}
	}

	return replyBytes, nil
}

func createResponseHeader(reqHeader MsgHeader) MsgHeader {
	mHeader := MsgHeader{
		ResponseTo: reqHeader.RequestID, // requestID from the original request
		OpCode:     1,
	}
	return mHeader
}

func setMessageSize(resp []byte) []byte {
	// Write actual message size
	respSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(respSize, uint32(len(resp)))
	resp[0] = respSize[0]
	resp[1] = respSize[1]
	resp[2] = respSize[2]
	resp[3] = respSize[3]

	return resp
}

// EncodeBSON encodes a BSON object in an OP_REPLY wire protocol message
// as a response to the request with header reqHeader. Not to be used with
// find or getMore command responses, as it disregards some flags that are important
// to those two commands.
// http://docs.mongodb.org/meta-driver/latest/legacy/mongodb-wire-protocol/
func EncodeBSON(reqHeader MsgHeader, b bson.M) ([]byte, error) {
	resHeader := createResponseHeader(reqHeader)

	// we just return 1 object, which is b.
	numberReturned := 1

	// The default is for the flags to be 8, as the AwaitCapable flag
	// is always set to true after MongoDB 1.6., while the other two flags
	// are not used by generic commands.
	// we should have ways to change it, though.
	flags := int32(8)

	buf := bytes.NewBuffer([]byte{})
	err := buffer.WriteToBuf(buf, resHeader, int32(flags),
		int64(0),              // cursorID. Not used for generic command responses
		int32(0),              // startingFrom. Not used for generic command responses
		int32(numberReturned)) // the number of documents returned, 1 in this case
	if err != nil {
		return nil, fmt.Errorf("error writing prepared response %v", err)
	}

	docBytes, err := marshalReplyDocs(b, nil)
	if err != nil {
		return nil, fmt.Errorf("error marshaling documents")
	}
	resp := append(buf.Bytes(), docBytes...)

	resp = setMessageSize(resp)

	return resp, nil
}

// Encodes a response into a byte slice that represents an OP_REPLY wire protocol message.
func Encode(reqHeader MsgHeader, res ModuleResponse) ([]byte, error) {

	Log(DEBUG, "Response: %#v\n", res)

	// handle error
	hasError := res.CommandError != nil

	// error checking
	if hasError {
		// reply with an error instead of the actual documents
		r := bson.M{}
		r["ok"] = 0
		r["errmsg"] = res.CommandError.Message
		r["code"] = res.CommandError.ErrorCode

		return EncodeBSON(reqHeader, r)
	}

	if res.Writer == nil {
		return nil, fmt.Errorf("No response was returned.")
	}

	return res.Writer.ToBytes(reqHeader)

}

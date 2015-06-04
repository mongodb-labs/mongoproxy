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
			return nil, fmt.Errorf("error marshaling response document: %v\n", err)
		}
	} else {
		replyBytes = make([]byte, 0)
	}

	if docs != nil {
		for _, doc := range docs {
			docBytes, err := bson.Marshal(doc)
			if err != nil {
				return nil, fmt.Errorf("error marshaling response document: %v\n", err)
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

// encodes a BSON object
func EncodeBSON(reqHeader MsgHeader, b bson.M) ([]byte, error) {
	resHeader := createResponseHeader(reqHeader)
	// the 1 comes from res.Reply, which is sent as the first document
	numberReturned := 1

	flags := int32(8) // The default is for the flags to be 8.

	buf := bytes.NewBuffer([]byte{})
	err := buffer.WriteToBuf(buf, resHeader, int32(flags), int64(0), int32(0),
		int32(numberReturned))
	if err != nil {
		return nil, fmt.Errorf("error writing prepared response %v\n", err)
	}
	resp := buf.Bytes()

	docBytes, err := marshalReplyDocs(b, nil)
	if err != nil {
		return nil, fmt.Errorf("error marshaling documents")
	}
	resp = append(resp, docBytes...)

	resp = setMessageSize(resp)

	return resp, nil
}

// Encodes a response into a byte slice that represents an OP_REPLY wire protocol message.
func Encode(reqHeader MsgHeader, res ModuleResponse) ([]byte, error) {

	Log(DEBUG, "%#v\n", res)

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

	return res.Writer.ToBytes(reqHeader)

}

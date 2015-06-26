package mongod

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/buffer"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"gopkg.in/mgo.v2/bson"
)

func commandToBSONDoc(c messages.Command) bson.D {
	nameArg, ok := c.Args[c.CommandName]
	if !ok {
		nameArg = 1
	}
	args := bson.D{
		{c.CommandName, nameArg},
	}

	for arg, value := range c.Args {
		if arg != c.CommandName {
			args = append(args, bson.DocElem{arg, value})
		}
	}

	return args
}

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

func CommandToBytes(c messages.Command) ([]byte, error) {
	// we're still creating a request message, so no need to make a response header.

	header := messages.MsgHeader{
		RequestID: c.ID(),
		OpCode:    int32(2004),
	}

	// commands ignore flags.
	flags := int32(0)
	nameSpace := c.Database + ".$cmd"
	nameSpaceBytes := []byte(nameSpace)
	nameSpaceBytes = append(nameSpaceBytes, byte('\x00'))

	buf := bytes.NewBuffer([]byte{})
	err := buffer.WriteToBuf(buf, header, int32(flags), nameSpaceBytes, int32(0), int32(1))
	if err != nil {
		return nil, fmt.Errorf("error writing prepared response %v\n", err)
	}

	// encode the args with the OP_QUERY for now.
	args := commandToBSONDoc(c)
	docBytes, err := marshalReplyDocs(args, nil)
	if err != nil {
		return nil, fmt.Errorf("error marshaling documents")
	}

	resp := append(buf.Bytes(), docBytes...)

	resp = setMessageSize(resp)

	return resp, nil
}

func FindToBytes(f messages.Find) ([]byte, error) {
	/*
		header := messages.MsgHeader{
			RequestID: f.ID(),
			OpCode:    int32(2004),
		}
	*/
	return nil, fmt.Errorf("unimplemented.")
}

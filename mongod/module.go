package mongod

import (
	"encoding/binary"
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/buffer"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/modules"
	"gopkg.in/mgo.v2/bson"
	"net"
)

var port = 8000

type RawBytesWriter struct {
	Data []byte
}

func (r RawBytesWriter) ToBSON() bson.M {
	return bson.M{
		"data": r.Data,
	}
}

func (r RawBytesWriter) ToBytes(header messages.MsgHeader) ([]byte, error) {
	return r.Data, nil
}

type MongodModule struct{}

func (m MongodModule) Process(req messages.Requester, res messages.Responder,
	next modules.PipelineFunc) {

	// takes the request, throws it back into a wire protocol message, and
	// sends it to the responder
	switch req.Type() {
	case messages.CommandType:
		command := req.(messages.Command)
		b, err := CommandToBytes(command)
		if err != nil {
			Log(ERROR, "%#v\n", err)
			next(req, res)
			return
		}
		conn, err := net.Dial("tcp", fmt.Sprintf(":%v", port))
		Log(DEBUG, "%#v\n", b)
		conn.Write(b)

		// listen to reply
		// read the first four bytes
		docSize, err := buffer.ReadInt32LE(conn)
		if err != nil {
			Log(ERROR, "%#v\n", err)
			conn.Close()
			next(req, res)
			return
		}

		msg := make([]byte, docSize-4)
		n, err := conn.Read(msg)

		if n < len(msg) {
			Log(ERROR, "Read in too few bytes\n")
			conn.Close()
			next(req, res)
			return
		}
		if err != nil {
			conn.Close()
			next(req, res)
			return
		}

		reqSize := make([]byte, 4)
		binary.LittleEndian.PutUint32(reqSize, uint32(docSize))

		wholeMsg := append(reqSize, msg...)

		r := RawBytesWriter{
			Data: wholeMsg,
		}

		res.Write(r)

		err = conn.Close()
		if err != nil {
			Log(ERROR, "%v\n", err)
		}

	default:
		Log(ERROR, "Unsupported operation")
	}
	next(req, res)

}

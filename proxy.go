package mongoproxy

import (
	"fmt"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/server"
	"net"
)

func Start(port int, pipeline server.PipelineFunc) {

	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		Log(ERROR, "Error listening on port %v: %v\n", port, err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			Log(ERROR, "error accepting connection: %v\n", err)
			continue
		}

		Log(NOTICE, "accepted connection from: %v\n", conn.RemoteAddr())
		go handleConnection(conn, pipeline)
	}

}

func handleConnection(conn net.Conn, pipeline server.PipelineFunc) {
	for {

		message, msgHeader, err := messages.Decode(conn)

		if err != nil {
			Log(ERROR, "%#v", err)
			conn.Close()
			return
		}

		res := &messages.ModuleResponse{}
		pipeline(message, res)

		Log(DEBUG, "%#v\n", res)

		bytes, err := messages.Encode(msgHeader, *res)

		// update, delete, and insert messages do not have a response, so we continue and write the
		// response on the getLastError that will be called immediately after. Kind of a hack.
		if msgHeader.OpCode == messages.OP_UPDATE || msgHeader.OpCode == messages.OP_INSERT ||
			msgHeader.OpCode == messages.OP_DELETE {
			continue
		}
		if err != nil {
			Log(ERROR, "%#v", err)
			conn.Close()
			return
		}
		_, err = conn.Write(bytes)
		if err != nil {
			Log(ERROR, "%#v", err)
			conn.Close()
			return
		}

	}
}

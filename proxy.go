package mongoproxy

import (
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/convert"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/server"
	_ "github.com/mongodbinc-interns/mongoproxy/server/config"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net"
)

func Start(port int, chain *server.ModuleChain) {

	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		Log(ERROR, "Error listening on port %v: %v\n", port, err)
		return
	}

	pipeline := server.BuildPipeline(chain)

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

func StartWithConfig(port int, config bson.M) {
	chain := server.CreateChain()
	modules, err := convert.ConvertToBSONMapSlice(config["modules"])
	if err != nil {
		Log(ERROR, "Invalid module configuration, or proxy was started with no modules.")
		return
	}

	for i := 0; i < len(modules); i++ {
		moduleName := convert.ToString(modules[i]["name"])
		moduleType, ok := server.Registry[moduleName]
		if !ok {
			Log(WARNING, "Module doesn't exist in the registry: %v", moduleName)
			continue // module doesn't exist
		}
		module := moduleType.New()

		// TODO: allow links to other collections
		moduleConfig := convert.ToBSONMap(modules[i]["config"])
		err := module.Configure(moduleConfig)
		if err != nil {
			Log(WARNING, "Invalid configuration for module: %v", moduleName)
			continue
		}
		chain.AddModule(module)
	}
	Start(port, chain)
}

func handleConnection(conn net.Conn, pipeline server.PipelineFunc) {
	for {

		message, msgHeader, err := messages.Decode(conn)

		if err != nil {
			if err != io.EOF {
				Log(ERROR, "Decoding error: %#v", err)
			}
			conn.Close()
			return
		}

		Log(DEBUG, "Request: %#v\n", message)

		res := &messages.ModuleResponse{}
		pipeline(message, res)

		bytes, err := messages.Encode(msgHeader, *res)

		// update, delete, and insert messages do not have a response, so we continue and write the
		// response on the getLastError that will be called immediately after. Kind of a hack.
		if msgHeader.OpCode == messages.OP_UPDATE || msgHeader.OpCode == messages.OP_INSERT ||
			msgHeader.OpCode == messages.OP_DELETE {
			Log(INFO, "Continuing on OpCode: %v", msgHeader.OpCode)
			continue
		}
		if err != nil {
			Log(ERROR, "Encoding error: %#v", err)
			conn.Close()
			return
		}
		_, err = conn.Write(bytes)
		if err != nil {
			Log(ERROR, "Error writing to connection: %#v", err)
			conn.Close()
			return
		}

	}
}

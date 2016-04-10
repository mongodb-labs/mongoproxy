package mongoproxy

import (
	"encoding/json"
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/convert"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/server"
	_ "github.com/mongodbinc-interns/mongoproxy/server/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"net"
)

// ParseConfigFromFile takes a filename for a JSON file, and returns a configuration
// object from the file, and an error if there was an error reading or unmarshalling the file.
func ParseConfigFromFile(configFilename string) (bson.M, error) {
	var result bson.M

	file, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return nil, fmt.Errorf("Error reading configuration file: %v", err)
	}

	err = json.Unmarshal(file, &result)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON Configuration %v: %v", configFilename, err)
	}
	return result, nil
}

// ParseConfigFromDB takes a MongoURI string and a namespace to query a MongoDB instance
// for a configuration document, and returns the document and any errors finding the document.
// If there are multiple documents in the collection, by default the latest one (the first result
// in a find) will be returned.
func ParseConfigFromDB(mongoURI string, configNamespace string) (bson.M, error) {
	var result bson.M

	mongoSession, err := mgo.Dial(mongoURI)
	defer mongoSession.Close()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to MongoDB instance: %v", err)
	}

	database, collection, err := messages.ParseNamespace(configNamespace)
	if err != nil {
		return nil, fmt.Errorf("Invalid namespace: %v", err)
	}

	err = mongoSession.DB(database).C(collection).Find(bson.M{}).One(&result)
	if err != nil {
		return nil, fmt.Errorf("Error querying MongoDB for configuration: %v", err)
	}

	return result, nil
}

// Start starts the server at the provided port and with the given module chain.
func Start(port int, chain *server.ModuleChain) {

	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		Log(ERROR, "Error listening on port %v: %v", port, err)
		return
	}

	pipeline := server.BuildPipeline(chain)
	Log(INFO, "Server running on port %v", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			Log(ERROR, "error accepting connection: %v", err)
			continue
		}

		Log(NOTICE, "accepted connection from: %v", conn.RemoteAddr())
		go handleConnection(conn, pipeline)
	}

}

// StartWithConfig starts the server at the provided port, creating a module chaine
// with the given configuration.
func StartWithConfig(port int, config bson.M) {
	chain := server.CreateChain()
	var modules []bson.M
	var err error
	modulesRaw, ok := config["modules"]
	if ok {
		modules, err = convert.ConvertToBSONMapSlice(modulesRaw)
		if err != nil {
			Log(WARNING, "Invalid module configuration: %v. Proxy will start with no modules.", err)
		}
	} else {
		Log(WARNING, "No modules provided. Proxy will start without modules.")
	}

	for i := 0; i < len(modules); i++ {
		moduleNameRaw, ok := modules[i]["name"]
		if !ok {
			Log(WARNING, "Module in configuration does not have a name")
			continue
		}
		moduleName := convert.ToString(moduleNameRaw)
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
			Log(WARNING, "Invalid configuration for module %v: %v", moduleName, err)
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
				Log(ERROR, "Decoding error: %v", err)
			}
			conn.Close()
			return
		}

		Log(DEBUG, "Request: %#v", message)

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
			Log(ERROR, "Encoding error: %v", err)
			conn.Close()
			return
		}
		_, err = conn.Write(bytes)
		if err != nil {
			Log(ERROR, "Error writing to connection: %v", err)
			conn.Close()
			return
		}

	}
}

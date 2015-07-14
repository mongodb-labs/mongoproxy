package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/convert"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi/frontend"
	"github.com/mongodbinc-interns/mongoproxy/server"
	_ "github.com/mongodbinc-interns/mongoproxy/server/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
)

var (
	port            int
	logLevel        int
	mongoURI        string
	configNamespace string
	configFilename  string
)

func parseFlags() {
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.IntVar(&logLevel, "logLevel", 3, "verbosity for logging")
	flag.StringVar(&mongoURI, "m", "mongodb://localhost:27017", "MongoDB instance to connect to for configuration.")
	flag.StringVar(&configNamespace, "c", "test.config", "Namespace to query for configuration.")
	flag.StringVar(&configFilename, "f", "", "JSON config filename. If set, will be used instead of mongoDB configuration.")
	flag.Parse()
}

func main() {

	parseFlags()
	SetLogLevel(logLevel)

	// grab config file
	var result bson.M
	if len(configFilename) == 0 {
		mongoSession, err := mgo.Dial(mongoURI)
		if err != nil {
			Log(ERROR, "Error connecting to MongoDB instance: %#v\n", err)
			return
		}

		database, collection, err := messages.ParseNamespace(configNamespace)
		if err != nil {
			Log(ERROR, "Invalid namespace: %#v\n", err)
			return
		}

		err = mongoSession.DB(database).C(collection).Find(bson.M{}).One(&result)
		if err != nil {
			Log(ERROR, "Error querying MongoDB for configuration: %#v\n", err)
			return
		}
	} else {
		file, err := ioutil.ReadFile(configFilename)
		if err != nil {
			Log(ERROR, "Error reading configuration file: %#v\n", err)
			return
		}

		json.Unmarshal(file, &result)
		if result == nil {
			Log(ERROR, "Invalid JSON configuration")
		}
	}

	modules, err := convert.ConvertToBSONMapSlice(result["modules"])
	if err != nil {
		Log(ERROR, "Invalid module configuration: %v.", err)
		return
	}

	var module server.Module
	var moduleConfig bson.M
	for i := 0; i < len(modules); i++ {
		moduleName := convert.ToString(modules[i]["name"])
		if moduleName == "bi" {
			module = server.Registry["bi"].New()
			// TODO: allow links to other collections
			moduleConfig = convert.ToBSONMap(modules[i]["config"])
			break
		}
	}
	if module == nil {
		Log(ERROR, "No BI module found in configuration")
		return
	}

	module.Configure(moduleConfig)
	biModule, ok := module.(*bi.BIModule)
	if !ok {
		Log(ERROR, "Not a BI Module")
	}

	r := frontend.Start(biModule, "modules/bi/frontend")
	r.Run(fmt.Sprintf(":%v", port))
}

package main

import (
	"encoding/json"
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
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
	flag.IntVar(&port, "port", 8124, "port to listen on")
	flag.IntVar(&logLevel, "logLevel", 3, "verbosity for logging")
	flag.StringVar(&mongoURI, "m", "mongodb://localhost:27017",
		"MongoDB instance to connect to for configuration.")
	flag.StringVar(&configNamespace, "c", "test.config",
		"Namespace to query for configuration.")
	flag.StringVar(&configFilename, "f", "",
		"JSON config filename. If set, will be used instead of mongoDB configuration.")
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

	mongoproxy.StartWithConfig(port, result)
}

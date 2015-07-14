package main

import (
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	_ "github.com/mongodbinc-interns/mongoproxy/server/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	port            int
	logLevel        int
	mongoURI        string
	configNamespace string
)

func parseFlags() {
	flag.IntVar(&port, "port", 8124, "port to listen on")
	flag.IntVar(&logLevel, "logLevel", 3, "verbosity for logging")
	flag.StringVar(&mongoURI, "m", "mongodb://localhost:27017", "MongoDB instance to connect to for configuration.")
	flag.StringVar(&configNamespace, "c", "test.config", "Namespace to query for configuration.")

	flag.Parse()
}

func main() {

	parseFlags()
	SetLogLevel(logLevel)

	// grab config file
	mongoSession, err := mgo.Dial(mongoURI)
	if err != nil {
		Log(ERROR, "%#v\n", err)
		return
	}

	database, collection, err := messages.ParseNamespace(configNamespace)
	if err != nil {
		Log(ERROR, "%#v\n", err)
		return
	}

	var result bson.M
	err = mongoSession.DB(database).C(collection).Find(bson.M{}).One(&result)
	if err != nil {
		Log(ERROR, "%#v\n", err)
		return
	}

	mongoproxy.StartWithConfig(port, result)
}

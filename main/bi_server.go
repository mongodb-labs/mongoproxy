package main

import (
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/server"
	"gopkg.in/mgo.v2/bson"
)

var (
	port     int
	logLevel int
)

func parseFlags() {
	flag.IntVar(&port, "port", 8124, "port to listen on")
	flag.IntVar(&logLevel, "logLevel", DEBUG, "verbosity for logging")

	flag.Parse()
}

func main() {

	parseFlags()
	SetLogLevel(logLevel)

	// initialize the mockule
	mockule := server.Registry["mockule"]

	// initialize BI module
	biModule := server.Registry["bi"]

	ruleBSON := bson.M{
		"origin":          "test.foo",
		"prefix":          "db.metrics",
		"timeGranularity": []string{bi.Daily, bi.Secondly},
		"valueField":      "price",
	}

	var biConfig = bson.M{}
	connection := bson.M{}
	connection["addresses"] = []string{"localhost:27017"}
	biConfig["connection"] = connection
	biConfig["rules"] = []bson.M{ruleBSON}

	// initialize the pipeline
	chain := server.CreateChain()

	chain.AddModule(biModule)
	biModule.Configure(biConfig)
	chain.AddModule(mockule)

	mongoproxy.Start(port, chain)
}

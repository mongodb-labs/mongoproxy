package main

import (
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi"
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

	modules := make([]bson.M, 2)
	modules[0] = bson.M{
		"name":   "bi",
		"config": biConfig,
	}
	modules[1] = bson.M{
		"name":   "mockule",
		"config": bson.M{},
	}

	// initialize the pipeline
	config := bson.M{
		"modules": modules,
	}

	mongoproxy.StartWithConfig(port, config)
}

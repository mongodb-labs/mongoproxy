package main

import (
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
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

	connection := bson.M{}
	connection["addresses"] = []string{"localhost:27017"}

	modules := make([]bson.M, 1)
	modules[0] = bson.M{
		"name":   "mongod",
		"config": connection,
	}

	mongoproxy.StartWithConfig(port, bson.M{
		"modules": modules,
	})
}

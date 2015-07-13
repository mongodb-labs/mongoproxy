package main

import (
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/modules/mongod"
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
	module := mongod.MongodModule{}
	connection := bson.M{}
	connection["addresses"] = []string{"localhost:27017"}

	// initialize the pipeline
	chain := server.CreateChain()
	chain.AddModule(module)
	module.Configure(connection)

	mongoproxy.Start(port, chain)
}

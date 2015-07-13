package main

import (
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/server"
	_ "github.com/mongodbinc-interns/mongoproxy/server/config"
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

	module := server.Registry["mongod"].New()

	connection := bson.M{}
	connection["addresses"] = []string{"localhost:27017"}

	// initialize the pipeline
	chain := server.CreateChain()
	module.Configure(connection)
	chain.AddModule(module)

	mongoproxy.Start(port, chain)
}

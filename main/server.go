package main

import (
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/modules/mockule"
	"github.com/mongodbinc-interns/mongoproxy/server"
)

var (
	port     int
	logLevel int
)

func parseFlags() {
	flag.IntVar(&port, "port", 8124, "port to listen on")
	flag.IntVar(&logLevel, "logLevel", 3, "verbosity for logging")

	flag.Parse()
}

func main() {

	parseFlags()
	SetLogLevel(logLevel)

	// initialize the mockule
	mockule := mockule.Mockule{}

	// initialize the pipeline
	chain := server.CreateChain()
	chain.AddModule(mockule)
	pipeline := server.BuildPipeline(chain)

	mongoproxy.Start(port, pipeline)
}

package main

import (
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/modules"
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
	mockule := mongoproxy.Mockule{}

	// initialize the pipeline
	chain := modules.CreateChain()
	modules.AddModule(chain, mockule)
	pipeline := modules.BuildPipeline(chain)

	mongoproxy.Start(port, pipeline)
}

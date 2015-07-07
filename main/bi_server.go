package main

import (
	"flag"
	"github.com/mongodbinc-interns/mongoproxy"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi"
	"github.com/mongodbinc-interns/mongoproxy/modules/mockule"
	"github.com/mongodbinc-interns/mongoproxy/server"
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
	mockule := mockule.Mockule{}

	// initialize BI module
	biModule := bi.BIModule{}

	t := make([]string, 2)
	t[0] = bi.Daily
	t[1] = bi.Minutely

	rule := bi.Rule{
		OriginDatabase:    "test",
		OriginCollection:  "foo",
		PrefixDatabase:    "db",
		PrefixCollection:  "metrics",
		TimeGranularities: t,
		ValueField:        "price",
	}
	biModule.Rules = append(biModule.Rules, rule)

	// initialize the pipeline
	chain := server.CreateChain()

	chain.AddModule(biModule)
	chain.AddModule(mockule)

	pipeline := server.BuildPipeline(chain)

	mongoproxy.Start(port, pipeline)
}

package main

import (
	"flag"
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy"
	"github.com/mongodbinc-interns/mongoproxy/convert"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi/frontend"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi/frontend/controllers"
	_ "github.com/mongodbinc-interns/mongoproxy/server/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	port            int
	logLevel        int
	mongoURI        string
	configNamespace string
	configFilename  string
)

func parseFlags() {
	flag.IntVar(&port, "port", 8080, "port to listen on")
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
	// Currently, it will take the configuration of the first BI module found in the chain.
	var result bson.M
	var configLocation *controllers.ConfigLocation
	var err error
	if len(configFilename) == 0 {
		result, err = mongoproxy.ParseConfigFromDB(mongoURI, configNamespace)
		mongoSession, err := mgo.Dial(mongoURI)

		if err == nil {
			database, collection, err := messages.ParseNamespace(configNamespace)
			if err == nil {
				configLocation = &controllers.ConfigLocation{
					Session:    mongoSession,
					Database:   database,
					Collection: collection,
				}
			} else {
				Log(WARNING, "Invalid namespace for configuration location.")
			}

		} else {
			Log(WARNING, "Unable to find configuration location.")
		}

	} else {
		result, err = mongoproxy.ParseConfigFromFile(configFilename)
	}

	if err != nil {
		Log(WARNING, "%v", err)
	}

	modules, err := convert.ConvertToBSONMapSlice(result["modules"])
	if err != nil {
		Log(WARNING, "Invalid module configuration: %v.", err)
	}

	var moduleConfig bson.M
	if modules != nil {
		for i := 0; i < len(modules); i++ {
			moduleName := convert.ToString(modules[i]["name"])
			if moduleName == "bi" {
				// TODO: allow links to other collections
				moduleConfig = convert.ToBSONMap(modules[i]["config"])
				break
			}
		}
	}

	if moduleConfig == nil {
		Log(WARNING, "No BI module found in configuration")
	}
	r, err := frontend.Start(moduleConfig, "modules/bi/frontend", configLocation)
	if err != nil {
		Log(ERROR, "Error starting frontend: %v", err)
		return
	}
	r.Run(fmt.Sprintf(":%v", port))
}

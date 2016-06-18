# Mongoproxy

A server that speaks the MongoDB wire protocol, and can analyze and transform requests and responses to and from a client and a MongoDB server instance. All requests go through a module pipeline, and modules include a proxy to send requests to `mongod`, as well as a module to collect real-time analytics about the requests passing through the server.

## Building and running

To grab dependencies (should not be necessary, as dependencies are checked into the repository):

	chmod 755 ./vendor.sh # only needs to be done once
	./vendor.sh

Note that currently, mongoproxy requires a specific fork of mgo because it requires some features not yet pulled into the main repo.

To run:

	chmod 755 ./start.sh # only needs to be done once
	./start.sh <options>

`start.sh` sets up the go path, and runs the `main/server.go` file. It's equivalent to the following two commands:

	. ./set_gopath.go
	go run main/server.go

### Configurations

Configurations tell the server which modules to load and run, and also specifies configurations for each module. They can either be a document in a MongoDB instance (as BSON), or a JSON file in the local directory.


By default, the server expects a `mongod` instance to be running on `localhost:27017` with a configuration document in the `test.config` collection. The location for the configuration document can be set as a command line option, and can also be a file. 

Configurations have one field `modules`, which is an array. Each object in the array has a `name` field for the name of the module, and a `config` field for the module's configuration.

A configuration can be found in the project directory named `example_bi_config.json`, which is run with the following command:

	./start.sh -f example_bi_config.json

### Command Line Options

	-port 		Port number to run the server on. Defaults to 8124.
	-logLevel 	Sets verbosity of the logs from 1 to 5, with 1 being least verbose and 5 being the most. Defaults to 3.
	-m 			URL of a mongod server to connect to to retrieve configuration information from. Defaults to localhost:27017
	-c 			Namespace of the collection in the mongod server to retrieve configuration information from. Defaults to test.config
	-f 			Path to a configuration file. If set, the m and c flags are ignored.

## Tests

To run unit tests:
	
	chmod 755 ./test.sh
	./test.sh

To run integration tests:

	# single test
	node tests/test <js file to test>

	# a directory full of test files
	node tests/test_dir <directory of files to test>

## Modules

The only thing that the server does to requests is to translate them into a Go struct. Everything else happens in a module. 

### Usage

The server can use modules that are defined in its configuration.

All modules have a unique name to identify them, which is also used in the `name` field of the configuration. Requests are passed from module to module in the order defined in the configuration, and the response sent to the server, and finally back up to the client.

The following modules are implemented and included in the source:

	mockule 	A mock module that stores insert requests in memory and can dump them back out. It also pretends it is a 1-node replica set.
	mongod 		A module that forwards the request to a MongoDB instance and passes back the response to the server.
	bi 			A module with pre-configured rules that analyzes requests and aggregates them into metrics.

### Developing Modules

All modules implement the `Module` interface, defined in `server/modules.go`. `Configure()` is called at the server startup, and `Process(req, res, next)` is called every time a request passes through the server. 

A module is responsible for calling the next module in the pipeline via the `next` argument in the `Process` function, which is a function that takes two arguments: a request and a response.

Modules also have to be added to the registry in order for the server to know they exist. Each module should live in their own package, and have an `init` function with the following line:

	server.Publish(<Module>)

where `server` is the imported `github.com/mongodb-labs/mongoproxy/server` package.

Then, in the `server/registry.go` file, add the import path of the module to the file preceded by an underscore, to add the module to the registry.

#### Example Module

	package examplemodule

	import (
		"github.com/mongodb-labs/mongoproxy/messages"
		"github.com/mongodb-labs/mongoproxy/server"
		"gopkg.in/mgo.v2/bson"
	)

	type ExampleModule struct {
	}

	func init() {
		server.Publish(ExampleModule{})
	}

	func (m ExampleModule) Name() string {
		return "find-only"
	}

	func (m ExampleModule) Configure(config bson.M) error {
		// configure the module here. Returns an error if configuration fails
		return nil
	}

	// this module will drop all requests except for Find requests
	func (m ExampleModule) Process(req messages.Requester, res messages.Responder,
	next server.PipelineFunc) {
		switch req.Type() {
		case messages.FindType:
			// send the request and response to the next module
			next(req, res)
		default:
			res.Write(<Request dropped>)
			return
	}

This module can then be run with a configuration like the following:

	{
		"modules": [
			...
			{
				"name": "find-only",
				"config": {}
			},
			...
		]
	}

#!/bin/bash

chmod 755 ./set_gopath.sh
. ./set_gopath.sh
go run main/mongod-server.go -logLevel 5

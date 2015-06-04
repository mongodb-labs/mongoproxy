#!/bin/bash

chmod 755 ./set_gopath.sh
. ./set_gopath.sh
go run main/server.go -logLevel 5
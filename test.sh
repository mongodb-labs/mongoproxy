#!/bin/bash

# usage: test.sh <name of package>
chmod 755 ./set_gopath.sh
. ./set_gopath.sh
go test github.com/mongodbinc-interns/mongoproxy/$1 -coverprofile=coverage.out
go tool cover -html=coverage.out
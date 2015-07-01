#!/bin/bash

# usage: test.sh <name of package>
chmod 755 ./set_gopath.sh
. ./set_gopath.sh
go test github.com/mongodbinc-interns/mongoproxy/buffer -coverprofile=coverage.out
go test github.com/mongodbinc-interns/mongoproxy/convert -coverprofile=coverage.out
go test github.com/mongodbinc-interns/mongoproxy/messages -coverprofile=coverage.out
go test github.com/mongodbinc-interns/mongoproxy/server -coverprofile=coverage.out

#!/bin/bash

# usage: test.sh <name of package>
chmod 755 ./set_gopath.sh
. ./set_gopath.sh

packages=(bsonutil buffer convert messages server modules/bi)
for i in ${packages[@]}; do
	go test github.com/mongodbinc-interns/mongoproxy/${i} -coverprofile=coverage.out $1
done

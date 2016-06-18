#!/bin/bash
# usage: test.sh <name of package>
packages=(bsonutil buffer convert messages server modules/bi)
for i in ${packages[@]}; do
	go test github.com/mongodb-labs/mongoproxy/${i} -coverprofile=coverage.out $1
done

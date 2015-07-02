#!/bin/bash

# usage: test.sh <name of package>
chmod 755 ./set_gopath.sh
. ./set_gopath.sh

packages=(buffer convert messages server)
for i in ${packages[@]}; do
	go test github.com/mongodbinc-interns/mongoproxy/${i} -coverprofile=coverage.out	
done

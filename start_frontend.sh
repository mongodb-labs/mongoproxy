#!/bin/bash

chmod 755 ./set_gopath.sh
. ./set_gopath.sh
cd modules/bi/frontend
npm run build
cd ../../../
go run main/frontend.go $@

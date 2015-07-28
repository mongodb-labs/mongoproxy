#!/bin/bash

PROXY_PKG='github.com/mongodbinc-interns/mongoproxy';

setgopath() {
	if [ "Windows_NT" != "$OS" ]; then
		SOURCE_GOPATH=`pwd`.gopath
		VENDOR_GOPATH=`pwd`/vendor

		# set up the $GOPATH to use the vendored dependencies as
		# well as the source for the mongo tools
		rm -rf .gopath/
		mkdir -p .gopath/src/"$(dirname "${PROXY_PKG}")"
		ln -sf `pwd` .gopath/src/$PROXY_PKG
		export GOPATH=`pwd`/.gopath:`pwd`/vendor
	else
		# This is assuming the use of git bash. Cygwin might require a different
		# configuration.
		SOURCE_GOPATH=`pwd`/.gopath
		VENDOR_GOPATH=`pwd`/vendor


		# set up the $GOPATH to use the vendored dependencies as
		# well as the source for the mongo tools
		rm -rf .gopath/
		mkdir -p .gopath/src/"$PROXY_PKG"

		packages=(bsonutil buffer convert log main messages mock modules server tests)
		for i in ${packages[@]}; do
			cp -r `pwd`/${i} .gopath/src/$PROXY_PKG/
		done
		cp * .gopath/src/$PROXY_PKG/
		export GOPATH=`pwd`/.gopath:`pwd`/vendor
	fi;
}

setgopath

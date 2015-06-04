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
		local SOURCE_GOPATH=`pwd`/.gopath
		local VENDOR_GOPATH=`pwd`/vendor
		SOURCE_GOPATH=$(cygpath -w $SOURCE_GOPATH);
		VENDOR_GOPATH=$(cygpath -w $VENDOR_GOPATH);

		# set up the $GOPATH to use the vendored dependencies as
		# well as the source for the mongo tools
		rm -rf .gopath/
		mkdir -p .gopath/src/"$PROXY_PKG"
		cp -r `pwd`/play .gopath/src/$PROXY_PKG
		export GOPATH="$SOURCE_GOPATH;$VENDOR_GOPATH"
	fi;
}

setgopath
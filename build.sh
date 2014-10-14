#!/bin/bash

echo "building eventeserver"
cd eventeserver/
export GOPATH=`pwd`

go get code.google.com/p/go.net/websocket && go get github.com/BurntSushi/xgbutil

go build main.go


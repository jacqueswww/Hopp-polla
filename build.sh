#!/bin/bash

cd ./eventserver
export GOPATH=`pwd`

if [[ ! -d pkg ]]
then
    echo "fetching required packages"
    go get code.google.com/p/go.net/websocket 
    go get github.com/BurntSushi/xgbutil
fi

echo "building eventeserver"
go build main.go


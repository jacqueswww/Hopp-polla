#!/bin/bash

if [[ -e ./eventserver/main ]]
then
    ./eventserver/main
else
    echo "server not built"
fi


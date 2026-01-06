#!/bin/sh
PROJ="kindlesync"

mkdir -p build
case "$1" in
    "native") 
        go build && mv ${PROJ} build/
        echo "OK"
    ;;
    "kual-kindle")
        env GOOS=linux GOARCH=arm GOARM=7 go build && mv ${PROJ} res/kual-extension-kindlesync/kindlesync/
        echo "OK"
    ;;
    *) echo "invalid target!"
    ;;
esac


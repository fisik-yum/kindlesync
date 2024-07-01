#!/bin/sh
KSYNC_LOC="/mnt/us/kindlesync/"
cd $KSYNC_LOC 
case "$1" in
    "sync") ./kindlesync sync 
    ;;
    "refresh") ./kindlesync refresh
    ;;
    *) echo "Invalid option selected"
    ;;
esac


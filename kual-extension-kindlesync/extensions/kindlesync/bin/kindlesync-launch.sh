#!/bin/sh
KSYNC_LOC="/mnt/us/kindlesync/"
cd $KSYNC_LOC 
./kindlesync $1
if [ $? -eq 0 ]; then
    /usr/sbin/eips 10 50 "Success!"
    /usr/sbin/eips ''
else
    /usr/sbin/eips 10 50 "Failure"
    /usr/sbin/eips ''
fi


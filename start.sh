#!/bin/bash

AUTHSERVER_LOG="$GOPATH/src/github.com/patkaehuaea/command/authserver/out/authserver.log"
TIMSERVER_LOG="$GOPATH/src/github.com/patkaehuaea/command/timeserver/out/timeserver.log"

echo "Starting authserver as background process..."
cd $GOPATH/src/github.com/patkaehuaea/command/authserver
$GOPATH/bin/authserver --dumpfile ~/users.json > /dev/null 2>&1 &

if [ "$?" -eq 0 ] ; then
    echo "Starting timeserver as background process.."
    cd $GOPATH/src/github.com/patkaehuaea/command/timeserver
    $GOPATH/bin/timeserver > /dev/null 2>&1 &
    tail -f $AUTHSERVER_LOG -f $TIMESERVER_LOG
else
    echo "Failed to start authserver. Will not start timeserver."
fi
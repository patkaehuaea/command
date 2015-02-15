#!/bin/bash
echo "Starting authserver as background process..."
cd $GOPATH/src/github.com/patkaehuaea/command/authserver
$GOPATH/bin/authserver --dumpfile ~/users.json > /dev/null 2>&1 &

if [ "$?" -eq 0 ] ; then
    echo "Starting timeserver as background process.."
    cd $GOPATH/src/github.com/patkaehuaea/command/timeserver
    $GOPATH/bin/timeserver > /dev/null 2>&1 &
    tail -f $GOPATH/src/github.com/patkaehuaea/command/*/out/*.log
else
    echo "Failed to start authserver. Will not start timeserver."
fi
#!/bin/bash

if pgrep authserver > /dev/null 2>&1 ; then
    echo "Authserver already running!"
else
    echo "Starting authserver as background process..."
    cd $GOPATH/src/github.com/patkaehuaea/command/authserver
    touch out/authserver.log
    $GOPATH/bin/authserver --dumpfile ~/users.json > /dev/null 2>&1 &
fi

if pgrep timeserver > /dev/null 2>&1 ; then
    echo "Timeserver already running!"
else
    echo "Starting timeserver as background process..."
    cd $GOPATH/src/github.com/patkaehuaea/command/timeserver
    touch out/timeserver.log
    $GOPATH/bin/timeserver > /dev/null 2>&1 &
fi

if pgrep authserver > /dev/null 2>&1 && pgrep timeserver > /dev/null 2>&1 ; then
    tail -f $GOPATH/src/github.com/patkaehuaea/command/*/out/*.log
else
    echo "Timeserver and/or authserver not running. Will not tail logs."
fi
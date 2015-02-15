#!/bin/bash

AUTHSERVER="authserver"
TIMESERVER="timeserver"

declare -a processes=($AUTHSERVER $TIMESERVER)

for process in "${processes[@]}"
do
    echo "Killing $process..."
    pkill -f $process
    if [ "$?" -eq 0 ] ; then
        echo "$process stopped successfully."
    else
        echo "Failed to stop $process."
    fi
done
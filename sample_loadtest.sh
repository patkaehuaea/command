#!/bin/bash
# Script is quick and dirty and intended to be run on
# dev machine, then rhel server to verify "completeness".
cd $GOPATH/src/github.com/patkaehuaea/command/
make

kill `pgrep authserver`
echo "Starting authserver."
cd $GOPATH/src/github.com/patkaehuaea/command/authserver
$GOPATH/bin/authserver --dumpfile ~/users.json > /dev/null 2>&1 &

kill `pgrep timeserver`
echo "Starting timeserver."
cd $GOPATH/src/github.com/patkaehuaea/command/timeserver
$GOPATH/bin/timeserver --max-inflight=80 --avg-response-ms=500ms  --deviation-ms=300ms > /dev/null 2>&1 &

cd $GOPATH/src/github.com/patkaehuaea/command/loadgen
echo "time $GOPATH/bin/loadgen --url='http://localhost:8080/time' --runtime=10s --rate=200 --burst=20 --timeout-ms=1000ms"
time $GOPATH/bin/loadgen --url='http://localhost:8080/time' --runtime=10s --rate=200 --burst=20 --timeout-ms=1000ms
echo ""
sleep 10
echo "time $GOPATH/bin/loadgen --url='http://localhost:8080/time' --runtime=10s --rate=2000 --burst=20 --timeout-ms=1000ms"
time $GOPATH/bin/loadgen --url='http://localhost:8080/time' --runtime=10s --rate=2000 --burst=20 --timeout-ms=1000ms
echo ""
sleep 10
echo "time $GOPATH/bin/loadgen --url='http://localhost:8080/time' --runtime=10s --rate=2000 --burst=2 --timeout-ms=500ms"
time $GOPATH/bin/loadgen --url='http://localhost:8080/time' --runtime=10s --rate=2000 --burst=2 --timeout-ms=500ms
[GENERAL]

Both authserver and timeserver can be launched using the start.sh script
or calling each executable directly. Launching with user defined parameters requires manual invocation.
Defaults for both applications are defined in `$GOPATH/src/github.com/patkaehuaea/command/config/config.go`.
Timeserver will run on defaults with no parameters passed, but authserver must receive the --dumpfile flag
or execution will halt. The functionality implemented by each server is described in the assignment-04
specification: http://goo.gl/ZmC2YD. Please review spec for list of supported features.


[CAVEATS]

1. The following timeserver flags have been implemented as time.Duration:

--authtimeout-ms
--avg-response-ms
--deviation-ms

Example usage (from timeserver directory):

$ $GOPATH/bin/timeserver --authtimeout-ms 1500ms --avg-response-ms 1000ms --deviation-ms 500ms


2. The following authserver flags have been implemented as time.Duration:

--checkpoint-interval

Example usage (from authserver directory):

$ $GOPATH/bin/authserver --dumpfile ~/users.json --checkpoint-interval 60s


[UNPACK]


Unzip go.zip and replace the go folder in your $GOPATH.
If $GOPATH == ~/go then:

$ mv ~/go ~/go_bk
$ unzip PATH_TO_DOWNLOADED_FILE/go.zip -d ~/


[EASY]


CD into the comand directory and make:

$ cd $GOPATH/src/github.com/patkaehuaea/command/
$ make


Run both servers with defaults defined in config.go:

$ chmod u+x start.sh
$ ./start.sh


Stop both servers:

$ chmod u+x stop.sh
$ ./stop.sh


[MANUAL]


CD into the comand directory and make:

$ cd $GOPATH/src/github.com/patkaehuaea/command/
$ make


Instructions for running authserver (dumpfile required):

$ cd $GOPATH/src/github.com/patkaehuaea/command/authserver
$ go install
$ $GOPATH/bin/authserver --dumpfile ~/users.json


Instructions for running timeserver:

$ cd $GOPATH/src/github.com/patkaehuaea/command/timeserver
$ go install
$ $GOPATH/bin/timeserver

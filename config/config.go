//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015
//
// Wraps command line parsing and log initialization for timeserver and authserver.
// Flag parameters exposed as package exports. Defaults for all flags defined in this
// package.
package config

import (
	"flag"
	log "github.com/cihub/seelog"
	"os"
	"path/filepath"
	"time"
)

const (
	// Constants for timeserver:
	AUTH_HOST       = "localhost"
	AUTH_PORT       = ":9080"
	AUTH_TIMEOUT_MS = 1000 * time.Millisecond
	AVG_RESP_MS     = 1000 * time.Millisecond
	DEV_MS          = 100 * time.Millisecond
	MAX_IN_FLIGHT   = 0
	TIME_PORT       = ":8080"
	TMPL_DIR        = "templates"

	// Constants for authserver:
	DUMP_FILE      = ""
	CHECKPOINT_INT = 60 * time.Second

	// Constants for loadgen:
	LOAD_RATE       = 10
	LOAD_BURST      = 10
	LOAD_TIMEOUT_MS = 1500 * time.Millisecond
	LOAD_RUNTIME    = 10 * time.Second
	LOAD_URL        = "http://localhost:8080/time"

	// Constants shared accross applications:
	SEELOG_CONF_DIR  = "etc"
	SEELOG_CONF_FILE = "seelog.xml"
)

var (
	// Variables for timserver:
	AuthHost      *string
	AuthPort      *string
	AuthTimeoutMS *time.Duration
	AvgRespMS     *time.Duration
	DeviationMS   *time.Duration
	MaxInFlight   *int

	// Variables for authserver:
	DumpFile      *string
	CheckpointInt *time.Duration

	// Variables for loadgen:
	Rate          *int
	Burst         *int
	LoadTimeoutMS *time.Duration
	Runtime       *time.Duration
	URL           *string

	// Variables shared across applications:
	TimePort *string
	TmplDir  *string
	Verbose  *bool
	Logger   log.LoggerInterface
)

func init() {
	// Parameters for timeserver:
	AuthHost = flag.String("authhost", AUTH_HOST, "Hostname of downstream authentication server.")
	AuthTimeoutMS = flag.Duration("authtimeout-ms", AUTH_TIMEOUT_MS, "Milliseconds to wait before terminating downstream auth request.")
	AvgRespMS = flag.Duration("avg-response-ms", AVG_RESP_MS, "Average time to delay response to upstream time request.")
	DeviationMS = flag.Duration("deviation-ms", DEV_MS, "Average standard deviation in response delay to upstream time request.")
	MaxInFlight = flag.Int("max-inflight", MAX_IN_FLIGHT, "Maximum number of in-flight time requests the timeserver can handle.")
	TimePort = flag.String("port", TIME_PORT, "Time server binds to this port.")
	TmplDir = flag.String("templates", TMPL_DIR, "Directory relative to executable where templates are stored.")
	Verbose = flag.Bool("V", false, "Prints version number of program.")

	// Parameters for authserver:
	DumpFile = flag.String("dumpfile", DUMP_FILE, "Name of file storing state as JSON document.")
	CheckpointInt = flag.Duration("checkpoint-interval", CHECKPOINT_INT, "Dump state to file every checkpoint-interval seconds.")

	// Parameters for loadgen:
	Rate = flag.Int("rate", LOAD_RATE, "Average rate of requests (per second).")
	Burst = flag.Int("burst", LOAD_BURST, "Number of concurrent requests to issue.")
	LoadTimeoutMS = flag.Duration("timeout-ms", LOAD_TIMEOUT_MS, "Max time to wait for response.")
	Runtime = flag.Duration("runtime", LOAD_RUNTIME, "Number of seconds to process.")
	URL = flag.String("url", LOAD_URL, "URL to sample.")

	// Shared parameters:
	AuthPort = flag.String("authport", AUTH_PORT, "Auth server binds to this port.")

	// Local parameters:
	logConf := flag.String("log", SEELOG_CONF_FILE, "Name of log configuration file in etc directory relative to executable.")

	flag.Parse()

	// Will fail to default log configuration as defined by seelog package
	// if unable to open file. Assumes *LogConf is in SEELOG_CONF_DIR relative to cwd.
	cwd, _ := os.Getwd()
	Logger, _ = log.LoggerFromConfigAsFile(filepath.Join(cwd, SEELOG_CONF_DIR, *logConf))
}

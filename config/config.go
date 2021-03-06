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
	AUTH_HOST        = "localhost"
	AUTH_PORT        = ":9080"
	AUTH_TIMEOUT_MS  = 1000 * time.Millisecond
	AVG_RESP_MS      = 1000 * time.Millisecond
	CHECKPOINT_INT   = 60 * time.Second
	DEV_MS           = 100 * time.Millisecond
	DUMP_FILE        = ""
	MAX_IN_FLIGHT    = 0
	TIME_PORT        = ":8080"
	SEELOG_CONF_DIR  = "etc"
	SEELOG_CONF_FILE = "seelog.xml"
	TMPL_DIR         = "templates"
)

var (
	AuthHost      *string
	AuthPort      *string
	AuthTimeoutMS *time.Duration
	AvgRespMS     *time.Duration
	DeviationMS   *time.Duration
	DumpFile      *string
	CheckpointInt *time.Duration
	MaxInFlight   *int
	TimePort      *string
	TmplDir       *string
	Verbose       *bool
	Logger        log.LoggerInterface
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

	// Shared parameters:
	AuthPort = flag.String("authport", AUTH_PORT, "Auth server binds to this port.")

	// Local parameters:
	logConf := flag.String("log", SEELOG_CONF_FILE, "Name of log configuration file in etc directory relative to executable.")

	flag.Parse()

	// Will fail to default log configuration as defined by seelog package
	// if unable to open file. Assumes *LogConf is in SEELOG_CONF_DIR relative to cwd.
	cwd, _ := os.Getwd()
	var err error
	if Logger, err = log.LoggerFromConfigAsFile(filepath.Join(cwd, SEELOG_CONF_DIR, *logConf)); err != nil {
		log.Warn(err)
	}
}

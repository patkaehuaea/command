//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015
//
// Package contains an authentication server with backend users data store.
// The authserver reads configuration data from the config package and exposts
// two endpoints /get and /set. The former allows a caller to fetch the name of
// a user given a UUID, and the later allows setting a user in the data store
// given a UUID and name. For purposes of this assignment both endpoints are
// are implemented as HTTP GETs with data passed via query parameter.

package main

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/patkaehuaea/command/authserver/people"
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/counters"
	"io"
	"net/http"
	"os"
)

const (
	DEFAULT_DELTA      = 1
	GET_COOKIE_COUNTER = "get-cookie"
	KEY_200            = "200s"
	KEY_400            = "400s"
	KEY_500            = "500s"
	NO_COOKIE_COUNTER  = "no-cookie"
	SEELOG_CONF_DIR    = "etc"
	SEELOG_CONF_FILE   = "seelog.xml"
	SET_COOKIE_COUNTER = "set-cookie"
	VERSION_NUMBER     = "v0.0.1"
)

var (
	counter     *counters.Counter
	counterKeys = []string{
		GET_COOKIE_COUNTER,
		SET_COOKIE_COUNTER,
		NO_COOKIE_COUNTER,
		KEY_200,
		KEY_400,
		KEY_500,
	}
	users *people.UserStore
)

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	log.Info("authserver: Get user handler called.")

	if uuid := r.FormValue("cookie"); people.IsValidUUID(uuid) {
		log.Debug("authserver: Found valid uuid: " + uuid)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, users.Name(uuid))
	} else {
		log.Debug("authserver: UUID not valid, or not found in users.")
		w.WriteHeader(http.StatusBadRequest)
		counter.Increment(NO_COOKIE_COUNTER, DEFAULT_DELTA)
	}
	counter.Increment(KEY_200, DEFAULT_DELTA)
	counter.Increment(GET_COOKIE_COUNTER, DEFAULT_DELTA)
}

func handleMonitor(w http.ResponseWriter, r *http.Request) {
	log.Info("authserver: Monitor called.")
	counter.Increment(KEY_200, DEFAULT_DELTA)
	copy := counter.Copy()
	data, err := json.Marshal(&copy)
	if err != nil {
		counter.Increment(KEY_500, DEFAULT_DELTA)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func handleSetUser(w http.ResponseWriter, r *http.Request) {
	log.Info("authserver: Set user handler called.")

	uuid := r.FormValue("cookie")
	name := r.FormValue("name")

	if people.IsValidUUID(uuid) && people.IsValidName(name) {
		users.Add(uuid, name)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Debug("authserver: Invalid uuid and/or name.")
		w.WriteHeader(http.StatusBadRequest)
	}
	counter.Increment(KEY_200, DEFAULT_DELTA)
	counter.Increment(SET_COOKIE_COUNTER, DEFAULT_DELTA)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Info("authserver: Not found handler called.")
	w.WriteHeader(http.StatusNotFound)
	counter.Increment(KEY_400, DEFAULT_DELTA)
}

func init() {

	log.ReplaceLogger(config.Logger)
	counter = counters.New(counterKeys)

	// DumpFile needs to be specified, but dumpfile need
	// not be present at startup.
	if *config.DumpFile == config.DUMP_FILE {
		log.Critical("database: Dumpfile not specified.")
		os.Exit(1)
	}

	// Initialization of the backend data store should be
	// transparent to the authserver. Future project to move
	// into its own pacakge's init() function and have authserver
	// reference a public member.
	users = people.NewUsers()
	if err := users.Load(*config.DumpFile); err != nil {
		log.Info("database: Backup not found at initialization.")
	}
	go users.Persist(*config.DumpFile, *config.CheckpointInt)
}

func main() {

	/*
	   Paramters surfaced via config pacakge used in this program:
	   *config.AuthPort
	   config.Logger
	   database.Users
	*/

	r := mux.NewRouter()
	r.HandleFunc("/get", handleGetUser).Methods("GET")
	r.HandleFunc("/set", handleSetUser).Methods("GET")
	r.HandleFunc("/monitor", handleMonitor)
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)
	http.Handle("/", r)
	if err := (http.ListenAndServe(*config.AuthPort, nil)); err != nil {
		log.Critical(err)
	}
}

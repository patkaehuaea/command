//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February2015

package main

import (
	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/timeserver/people"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	VERSION_NUMBER   = "v0.0.1"
	SEELOG_CONF_DIR  = "etc"
	SEELOG_CONF_FILE = "seelog.xml"
)

var users = people.NewUsers()

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	log.Info("Get user handler called.")

	if uuid := r.Form.Get("cookie"); people.IsValidUUID(uuid) {
		log.Debug("Found valid uuid: " + uuid)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, users.Name(uuid))
	} else {
		log.Debug("UUID not valid, or not found in users.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handleSetUser(w http.ResponseWriter, r *http.Request) {
	log.Info("Set user handler called.")

	person, err := people.NewPerson(r.Form.Get("cookie"), r.Form.Get("name"))
	// TODO: Check if user already exists with key.
	if err == nil {
		users.Add(person)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	users.DumpFile()
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Info("Not found handler called.")
	w.WriteHeader(http.StatusNotFound)
}

func parseFormWrapper(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err == nil {
			fn(w, r)
		} else {
			log.Error(err)
		}
		// TODO: Maybe return 500 here if unable to ParseForm
	}
}

func main() {

	/*
	   Paramters surfaced via config pacakge used in this program:
	   *config.AuthPort
	   *config.CheckpointInt
	   *config.DumpFile
	*/

	// Server will fail to default log configuration as defined by seelog package
	// if unable to open file. Assumes *logConf is in SEELOG_CONF_DIR relative to cwd.
	cwd, _ := os.Getwd()
	logger, err := log.LoggerFromConfigAsFile(filepath.Join(cwd, SEELOG_CONF_DIR, *config.LogConf))
	if err != nil {
		log.Error(err)
	}
	log.ReplaceLogger(logger)

	r := mux.NewRouter()
	r.HandleFunc("/get", parseFormWrapper(handleGetUser)).Methods("GET")
	r.HandleFunc("/set", parseFormWrapper(handleSetUser)).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)
	http.Handle("/", r)
	if err := (http.ListenAndServe(*config.AuthPort, nil)); err != nil {
		log.Critical(err)
	}
}

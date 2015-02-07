//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February2015

package main

import (
	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/patkaehuaea/command/authserver/people"
	"github.com/patkaehuaea/command/config"
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
	log.Info("authserver: Get user handler called.")

	if uuid := r.FormValue("cookie"); people.IsValidUUID(uuid) {
		log.Debug("authserver: Found valid uuid: " + uuid)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, users.Name(uuid))
	} else {
		log.Debug("authserver: UUID not valid, or not found in users.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handleSetUser(w http.ResponseWriter, r *http.Request) {
	log.Info("authserver: Set user handler called.")

	person, err := people.NewPerson(r.FormValue("cookie"), r.FormValue("name"))
	if err == nil {
		users.Add(person)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Info("authserver: Not found handler called.")
	w.WriteHeader(http.StatusNotFound)
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
	r.HandleFunc("/get", handleGetUser).Methods("GET")
	// Should be POST, but assignment spec requires GET.
	r.HandleFunc("/set", handleSetUser).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)
	http.Handle("/", r)
	if err := (http.ListenAndServe(*config.AuthPort, nil)); err != nil {
		log.Critical(err)
	}
}

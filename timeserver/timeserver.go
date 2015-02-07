//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, January 2015
//
// Package contains simple web server that binds to specified --port or 8080.
// Exectuable accepts two parameters, --port to designate listen port,
// and -V to output the version number of the program. Additional flag --templates
// determines location of templates on filesystem and --log parameter provides
// name of seelog configuration file in etc/. Server provieds '/time' endpoint as
// well as '/login' '/logout' and root pages '/', 'index.html'. State is lost
// upon program termination.
package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/timeserver/auth"
	"github.com/patkaehuaea/command/timeserver/cookie"
	"github.com/patkaehuaea/command/timeserver/people"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	VERSION_NUMBER       = "v1.3.0"
	SEELOG_CONF_DIR      = "etc"
	SEELOG_CONF_FILE     = "seelog.xml"
	TEMPL_DIR            = "templates"
	TEMPL_FILE_EXTENSION = ".tmpl"
	LOCAL_TIME_LAYOUT    = "3:04:05 PM"
	UTC_TIME_LAYOUT      = "15:04:05 UTC"
)

var (
	templates *template.Template
	authClient *auth.AuthClient
)

func getUUIDThenName(r *http.Request) (name string, err error) {
	log.Info("Called getUUIDThenname function.")
	uuid, uuidErr := cookie.UUID(r)
	if uuidErr != nil {
		log.Warn(uuidErr)
		err = uuidErr
		return
	}
	response, authErr := authClient.Get(uuid)
	if authErr != nil {
		log.Warn(authErr)
		err = authErr
		return
	}
	name = response
	return
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	log.Info("Default handler called.")
	if name, err := getUUIDThenName(r); err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		log.Debug(name + " viewing site.")
		renderTemplate(w, "greetings", name)
	}
}

func handleDisplayLogin(w http.ResponseWriter, r *http.Request) {
	log.Info("Display login handler called.")
	renderTemplate(w, "login", "What is your name, Earthling?")
}

func handleProcessLogin(w http.ResponseWriter, r *http.Request) {
	log.Info("Process login handler called.")
	name := r.FormValue("name")

	if valid := people.IsValidName(name); valid {
		log.Debug("Name matched regex.")
		uuid := people.UUID()
		if err := authClient.Set(uuid, name); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			// Wouldn't typicall print hacker friendly error
			// messages to webpage but useful for this project.
			renderTemplate(w, "500", err.Error())
		} else {
			http.SetCookie(w, cookie.NewCookie(uuid, cookie.MAX_AGE))
			http.Redirect(w, r, "/", http.StatusFound)
			log.Info(name + " registered on site.")
		}
	} else {
		log.Debug("Invalid username or registration failed.")
		w.WriteHeader(http.StatusBadRequest)
		renderTemplate(w, "login", "C'mon, I need a name.")
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	log.Info("Logout handler called.")
	http.SetCookie(w, cookie.NewCookie("deleted", cookie.DELETE_AGE))
	renderTemplate(w, "logged-out", nil)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Info("Not found handler called.")
	w.WriteHeader(http.StatusNotFound)
	renderTemplate(w, "404", nil)
}

func handleTime(w http.ResponseWriter, r *http.Request) {
	log.Info("Time handler called.")
	name, _ := getUUIDThenName(r)
	// If name is blank, template will not render
	// personalized greeting.
	params := map[string]interface{}{
		"localTime": time.Now().Format(LOCAL_TIME_LAYOUT),
		"UTCTime":   time.Now().Format(UTC_TIME_LAYOUT),
		"name":      name,
	}
	renderTemplate(w, "time", params)
}

// credit: http://tinyurl.com/kwc4hls
func logFileRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("File server called.")
		h.ServeHTTP(w, r)
	})
}

// credit: https://golang.org/doc/articles/wiki/#tmp_10
func renderTemplate(w http.ResponseWriter, templ string, d interface{}) {
	err := templates.ExecuteTemplate(w, templ+TEMPL_FILE_EXTENSION, d)
	if err != nil {
		log.Error("Error looking for template: " + templ)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	/*
		Paramters surfaced via config pacakge used in this program:
		*config.AuthHost
		*config.AuthPort
		*config.AuthTimeoutMS
		*config.AvgRespMS
		*config.DeviationMS
		*config.LogConf
		*config.MaxInFlight
		*config.TimePort
		*config.TmplDir
		*config.Verbose
	*/

	if *config.Verbose {
		fmt.Printf("Version number: %s \n", VERSION_NUMBER)
		os.Exit(1)
	}

	// Restrict parsing to *.templ to prevent fail on non-template files in a given directory
	// like .DS_STORE.
	var err error
	templates, err = template.ParseGlob(filepath.Join(*config.TmplDir, "*"+TEMPL_FILE_EXTENSION))
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}

	authClient = auth.NewAuthClient(*config.AuthHost, *config.AuthPort, *config.AuthTimeoutMS)

	// Server will fail to default log configuration as defined by seelog package
	// if unable to open file. Assumes *logConf is in SEELOG_CONF_DIR relative to cwd.
	cwd, _ := os.Getwd()
	logger, err := log.LoggerFromConfigAsFile(filepath.Join(cwd, SEELOG_CONF_DIR, *config.LogConf))
	if err != nil {
		log.Error(err)
	}
	log.ReplaceLogger(logger)

	r := mux.NewRouter()
	r.HandleFunc("/", handleDefault)
	r.PathPrefix("/css/").Handler(logFileRequest(http.StripPrefix("/css/", http.FileServer(http.Dir("css/")))))
	r.HandleFunc("/index.html", handleDefault)
	r.HandleFunc("/login", handleDisplayLogin).Methods("GET")
	r.HandleFunc("/login", handleProcessLogin).Methods("POST")
	r.HandleFunc("/logout", handleLogout)
	r.HandleFunc("/time", handleTime)
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)
	http.Handle("/", r)
	if err := (http.ListenAndServe(*config.TimePort, nil)); err != nil {
		log.Critical(err)
	}
}

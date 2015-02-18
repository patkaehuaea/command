//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015
//
// Package contains simple web server that provides '/time' endpoint as
// well as '/login', '/logout', '/', and 'index.html'. Operations to
// find a user given a UUID, and create a user are conducted via the
// client package that abstracts HTTP communication with authserver from
// this program. Configuration data for btoh timeserver and authserver
// are exposed in the config pacakge. Timeserver will only throttle
// requests to the time endpoint and stats are only available if requests
// are configured to throttle. The latter two features aren't derived from
// customer use case's but from the desire to simulate load on timeserver.
package main

import (
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/patkaehuaea/command/authserver/client"
	"github.com/patkaehuaea/command/authserver/people"
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/timeserver/cookie"
	"github.com/patkaehuaea/command/timeserver/stats"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	VERSION_NUMBER       = "v2.3.2"
	TEMPL_DIR            = "templates"
	TEMPL_FILE_EXTENSION = ".tmpl"
	LOCAL_TIME_LAYOUT    = "3:04:05 PM"
	UTC_TIME_LAYOUT      = "15:04:05 UTC"
)

var (
	authClient *client.AuthClient
	inFlight   *stats.ConcurrentRequests
	templates  *template.Template
)

// Credit: http://goo.gl/MsxPHk
func delay(average time.Duration, deviation time.Duration) {
	log.Trace("timeserver: delay average - " + average.String() + " ; " + "delay deviation = " + deviation.String())
	load := time.Duration(rand.NormFloat64())*deviation + average
	log.Debug("timeserver: Sleeping for " + load.String() + ".")
	time.Sleep(load)
}

func getUUIDThenName(r *http.Request) (name string, err error) {
	log.Info("timeserver: Called getUUIDThenName function.")

	var uuid string
	if uuid, err = cookie.UUID(r); err != nil {
		log.Warn(err)
		return
	}

	if name, err = authClient.Get(uuid); err != nil {
		log.Warn(err)
		return
	}

	// Prevents issues where cookies persists in browser but
	// does not persist in authserver. Caller should be notified
	// that authserver contains empty result.
	if name == "" {
		err = errors.New("timeserver: Empty result from get user.")
		log.Warn(err)
	}

	return
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	log.Info("timeserver: Default handler called.")

	name, err := getUUIDThenName(r)

	if err != nil {
		http.SetCookie(w, cookie.NewCookie(cookie.DELETE_VALUE, cookie.DELETE_AGE))
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	log.Debug("timeserver: " + name + " viewing site.")
	renderTemplate(w, "greetings", name)
}

func handleDisplayLogin(w http.ResponseWriter, r *http.Request) {
	log.Info("timeserver: Display login handler called.")
	renderTemplate(w, "login", "What is your name, Earthling?")
}

func handleProcessLogin(w http.ResponseWriter, r *http.Request) {
	log.Info("timeserver: Process login handler called.")

	name := r.FormValue("name")

	if people.IsValidName(name) {
		log.Trace("timeserver: Name matched regex.")
		uuid := people.UUID()

		if err := authClient.Set(uuid, name); err != nil {
			http.SetCookie(w, cookie.NewCookie(cookie.DELETE_VALUE, cookie.DELETE_AGE))
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate(w, "500", nil)
			log.Error(err)
			return
		}

		http.SetCookie(w, cookie.NewCookie(uuid, cookie.MAX_AGE))
		http.Redirect(w, r, "/", http.StatusFound)
		log.Info("timeserver: " + name + " registered on site.")
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	renderTemplate(w, "login", "C'mon, I need a name.")
	log.Warn("timeserver: Invalid username or registration failed.")
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	log.Info("timeserver: Logout handler called.")

	http.SetCookie(w, cookie.NewCookie(cookie.DELETE_VALUE, cookie.DELETE_AGE))
	renderTemplate(w, "logged-out", nil)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Info("timeserver: Not found handler called.")

	w.WriteHeader(http.StatusNotFound)
	renderTemplate(w, "404", nil)
}

func handleTime(w http.ResponseWriter, r *http.Request) {
	log.Info("timeserver: Time handler called.")

	// Simulate load with delay function.
	delay(*config.AvgRespMS, *config.DeviationMS)

	name, err := getUUIDThenName(r)

	if err != nil {
		http.SetCookie(w, cookie.NewCookie(cookie.DELETE_VALUE, cookie.DELETE_AGE))
	}

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
		log.Info("timeserver: File server called.")
		h.ServeHTTP(w, r)
	})
}

// credit: https://golang.org/doc/articles/wiki/#tmp_10
func renderTemplate(w http.ResponseWriter, templ string, d interface{}) {
	err := templates.ExecuteTemplate(w, templ+TEMPL_FILE_EXTENSION, d)
	if err != nil {
		log.Error("timeserver: Error looking for template: " + templ)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func throttle(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := inFlight.Add(); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate(w, "500", nil)
			return
		}

		fn(w, r)
		// Only subtract if stat was incrememted otherwise
		// may attempt to subtract below stats.MIN_VALUE.
		if err := inFlight.Subtract(); err != nil {
			log.Error(err)
		}
	}
}

func init() {

	// Restrict parsing to *.templ to prevent fail on non-template files in a given directory
	// like .DS_STORE.
	var err error
	if templates, err = template.ParseGlob(filepath.Join(*config.TmplDir, "*"+TEMPL_FILE_EXTENSION)); err != nil {
		log.Critical(err)
		os.Exit(1)
	}

	log.ReplaceLogger(config.Logger)
	authClient = client.NewAuthClient(*config.AuthHost, *config.AuthPort, *config.AuthTimeoutMS)
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
		config.Logger
		*config.MaxInFlight
		*config.TimePort
		*config.TmplDir
		*config.Verbose
	*/

	if *config.Verbose {
		fmt.Printf("Version number: %s \n", VERSION_NUMBER)
		os.Exit(0)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handleDefault)
	r.PathPrefix("/css/").Handler(logFileRequest(http.StripPrefix("/css/", http.FileServer(http.Dir("css/")))))
	r.HandleFunc("/index.html", handleDefault)
	r.HandleFunc("/login", handleDisplayLogin).Methods("GET")
	r.HandleFunc("/login", handleProcessLogin).Methods("POST")
	r.HandleFunc("/logout", handleLogout)
	if *config.MaxInFlight != 0 {
		log.Infof("%s - %d", "timeserver: Max concurrent time connections", *config.MaxInFlight)
		inFlight = stats.NewCR(*config.MaxInFlight)
		r.HandleFunc("/time", throttle(handleTime))
	}
	r.HandleFunc("/time", handleTime)
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)
	http.Handle("/", r)
	if err := (http.ListenAndServe(*config.TimePort, nil)); err != nil {
		log.Critical(err)
		os.Exit(1)
	}
}

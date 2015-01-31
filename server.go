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
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/patkaehuaea/server/cookie"
	"github.com/patkaehuaea/server/people"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	VERSION_NUMBER       = "v1.2.4"
	SERVER_PORT          = ":8080"
	SEELOG_CONF_DIR      = "etc"
	SEELOG_CONF_FILE     = "seelog.xml"
	TEMPL_DIR            = "templates"
	TEMPL_FILE_EXTENSION = ".tmpl"
	LOCAL_TIME_LAYOUT    = "3:04:05 PM"
	UTC_TIME_LAYOUT      = "15:04:05 UTC"
)

var (
	templates *template.Template
	users     = people.NewUsers()
)

func handleDefault(w http.ResponseWriter, r *http.Request, name string) {
	log.Info("Default handler called.")
	if name != "" {
		log.Debug("User: " + name + " viewing site.")
		renderTemplate(w, "greetings", name)
	} else {
		log.Debug("Name is empty. Redirecting to login.")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func handleDisplayLogin(w http.ResponseWriter, r *http.Request) {
	log.Info("Display login handler called.")
	renderTemplate(w, "login", "What is your name, Earthling?")
}

func handleProcessLogin(w http.ResponseWriter, r *http.Request) {
	log.Info("Process login handler called.")
	name := r.FormValue("name")
	if valid, _ := people.FirstAndOrLastName(name); valid {
		log.Debug("Name matched regex.")
		person := people.NewPerson(name)
		users.Add(person)
		http.SetCookie(w, cookie.NewCookie(person.ID, cookie.MAX_AGE))
		http.Redirect(w, r, "/", http.StatusFound)
		log.Debug("User: " + person.Name + " logged in.")
	} else {
		log.Debug("Invalid username. Rendering login page.")
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

func handleTime(w http.ResponseWriter, r *http.Request, name string) {
	log.Info("Time handler called.")
	// Template will not render personal greeting if name blank.
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

// Read 'uuid' cookie, then perform lookup in users before passing name
// to calling handler. Centralizes cookie parsing and data lookup. Can
// extend to make remote call and/or pass additional data to callers.
func uuidCookieToName(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := ""
		if uuid, _ := cookie.UUIDCookieValue(r); uuid != "" {
			name = users.Name(uuid)
		}
		log.Debug("Cookie value not found, or user not found.")
		fn(w, r, name)
	}
}

func main() {
	logConf := flag.String("log", SEELOG_CONF_FILE, "Name of log configuration file in etc directory relative to executable.")
	port := flag.String("port", SERVER_PORT, "Web server binds to this port. Default is 8080.")
	templDir := flag.String("templates", TEMPL_DIR, "Directory relative to executable where templates are stored.")
	verbose := flag.Bool("V", false, "Prints version number of program.")
	flag.Parse()

	if *verbose {
		fmt.Printf("Version number: %s \n", VERSION_NUMBER)
		os.Exit(1)
	}

	// Restrict parsing to *.templ to prevent fail on non-template files in a given directory
	// like .DS_STORE.
	var err error
	templates, err = template.ParseGlob(filepath.Join(*templDir, "*"+TEMPL_FILE_EXTENSION))
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}

	// Server will fail to default log configuration as defined by seelog package
	// if unable to open file. Assumes *logConf is in SEELOG_CONF_DIR relative to cwd.
	cwd, _ := os.Getwd()
	logger, err := log.LoggerFromConfigAsFile(filepath.Join(cwd, SEELOG_CONF_DIR, *logConf))
	if err != nil {
		log.Error(err)
	}
	log.ReplaceLogger(logger)

	r := mux.NewRouter()
	r.HandleFunc("/", uuidCookieToName(handleDefault))
	r.PathPrefix("/css/").Handler(logFileRequest(http.StripPrefix("/css/", http.FileServer(http.Dir("css/")))))
	r.HandleFunc("/index.html", uuidCookieToName(handleDefault))
	r.HandleFunc("/login", handleDisplayLogin).Methods("GET")
	r.HandleFunc("/login", handleProcessLogin).Methods("POST")
	r.HandleFunc("/logout", handleLogout)
	r.HandleFunc("/time", uuidCookieToName(handleTime))
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)
	http.Handle("/", r)
	if err := (http.ListenAndServe(*port, nil)); err != nil {
		log.Critical(err)
	}
}

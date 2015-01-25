//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, January 2015
//
// Package contains simple web server that binds to port 8080. Exectuable accepts
// two parameters, --port to designate listen port, and -V to output the version
// number of the program. Server provieds '/time' endpoint as well as '/login' '/logout'
// and root pages '/', 'index.html'. Pages are rendered from templates that must be
// located in a 'templates/' directory relative to the executable. This package uses
// adjacent people package to maintain state as it relates to visits. State is lost
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
	"regexp"
	"time"
)

const (
	VERSION_NUMBER       = "v1.2.2"
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

func handleDefault(w http.ResponseWriter, r *http.Request) {
	log.Info("Default handler called.")
	id, _ := cookie.UUIDValue(r)
	if name := users.Name(id); name != "" {
		log.Debug("User: " + name + " viewing site.")
		renderTemplate(w, "greetings", name)
	} else {
		log.Debug("No cookie found or value empty. Redirecting to login.")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// Handling GET and POST methods can be implemented on separate /login
// handlers with mux. Left as-is for clarity of flow.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Info("Login handler called.")
	if r.Method == "GET" {
		log.Debug("Login GET method detected.")
		renderTemplate(w, "login", "What is your name, Earthling?")
	} else if r.Method == "POST" {
		log.Debug("Login POST method detected.")
		name := r.FormValue("name")
		// Allows first name, or first and last name in English characters with intervening space.
		// Minimum length of name is two characters and maximum length of field is 71 characters
		// including space.
		if valid, _ := regexp.MatchString("^[a-zA-Z]{2,35} {0,1}[a-zA-Z]{0,35}$", name); valid {
			log.Debug("Name matched regex.")
			person := people.NewPerson(name)
			users.Add(person)
			http.SetCookie(w, cookie.NewCookie(person.ID, cookie.MAX_AGE))
			http.Redirect(w, r, "/", http.StatusFound)
			log.Debug("User: " + person.Name + " logged in.")
			return
		} else {
			log.Debug("Invalid username. Rendering login page.")
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate(w, "login", "C'mon, I need a name.")
		}
	} else {
		log.Debug("Login request method not handled.")
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	log.Info("Logout handler called.")
	// Invalidate data along and set MaxAge to avoid accidental persistence issues.
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
	id, _ := cookie.UUIDValue(r)
	// Personalized message will only display if user's cookie contains an id
	// and that id is found in the users table. Template handles display logic.
	params := map[string]interface{}{"localTime": time.Now().Format(LOCAL_TIME_LAYOUT),
		"UTCTime": time.Now().Format(UTC_TIME_LAYOUT),
		"name":    users.Name(id)}
	renderTemplate(w, "time", params)
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
	logConf := flag.String("log", SEELOG_CONF_FILE, "Name of log configuration file in etc directory relative to executable.")
	port := flag.String("port", SERVER_PORT, "Web server binds to this port. Default is 8080.")
	templDir := flag.String("templates", TEMPL_DIR, "Directory relative to executable where templates are stored.")
	verbose := flag.Bool("V", false, "Prints version number of program.")
	flag.Parse()

	if *verbose {
		fmt.Printf("Version number: %s \n", VERSION_NUMBER)
		os.Exit(1)
	}

	cwd, _ := os.Getwd()

	// TODO: Handle error when no *.templ files exist in directory and when bad format.
	// Restricting parsing to *.templ prevents server from attempting to parse irrelevant files
	// in a given directory like .DS_STORE.
	templates = template.Must(template.ParseGlob(filepath.Join(cwd, *templDir, "*"+TEMPL_FILE_EXTENSION)))

	// TODO: Handle errors opening log file.
	// Server will fail to default log configuration as defined by seelog package
	// if unable to open file.
	logger, err := log.LoggerFromConfigAsFile(filepath.Join(cwd, SEELOG_CONF_DIR, *logConf))
	if err != nil {
		log.Critical(err)
	}
	log.ReplaceLogger(logger)

	r := mux.NewRouter()
	r.HandleFunc("/", handleDefault)
	// TODO: Log file system css access.
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	r.HandleFunc("/time", handleTime)
	r.HandleFunc("/index.html", handleDefault)
	r.HandleFunc("/login", handleLogin)
	r.HandleFunc("/logout", handleLogout)
	r.HandleFunc("/time", handleTime)
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)
	http.Handle("/", r)
	log.Critical(http.ListenAndServe(*port, nil))
}

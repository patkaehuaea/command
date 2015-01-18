//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, January 2015

// Package contains simple web server that binds to port 8080. Exectuable accepts
// two parameters, --port to designate listen port, and -V to output the version
// number of the program. Server responsds to only one request at /time and responds
// with the current time. All other requests should generate an http 404 status and
// custom not found page.

// TODO: Extend Cookie Struct
// TODO: Determine if logout needs to have user removed from map.

package main

import (
	// "errors"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/patkaehuaea/people"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const timeLayout = "3:04:05 PM"
const COOKIE_NAME = "uuid"
const COOKIE_MAX_AGE = 86400

// Program expects html template directory to be in same path as executable is run.
var cwd, _ = os.Getwd()
var templates = template.Must(template.ParseGlob(filepath.Join(cwd, "templates", "*.html")))
var users = people.NewUsers()

func handleDefault(w http.ResponseWriter, r *http.Request) {
	debug("Default handler called.", r)
	id, _ := idFromUUIDCookie(r)
	if name := users.Name(id) ; name != "" {
		info("ID found in users table.", r)
		renderTemplate(w, "greetings", name)
	} else {
		debug("No cookie found or value empty. Redirecting to login.", r)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	debug("Login handler called.", r)
	if r.Method == "GET" {
		debug("Login GET method detected.", r)
		renderTemplate(w, "login", nil)
	} else if r.Method == "POST" {
		debug("Login POST method detected.", r)
		name := r.FormValue("name")
		// Allows first name, or first and last name in English characters with intervening space. 
		// Minimum length of name is two characters and maximum length of field is 71 characters 
		// including space. Case where field is completely empty handled by javascript in template.
		if valid, _ := regexp.MatchString("^[a-zA-Z]{2,35} {0,1}[a-zA-Z]{0,35}$", name) ; valid {
			debug("Name matched regex.", r)
			// uuid := uuid()
			person := people.NewPerson(name)
			users.Add(person)
			setCookie(w, person.ID, COOKIE_MAX_AGE)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		} else {
			// Fail and require new input rather than cleaning and 
			// passing on.
			debug("Invalid username. Redirecting to root.", r)
			//w.WriteHeader(http.StatusBadRequest)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		debug("Login request method not handled.", r)
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	debug("Logout handler called.", r)
	// Invalidate data along and set MaxAge to avoid accidental persistence issues.
	setCookie(w, "deleted", -1)
	renderTemplate(w, "logged-out", nil)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	debug("Not found handler called.", r)
	w.WriteHeader(http.StatusNotFound)
	renderTemplate(w, "404", nil)
}

func handleTime(w http.ResponseWriter, r *http.Request) {
	debug("Time handler called.", r)
	id, _ := idFromUUIDCookie(r)
	params := map[string]interface{}{"time": time.Now().Format(timeLayout), "name": users.Name(id)}
	renderTemplate(w, "time", params)
}

// Possible to define customer formatter to avoid having
// to call WithFields on each log call. In the meantime
// implementing helpers function as is.
func info(msg string, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"time":   time.Now().Format(timeLayout),
		"url":    r.URL,
	}).Info(msg)
}

func debug(msg string, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"time":   time.Now().Format(timeLayout),
		"url":    r.URL,
	}).Debug(msg)
}

// credit: https://golang.org/doc/articles/wiki/#tmp_10
func renderTemplate(w http.ResponseWriter, templ string, d interface{}) {
	err := templates.ExecuteTemplate(w, templ + ".html", d)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func setCookie(w http.ResponseWriter, uuid string, maxAge int) {
	c := http.Cookie{Name: COOKIE_NAME, Value: uuid, Path: "/", MaxAge: maxAge}
	http.SetCookie(w, &c)
}

func idFromUUIDCookie(r *http.Request) (string, error) {
	log.Debug("Reading cookie 'uuid'.")
	cookie, err := r.Cookie(COOKIE_NAME)
	if err == http.ErrNoCookie {
		log.Debug("Cookie not found.")
		return "", http.ErrNoCookie
	}
	return cookie.Value, nil
}

func main() {

	const VERSION_NUMBER = "v1.0.6"

	portPtr := flag.String("port", "8080", "Web server binds to this port. Default is 8080.")
	verbosePtr := flag.Bool("V", false, "Prints version number of program.")
	flag.Parse()
	portParam := ":" + *portPtr

	if *verbosePtr {
		fmt.Printf("Version number: %s \n", VERSION_NUMBER)
		os.Exit(1)
	}

	log.SetLevel(log.DebugLevel)

	// The gorilla web toolkit (http://www.gorillatoolkit.org/) seems like it provides a cleaner way
	// to handle notFound and provides some additional functionality.
	r := mux.NewRouter()
	r.HandleFunc("/", handleDefault)
	r.HandleFunc("/index.html", handleDefault)
	r.HandleFunc("/login", handleLogin)
	r.HandleFunc("/logout", handleLogout)
	r.HandleFunc("/time", handleTime)
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(portParam, nil))
}

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
	"errors"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const timeLayout = "3:04:05 PM"
const COOKIE_NAME = "uuid"
const COOKIE_MAX_AGE = 86400

// Program expects html template directory to be in same path as executable is run.
var cwd, _ = os.Getwd()
var templates = template.Must(template.ParseGlob(filepath.Join(cwd, "templates", "*.html")))

// credit: https://blog.golang.org/go-maps-in-action
var users = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// credit: https://blog.golang.org/go-maps-in-action
func addName(uuid string, name string) bool {
	if validateName(name) {
		users.Lock()
		users.m[uuid] = name
		users.Unlock()
		return true
	}
	return false
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	logInfo("Default handler called.", r)

	name, err := uuidCookieToName(r)
	if name == "" || err != nil {
		log.Debug("No cookie found or value empty. Redirecting to login.")
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	log.Debug("Cookie uuid found in user table: " + name)
	templates.ExecuteTemplate(w, "greetings", name)
}

// credit: https://blog.golang.org/go-maps-in-action
func findName(uuid string) string {
	users.RLock()
	name := users.m[uuid]
	users.RUnlock()
	return name
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Login handler called.")

	if r.Method == "GET" {
		log.Debug("Login GET method detected.")
		templates.ExecuteTemplate(w, "login", nil)
		log.Debug("Login template rendered.")
	}

	if r.Method == "POST" {
		log.Debug("Login POST method detected.")
		// Form will not submit if name empty.
		name := r.FormValue("name")
		if validateName(name) {
			uuid := uuid()
			addName(uuid, name)
			setCookie(w, uuid, COOKIE_MAX_AGE)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		} else {
			// Redirect user with 4xx status code.
			log.Debug("Invalid username. Redirecting to root.")
			w.WriteHeader(http.StatusBadRequest)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	log.Debug("Request method not handled.")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Logout handler called.")
	// Invalidate data along and set MaxAge to avoid accidental persistence issues.
	setCookie(w, "deleted", -1)
	templates.ExecuteTemplate(w, "logged-out.html", nil)
}

func logInfo(msg string, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"time":   time.Now().Format(localFormat),
		"url":    r.URL,
	}).Info(msg)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	logInfo("Not found handler called.", r)
	w.WriteHeader(http.StatusNotFound)
	templates.ExecuteTemplate(w, "404.html", nil)
}

func setCookie(w http.ResponseWriter, uuid string, maxAge int) {
	c := http.Cookie{Name: COOKIE_NAME, Value: uuid, Path: "/", MaxAge: maxAge}
	http.SetCookie(w, &c)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "time.html", time.Now().Format(timeLayout))
}

func uuid() string {
	// credit: http://golang.org/pkg/os/exec/#Cmd.Run
	log.Debug("Getting uuid.")
	out, err := exec.Command("/usr/bin/uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	uuid := strings.TrimSuffix(string(out), "\n")

	log.Debug("Clean uuid generated:" + uuid)
	return uuid
}

func uuidCookieToName(r *http.Request) (uName string, err error) {
	log.Debug("Reading cookie 'uuid' and finding name.")

	cookie, err := r.Cookie(COOKIE_NAME)
	// TODO: Implement additional cookie validation
	// like domain and expiry in own method.
	if err == http.ErrNoCookie {
		return "", http.ErrNoCookie
	}

	name := findName(cookie.Value)

	if name == "" {
		return "", errors.New("Cookie value not found in user table.")
	}

	return name, nil
}

func validateName(name string) bool {
	// TODO: better to implement in form?
	// TODO: implement name validation
	return true
}

func main() {

	const VERSION_NUMBER = "v1.0.4"

	portPtr := flag.String("port", "8080", "Web server binds to this port. Default is 8080.")
	verbosePtr := flag.Bool("V", false, "Prints version number of program.")
	flag.Parse()
	portParam := ":" + *portPtr

	if *verbosePtr {
		fmt.Printf("Version number: %s \n", VERSION_NUMBER)
		os.Exit(1)
	}

	// The gorilla web toolkit (http://www.gorillatoolkit.org/) seems like it provides a cleaner way
	// to handle notFound and provides some additional functionality.
	r := mux.NewRouter()
	go r.HandleFunc("/", defaultHandler)
	go r.HandleFunc("/index.html", defaultHandler)
	go r.HandleFunc("/login", loginHandler)
	go r.HandleFunc("/logout", logoutHandler)
	go r.HandleFunc("/time", timeHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(portParam, nil))
}

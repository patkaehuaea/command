//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, January 2015

// Package contains simple web server that binds to port 8080. Exectuable accepts
// two parameters, --port to designate listen port, and -V to output the version
// number of the program. Server responsds to only one request at /time and responds
// with the current time. All other requests should generate an http 404 status and
// custom not found page.

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Program expects html template directory to be in same path as executable is run.
var cwd, _ = os.Getwd()
var templates = template.Must(template.ParseGlob(filepath.Join(cwd, "templates", "*.html")))

const timeLayout = "3:04:05 PM"

func timeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "time.html", time.Now().Format(timeLayout))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.ExecuteTemplate(w, "http404.html", nil)
}

func main() {

	const VERSION_NUMBER = "v1.0.0"

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
	r.HandleFunc("/time", timeHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/", r)
	http.ListenAndServe(portParam, nil)
}

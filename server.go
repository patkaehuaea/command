/* Copyright (C) Pat Kaehuaea - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Pat Kaehuaea, January 2015
 */

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

//credit: http://stackoverflow.com/questions/17206467/go-how-to-render-multiple-templates-in-golang
var templates = template.Must(template.ParseGlob(htmlTemplPath()))

func htmlTemplPath() string {
	curDir, _ := os.Getwd()
	templatesPath := filepath.Join(curDir, "templates", "*.html")
	return templatesPath
}

func getTime() string {
	const layout = "3:04:05 PM"
	t := time.Now().Format(layout)
	return t
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "time", getTime())
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.ExecuteTemplate(w, "http404", nil)
}

func main() {

	const VERSION_NUMBER = "0.0.1"

	portPtr := flag.String("port", "8080", "Web server binds to this port. Default is 8080.")
	verbosePtr := flag.Bool("V", false, "Prints version number of program.")
	portString := ":" + *portPtr
	flag.Parse()

	if *verbosePtr {
		fmt.Printf("Version number: %s \n", VERSION_NUMBER)
		os.Exit(1)
	}

	//credit: http://stackoverflow.com/questions/9996767/showing-custom-404-error-page-with-standard-http-package
	r := mux.NewRouter()
	r.HandleFunc("/time", timeHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/", r)
	http.ListenAndServe(portString, nil)
}

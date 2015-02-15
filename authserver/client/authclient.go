//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015
//
// Package exposes AuthClient as interface to authserver. Exposes methods
// to construct a new AuthClient as well as Get() and Set() users. Both
// functions able to use request helper function because authserver implements
// endpoints as GET rather than GET and POST.
package client

import (
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	AUTH_SCHEME = "http"
)

// Host and port stored as strings, with
// port expected in form ':8080'.
type AuthClient struct {
	host   string
	port   string
	client *http.Client
}

// Returns new auth client to calling function after initializing
// new http client with timeoutMS as Timeout.
func NewAuthClient(host string, port string, timeoutMS time.Duration) (ac *AuthClient) {
	t := timeoutMS
	c := &http.Client{Timeout: t}
	ac = &AuthClient{host: host, port: port, client: c}
	return
}

// Calls private request method with "get" as parameter
// and map of cookie to uuid. Performs no error checking
// on UUID or name before submission. Returns name if found
// by authserver and empty otherwise. Error associated with
// HTTP request are returned to caller.
func (ac *AuthClient) Get(uuid string) (name string, err error) {
	log.Trace("auth: Get called.")
	params := map[string]string{"cookie": uuid}
	name, err = ac.request("get", params)
	log.Trace("auth: Get complete.")
	return
}

// Calls private request method with "set" as parameter
// and map of cookie to uuid, and name to name. Performs no error
// checking on UUID or name. Error associated with
// HTTP request is returned to caller.
func (ac *AuthClient) Set(uuid string, name string) (err error) {
	log.Trace("auth: Set called.")
	params := map[string]string{"cookie": uuid, "name": name}
	_, err = ac.request("set", params)
	log.Trace("auth: Set complete.")
	return
}

// Takes the request path as an argument along with a map of parameters. Map is encoded
// into URL then submitted via HTTP GET request to authserver. Returns the content of the
// response as a string and error if request failed.
func (ac *AuthClient) request(path string, params map[string]string) (contents string, err error) {
	log.Trace("auth: Request called.")

	var resp *http.Response
	var body []byte

	uri := url.URL{Scheme: AUTH_SCHEME, Host: ac.host + ac.port, Path: path}
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	uri.RawQuery = values.Encode()

	log.Debug("auth: Requesting URI - " + uri.String())
	if resp, err = ac.client.Get(uri.String()); err != nil {
		return
	}

	// Call to close response body will cause
	// panic unless error on call to client.Get
	// is non-nil. Calling here, after error checking
	// ensures response is valid.
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	contents = string(body)
	log.Trace("auth: Request complete.")
	return
}

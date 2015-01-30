//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, January 2015
//
// Package encapsulates cookie functionality needed by personal time server.
// Provides methods for creating a new cookie with relevant fields as well
// as returning the value from the uuid cookie.
package cookie

import (
	log "github.com/cihub/seelog"
	"net/http"
)

const (
	COOKIE_NAME = "uuid"
	COOKIE_PATH = "/"
	MAX_AGE     = 86400
	DELETE_AGE  = -1
)

// Returns address of new cookie with 'uuid' name, value set to value
// path to '/' and age set accordingly. Should utilize MAX_AGE when
// creating, and DELETE_AGE when intending to delete cookie with overwright.
func NewCookie(value string, age int) *http.Cookie {
	c := http.Cookie{Name: COOKIE_NAME, Value: value, Path: COOKIE_PATH, MaxAge: age}
	return &c
}

// Given *http.Request, returns value of uuid cookie if present.
// Returns "" and error if no cookie found named uuid. Performs
// no check on validity of uuid before returning.
func UUIDCookieValue(r *http.Request) (string, error) {
	log.Debug("Attempting to read " + COOKIE_NAME + " cookie from request.")
	cookie, err := r.Cookie(COOKIE_NAME)
	if err == http.ErrNoCookie {
		log.Debug(COOKIE_NAME + " cookie not found in request.")
		return "", http.ErrNoCookie
	}
	log.Debug("Cookie value: " + cookie.Value)
	return cookie.Value, nil
}

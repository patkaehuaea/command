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
	"errors"
	log "github.com/cihub/seelog"
	"github.com/patkaehuaea/command/authserver/people"
	"net/http"
)

const (
	COOKIE_NAME = "uuid"
	COOKIE_PATH = "/"
	MAX_AGE     = 86400
	DELETE_AGE  = -1
	DELETE_VALUE = "deleted"
)

// Returns address of new cookie with 'uuid' name, value set to value
// path to '/' and age set accordingly. Should utilize MAX_AGE when
// creating, and DELETE_AGE when intending to delete cookie with overwright.
func NewCookie(value string, age int) *http.Cookie {
	c := http.Cookie{Name: COOKIE_NAME, Value: value, Path: COOKIE_PATH, MaxAge: age}
	return &c
}

func UUID(r *http.Request) (uuid string, err error) {
	log.Trace("cookie: getting uuid from " + COOKIE_NAME + " cookie.")

	var cookie *http.Cookie
	if cookie, err = r.Cookie(COOKIE_NAME); err != nil {
		return
	}

	if people.IsValidUUID(cookie.Value) {
		uuid = cookie.Value
		return
	}

	err = errors.New("cookie: value not valid uuid")
	return
}

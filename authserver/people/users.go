//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, January 2015
//
// Package contains two types, Person and Users. Intended to be used as
// state tracking mechanism for simple server. Initialization of a
// Users type creates a map of id -> *Person. Access to the map is
// gated via RWMutex. Constructors exist for both Person and
// Users structs to allow for easy initialization.
package people

import (
	log "github.com/cihub/seelog"
	"github.com/patkaehuaea/command/authserver/backup"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// First name, or first and last name in English characters with intervening space.
// Minimum two characters and max length 71 characters including space.
const (
	NAME_REGEX = "^[a-zA-Z]{2,35} {0,1}[a-zA-Z]{0,35}$"
	UUID_REGEX = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
)

type UserStore struct {
	sync.RWMutex
	users map[string]string
}

// Adds a *Person to users map. Acquires RW lock before accessing resource.
func (u *UserStore) Add(id string, name string) {
	u.Lock()
	u.users[id] = name
	u.Unlock()
}

func (u *UserStore) Dump(dumpFile string) (err error) {
	copy := make(map[string]string)
	u.Lock()
	for uuid, name := range u.users {
		copy[uuid] = name
	}
	u.Unlock()

	if err = backup.Write(dumpFile, copy); err != nil {
		log.Error(err)
	}
	return
}

// Deletes *Person from users map whose ID is p.ID. Acquires RW lock before accessing resource.
func (u *UserStore) Delete(id string, name string) {
	u.Lock()
	delete(u.users, id)
	u.Unlock()
}

// Performs read lock on Users. Returns true
// if user with id exists in map. Returns false
// otherise.
func (u *UserStore) Exists(id string) bool {
	u.RLock()
	_, ok := u.users[id]
	u.RUnlock()
	return ok
}

// Uses people.NAME_REGEX to determine if name passed as
// parameter is valid.
func IsValidName(name string) bool {
	match, err := regexp.MatchString(NAME_REGEX, name)
	if err != nil {
		log.Error(err)
	}
	return match
}

func IsValidUUID(value string) bool {
	match, err := regexp.MatchString(UUID_REGEX, value)
	if err != nil {
		log.Error(err)
	}
	return match
}

func (u *UserStore) Load(dumpFile string) (err error) {
	u.Lock()
	err = backup.Read(dumpFile, u.users)
	u.Unlock()
	return
}

// Performs read lock on Users and returns
// name of user with id. If not found, returns
// empty string.
func (u *UserStore) Name(id string) (name string) {
	u.RLock()
	name = u.users[id]
	u.RUnlock()
	return
}

// Returns pointer to object of Users type. Map containing
// state is initialized and ready for use.
func NewUsers() *UserStore {
	return &UserStore{users: make(map[string]string)}
}

func (u *UserStore) Persist(dumpFile string, wait time.Duration) {
	for {
		log.Trace("database: Beginning persist dump.")
		if err := u.Dump(dumpFile); err != nil {
			log.Error(err)
		}
		log.Trace("database: Sleeping for " + wait.String())
		time.Sleep(wait)
	}
}

// For simplicity, was implimented as call to OS executable, but
// should be replaced with uuid package.
func UUID() string {
	out, err := exec.Command("/usr/bin/uuidgen").Output()
	if err != nil {
		log.Error(err)
		return ""
	}
	// Command returns newline at end and must be stripped before use
	// otherwise SetCookie will fail.
	uuid := strings.TrimSuffix(string(out), "\n")
	return uuid
}

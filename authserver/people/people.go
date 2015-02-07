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
	"errors"
	log "github.com/cihub/seelog"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

// First name, or first and last name in English characters with intervening space.
// Minimum two characters and max length 71 characters including space.
const (
	NAME_REGEX = "^[a-zA-Z]{2,35} {0,1}[a-zA-Z]{0,35}$"
    UUID_REGEX = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
)

var (
	ErrInvalidPerson = errors.New("person: uuid or name not valid")
)
type Person struct {
	Name string
	ID   string
}

// Initializes by setting name and calling method
// to create ID. Failure of call to uuid method
// will cause Person ID to be blank.
func NewPerson(uuid string, name string) (p *Person, err error) {
	if IsValidUUID(uuid) && IsValidName(name) {
		p = &Person{Name: name, ID: uuid}
	} else {
		err = ErrInvalidPerson
	}
	return
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

type Users struct {
	sync.RWMutex
	m map[string]*Person
}

// Returns pointer to object of Users type. Map containing
// state is initialized and ready for use.
func NewUsers() *Users {
	return &Users{m: make(map[string]*Person)}
}

// Adds a *Person to users map. Acquires RW lock before accessing resource.
func (u *Users) Add(p *Person) {
	u.Lock()
	u.m[p.ID] = p
	u.Unlock()
}

// Deletes *Person from users map whose ID is p.ID. Acquires RW lock before accessing resource.
func (u *Users) Delete(p *Person) {
	u.Lock()
	delete(u.m, p.ID)
	u.Unlock()
}

// Deletes *Person from users map whose ID is p.ID. Acquires RW lock before accessing resource.
func (u *Users) DumpFile() {
	u.Lock()
	for uuid, person := range u.m {
		log.Info("{ uuid : " + uuid + " , name: " + person.Name + " }")
	}
	u.Unlock()
}

// Performs read lock on Users. Returns true
// if user with id exists in map. Returns false
// otherise.
func (u *Users) Exists(id string) bool {
	u.RLock()
	_, ok := u.m[id]
	u.RUnlock()
	return ok
}

// Performs read lock on Users and returns
// name of user with id. If not found, returns
// empty string.
func (u *Users) Name(id string) (name string) {
	u.RLock()
	defer u.RUnlock()
	p := u.m[id]
	if p != nil {
		name = p.Name
	}
	return
}

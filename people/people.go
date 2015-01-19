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
	"log"
	"os/exec"
	"strings"
	"sync"
)

type Person struct {
	Name string
	ID   string
}

// Initializes by setting name and calling method
// to create ID. Failure of call to uuid method
// will cause Person ID to be blank.
func NewPerson(name string) *Person {
	return &Person{Name: name, ID: uuid()}
}

// For simplicity, was implimented as call to OS executable, but
// should be replaced with uuid package.
func uuid() string {
	out, err := exec.Command("/usr/bin/uuidgen").Output()
	if err != nil {
		// TODO: Handle error case.
		log.Fatal(err)
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

// Locks Users map before adding *Person.
func (u Users) Add(p *Person) {
	u.Lock()
	u.m[p.ID] = p
	u.Unlock()
}

// Locks Users map before deleting *Person.
func (u Users) Delete(p *Person) {
	u.Lock()
	delete(u.m, p.ID)
	u.Unlock()
}

// Performs read lock on Users. Returns true
// if user with id exists in map. Returns false
// otherise.
func (u Users) Exists(id string) bool {
	u.RLock()
	_, ok := u.m[id]
	u.RUnlock()
	return ok
}

// Performs read lock on Users and returns
// name of user with id. If not found, returns
// empty string.
func (u Users) Name(id string) string {
	u.RLock()
	p := u.m[id]
	u.RUnlock()
	if p != nil {
		return p.Name
	} else {
		return ""
	}

}

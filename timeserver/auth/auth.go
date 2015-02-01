package auth

import (
	log "github.com/cihub/seelog"
	"github.com/patkaehuaea/timeserver/cookie"
	"net/http"
)

func Login(uuid string) (name string, err error) {
	log.Debug("Attempting to perform login via remote system.")

	// Remote request returns 200 and name, or 400 if bad.
	// If no name, returns empty string.

	return name, err
}

func Register(name string) (err error) {
	log.Debug("Attempting to register user with remote system.")

	// Generate a new UUID for user.
	uuid := people.UUID()

	// Send Person.ID and Person.Name
}

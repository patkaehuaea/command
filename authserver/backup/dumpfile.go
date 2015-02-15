//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015
//
// Package intended to server as the interface between the in memory user's
// data store and the file system. Implements functions to Read(), and Write()
// a JSON encoded document to the file system along with Exists() and verify()
// helper methods. Common parameters include a filepath/filename and a user
// map[string]string. Read() and Write() methods are guarded by a method
// which checks for presence of the dumpFile before contuing.
package backup

import (
	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"reflect"
)

const (
	BACKUP_FILE_EXTENSION = ".bak"
	DEFAULT_MODE          = 0600
)

// Calls os.Stat() on dumpfile and passes file mode to
// caller if exists. Same expectation as Stat() method
// where err != nil indicates file not present.
func Exists(dumpFile string) (mode os.FileMode, err error) {
	var info os.FileInfo
	if info, err = os.Stat(dumpFile); err != nil {
		log.Trace("backup: File does not exist.")
		return
	}
	mode = info.Mode().Perm()
	return
}

// If dumpFile exists, read the JSON encoded documents into
// users. Undetermined behaviour if map is not string to string.
// Will not unmarshall into users unless file is read successfully.
func Read(dumpFile string, users map[string]string) (err error) {

	var contents []byte

	if _, err = Exists(dumpFile); err != nil {
		log.Trace("backup: Backup does not exist.")
		return
	}

	log.Trace("backup: Reading backup dumpFile.")
	if contents, err = ioutil.ReadFile(dumpFile); err != nil {
		return
	}

	log.Trace("backup: Deserializing into users.")
	err = json.Unmarshal(contents, &users)
	return
}

// Credit for advice on reflect package and DeepEqual: http://goo.gl/VqeDyZ
func verify(dumpFile string, original map[string]string) (err error) {
	compare := make(map[string]string)
	if err = Read(dumpFile, compare); err != nil {
		return
	}
	if equal := reflect.DeepEqual(original, compare); !equal {
		err = errors.New("backup: Backup data not equal to original.")
		return
	}

	return
}

// Expects map passed as parameter to be copy of main data store. Function
// writes JSON encoded document to disk given user parameter. Will rename
// existing dumpFile, but will not delete until new dumpFile can be parsed
// and verified to contain data that is identical to users.
func Write(dumpFile string, users map[string]string) (err error) {

	var mode os.FileMode
	var data []byte
	backup := dumpFile + BACKUP_FILE_EXTENSION

	if mode, err = Exists(dumpFile); err == nil {
		log.Trace("backup: Renaming original dumpFile.")
		if err = os.Rename(dumpFile, backup) ; err != nil {
			return
		}
	} else {
		// ioutil.WriteFile needs default mode or dumpFile will
		// be created with no permission bits set.
		mode = DEFAULT_MODE
	}

	log.Trace("backup: Serializing users map.")
	if data, err = json.Marshal(&users); err != nil {
		return
	}

	log.Trace("backup: Writing dumpFile to disk.")
	if err = ioutil.WriteFile(dumpFile, data, mode); err != nil {
		return
	}

	log.Trace("backup: Verifying dumpFile.")
	if err = verify(dumpFile, users); err != nil {
		return
	}

	// A .bak dumpFile will not exist if the original dumpFile was not present
	// when Write() was called. This block must use a separate error
	// variable or function will incorrectly report an error.
	if _, bakErr := Exists(backup) ; bakErr == nil {
		err = os.Remove(backup)
	}

	return
}

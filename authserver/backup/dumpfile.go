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

// Exists when err == nil, same expected use as simply
// calling os.state, but get additional data.
func Exists(dumpFile string) (mode os.FileMode, err error) {
	var info os.FileInfo
	if info, err = os.Stat(dumpFile); err != nil {
		log.Trace("backup: File does not exist.")
		return
	}
	mode = info.Mode().Perm()
	return
}

func Read(dumpFile string, users map[string]string) (err error) {

	var contents []byte

	if _, err = Exists(dumpFile); err != nil {
		log.Trace("backup: Backup does not exist.")
	}

	log.Trace("backup: Reading backup dumpFile.")
	if contents, err = ioutil.ReadFile(dumpFile); err != nil {
		return
	}

	log.Trace("backup: Deserializing into users.")
	if err = json.Unmarshal(contents, &users); err != nil {
		return
	}

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

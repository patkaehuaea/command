package backup

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    log "github.com/cihub/seelog"
    "os"
    "reflect"
)

const (
    BACKUP_FILE_EXTENSION = ".bak"
    DEFAULT_MODE = 0600
)

func exists(file string)  (mode os.FileMode, err error) {
    var info os.FileInfo
    if info , err = os.Stat(file) ; err != nil {
        return
    }
    mode = info.Mode().Perm()
    return
}

func Read(file string, users map[string]string) (err error) {

    var contents []byte

    if _ , err = exists(file) ; err != nil {
        log.Trace("dumpfile: backup does not exist")
        return
    }

    log.Trace("dumpfile: reading backup file")
    if contents , err = ioutil.ReadFile(file) ; err != nil {
        log.Error(err)
        return
    }

    log.Trace("dumpfile: deserializing into userss")
    if err = json.Unmarshal(contents, &users) ; err != nil {
        log.Error(err)
        return
    }

    return
}

func rename(file string) (err error) {
    if err :=  os.Rename(file, file + BACKUP_FILE_EXTENSION) ; err != nil {
      log.Error(err)
    }
    return
}

// Credit for advice on reflect package and DeepEqual: http://goo.gl/VqeDyZ
func verify(file string, original map[string]string) (err error) {
    compare := make(map[string]string)
    if err =  Read(file, compare) ; err != nil {
        log.Error(err)
        return
    }
    if equal := reflect.DeepEqual(original, compare) ; !equal {
        err = errors.New("dumpfile: backup data not equal to original")
        return
    }

    return
}

func Write(file string, users map[string]string) (err error) {

    var mode os.FileMode
    var data []byte

    if mode, err = exists(file) ; err == nil {
        log.Trace("dumpfile: renaming original file")
        if err = rename(file) ; err != nil {
            log.Error(err)
            return
        }
    } else {
        // ioutil.WriteFile needs default mode or file will
        // be created with no permission bits set.
        mode = DEFAULT_MODE
    }

    log.Trace("dumpfile: serializing users map")
    if data, err = json.Marshal(&users) ; err != nil {
        log.Error(err)
        return
    }

    log.Trace("dumpfile: writing file to disk")
    if err = ioutil.WriteFile(file, data, mode); err != nil {
        log.Error(err)
        return
    }

    log.Trace("dumpfile: verifying file")
    if err = verify(file, users) ; err != nil {
        log.Error(err)
        return
    }

    log.Trace("dumpfile: removing .bak file")
    if err = os.Remove(file + BACKUP_FILE_EXTENSION) ; err != nil {
        log.Error(err)
    }

    return
}
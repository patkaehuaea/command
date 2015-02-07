//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015

package stats

import (
    // log "github.com/cihub/seelog"
    "errors"
    "sync"
)

const (
    START_VALUE = 0
    MIN_VALUE = 0
)

type ConcurrentRequests struct {
    sync.RWMutex
    count int
    max int
}

func NewCR(max int) (cr *ConcurrentRequests) {
    cr = &ConcurrentRequests{count: START_VALUE, max: max}
    return
}

func (cr *ConcurrentRequests) Add() (success bool, err error) {
    cr.Lock()
    if cr.count < cr.max {
        cr.count = cr.count + 1
        success = true
    }
    cr.Unlock()
    if success == false {
        err = errors.New("stats: count >= " + string(cr.max) + " , Add() not allowed.")
    }
    return
}

func (cr *ConcurrentRequests) Subtract() (success bool, err error) {
    cr.Lock()
    if cr.count > MIN_VALUE {
        cr.count = cr.count - 1
        success = true
    }
    cr.Unlock()
    if success == false {
        err = errors.New("stats: count <= " + string(MIN_VALUE) + " , Subtract() not allowed.")
    }
    return
}

func (cr *ConcurrentRequests) Current() (current int) {
    cr.Lock()
    current = cr.count
    cr.Unlock()
    return
}
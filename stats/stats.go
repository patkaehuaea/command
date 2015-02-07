//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015

package stats

import (
    "sync"
)

const (
    START_VALUE = 0
    CONCURRENT_REQUESTS = "concurret requests"
)
type ConcurrentRequests struct {
    sync.RWMutex
    m map[string]int
}

func NewCR() (cr *ConcurrentRequests) {
    cr = &ConcurrentRequests{m: make(map[string]int)}
    cr.m[CONCURRENT_REQUESTS] = START_VALUE
    return
}

func (cr *ConcurrentRequests) Add() {
    cr.Lock()
    cr.m[CONCURRENT_REQUESTS] = cr.m[CONCURRENT_REQUESTS] + 1
    cr.Unlock()
}

func (cr *ConcurrentRequests) Subtract() {
    cr.Lock()
    cr.m[CONCURRENT_REQUESTS] = cr.m[CONCURRENT_REQUESTS] - 1
    cr.Unlock()
}

func (cr *ConcurrentRequests) Current() (current int){
    cr.Lock()
    current = cr.m[CONCURRENT_REQUESTS]
    cr.Unlock()
    return
}
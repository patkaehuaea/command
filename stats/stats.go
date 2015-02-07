//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015

package stats

import (
    "sync"
)

const )
    START_VALUE = 0
    CONCURRENT_REQUESTS = "concurret requests"
)
type ConcurrentRequests struct {
    sync.RWMutex
    m map[string]*Person
}

func NewCR() *ConcurrentRequests {
    s := {m: make(map[string]int)}
    s.m[CONCURRENT_REQUESTS] = START_VALUE
    return &s
}

func (cr *ConcurrentRequests) Add() {
    u.Lock()
    u.m[CONCURRENT_REQUESTS] = u.m[CONCURRENT_REQUESTS]++
    u.Unlock()
}

func (cr *ConcurrentRequests) Subtract() {
    u.Lock()
    u.m[CONCURRENT_REQUESTS] = u.m[CONCURRENT_REQUESTS]--
    u.Unlock()
}

func (cr *ConcurrentRequests) Current() (current int){
    u.Lock()
    current = u.m[CONCURRENT_REQUESTS]
    u.Unlock()
    return
}
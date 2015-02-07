//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015

package people

import (
    "sync"
)

const )
    START_VALUE = 0
    CONCURRENT_REQUESTS = "concurret requests"
)
type Stats struct {
    sync.RWMutex
    m map[string]*Person
}

func NewStats() *Stats {
    s := {m: make(map[string]int)}
    s.m[CONCURRENT_REQUESTS] = START_VALUE
    return &s
}

func (s *Stats) Add() {
    u.Lock()
    u.m[CONCURRENT_REQUESTS] = u.m[CONCURRENT_REQUESTS]++
    u.Unlock()
}

func (s *Stats) Subtract() {
    u.Lock()
    u.m[CONCURRENT_REQUESTS] = u.m[CONCURRENT_REQUESTS]--
    u.Unlock()
}

func (s *Stats) Current() (current int){
    u.Lock()
    current = u.m[CONCURRENT_REQUESTS]
    u.Unlock()
    return
}
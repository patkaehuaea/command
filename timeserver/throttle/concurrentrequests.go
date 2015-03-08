//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015

package throttle

import (
	"errors"
	"fmt"
	"sync"
)

const (
	START_VALUE = 0
	MIN_VALUE   = 0
)

type ConcurrentRequests struct {
	sync.RWMutex
	count int
	max   int
}

func New(max int) (cr *ConcurrentRequests) {
	cr = &ConcurrentRequests{count: START_VALUE, max: max}
	return
}

func (cr *ConcurrentRequests) Add() (err error) {
	cr.Lock()
	if cr.count < cr.max {
		cr.count = cr.count + 1
	} else {
		err = errors.New(fmt.Sprintf("%s %d %s", "stats: Exceded threashold of ", cr.count, " concurrent requests"))
	}
	cr.Unlock()
	return
}

func (cr *ConcurrentRequests) Current() (current int) {
	cr.Lock()
	current = cr.count
	cr.Unlock()
	return
}

func (cr *ConcurrentRequests) Subtract() (err error) {
	cr.Lock()
	if cr.count > MIN_VALUE {
		cr.count = cr.count - 1
	} else {
		err = errors.New(fmt.Sprintf("%s %d", "stats: Count already at MIN_VALUE ", cr.count))
	}
	cr.Unlock()
	return
}

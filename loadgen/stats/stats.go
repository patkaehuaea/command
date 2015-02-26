//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015

package stats

import (
	"strconv"
	"sync"
)

const (
	TOTAL_KEY   = "Total"
	KEY_100     = "100s"
	KEY_200     = "200s"
	KEY_300     = "300s"
	KEY_400     = "400s"
	KEY_500     = "500s"
	ERROR_KEY   = "Errors"
	START_VALUE = 0
)

type Statistics struct {
	sync.RWMutex
	counters map[string]int
}

func NewCounters() (s *Statistics) {
	s = &Statistics{counters: make(map[string]int)}
	s.counters[TOTAL_KEY] = START_VALUE
	s.counters[KEY_100] = START_VALUE
	s.counters[KEY_200] = START_VALUE
	s.counters[KEY_300] = START_VALUE
	s.counters[KEY_400] = START_VALUE
	s.counters[KEY_500] = START_VALUE
	s.counters[ERROR_KEY] = START_VALUE
	return
}

func (s *Statistics) Copy() (copy map[string]int) {
	copy = make(map[string]int)
	s.RLock()
	for statusCode, count := range s.counters {
		copy[statusCode] = count
	}
	s.RUnlock()
	return
}

// Expect one.
func (s *Statistics) Increment(statistic string, delta int) (err error) {
	s.Lock()
	s.counters[TOTAL_KEY] = s.counters[TOTAL_KEY] + delta
	s.counters[statistic] = s.counters[statistic] + delta
	s.Unlock()
	return
}

func (s *Statistics) Reset(statistic string) {
	s.Lock()
	s.counters[statistic] = START_VALUE
	s.Unlock()

}

func (s *Statistics) String() (output string) {
	copy := s.Copy()
	output = TOTAL_KEY + ":\t" + strconv.Itoa(copy[TOTAL_KEY]) + "\n"
	output = output + KEY_100 + ":\t" + strconv.Itoa(copy[KEY_100]) + "\n"
	output = output + KEY_200 + ":\t" + strconv.Itoa(copy[KEY_200]) + "\n"
	output = output + KEY_300 + ":\t" + strconv.Itoa(copy[KEY_300]) + "\n"
	output = output + KEY_400 + ":\t" + strconv.Itoa(copy[KEY_400]) + "\n"
	output = output + KEY_500 + ":\t" + strconv.Itoa(copy[KEY_500]) + "\n"
	output = output + ERROR_KEY + ":\t" + strconv.Itoa(copy[ERROR_KEY]) + "\n"
	return
}

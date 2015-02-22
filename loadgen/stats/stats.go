//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015

package stats

import (
	"errors"
	"regexp"
	"strconv"
	"sync"
)

const (
	TOTAL_KEY   = "Total"
	KEY_100     = "100"
	KEY_200     = "200"
	KEY_300     = "300"
	KEY_400     = "400"
	KEY_500     = "500"
	ERROR_KEY   = "Errors"
	START_VALUE = 0
	// Matches 100 through 599 so not exact, but
	// better than nothing.
	HTTP_STATUS_CODE_REGEX = "^[1-5][0-9][0-9]$"
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

// Converts string adhering to HTTP_STATUS_CODE_REGEX to
// one of the allowed key constants defined in the package.
// Behavior undefined if receives non-compliant string.
func convertToKey(httpStatusCode string) (code string) {
	code = string(httpStatusCode[0])
	code = code + "00"
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

func (s *Statistics) Error() {
	s.Lock()
	s.counters[TOTAL_KEY] = s.counters[TOTAL_KEY] + 1
	s.counters[ERROR_KEY] = s.counters[ERROR_KEY] + 1
	s.Unlock()
}

func (s *Statistics) Increment(httpStatusCode int) (err error) {
	code := strconv.Itoa(httpStatusCode)

	if match := isValidStatus(code); match {

		if code = convertToKey(code); err != nil {
			return
		}

		s.Lock()
		s.counters[TOTAL_KEY] = s.counters[TOTAL_KEY] + 1
		s.counters[code] = s.counters[code] + 1
		s.Unlock()
	}
	return
}

// Uses stats.HTTP_STATUS_CODE_REGEX to determine if int passed as
// parameter is valid.
func isValidStatus(httpStatusCode string) (match bool) {
	match, _ = regexp.MatchString(HTTP_STATUS_CODE_REGEX, httpStatusCode)
	return
}

func (s *Statistics) Reset(httpStatusCode int) (err error) {
	code := strconv.Itoa(httpStatusCode)
	if match := isValidStatus(code); match {
		s.Lock()
		s.counters[code] = START_VALUE
		s.Unlock()
	}

	err = errors.New("stats: Unable to reset " + code + " to 0.")
	return
}

func (s *Statistics) String() (output string) {
	copy := s.Copy()
	output = "Total:  " + strconv.Itoa(s.counters[TOTAL_KEY]) + "\n"
	output = output + "100s:   " + strconv.Itoa(copy[KEY_100]) + "\n"
	output = output + "200s:   " + strconv.Itoa(copy[KEY_200]) + "\n"
	output = output + "300s:   " + strconv.Itoa(copy[KEY_300]) + "\n"
	output = output + "400s:   " + strconv.Itoa(copy[KEY_400]) + "\n"
	output = output + "500s:   " + strconv.Itoa(copy[KEY_500]) + "\n"
	output = output + "Errors: " + strconv.Itoa(copy[ERROR_KEY]) + "\n"
	return
}

//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015
//
// Provides counters via map of string to int. Keys are defined
// and initialized to start value. All defined methods are protected
// via mutex and are thread safe. There is no error checking performed
// in methods that take statistic parameter as string. It is expected
// that callers will use existing counters. It is highly recommended
// to pass a status code to the Key() method when calling methods with
// statistic parameters.

package stats

import (
	"fmt"
	"sync"
)

const (
	TOTAL_KEY          = "Total"
	KEY_100            = "100s"
	KEY_200            = "200s"
	KEY_300            = "300s"
	KEY_400            = "400s"
	KEY_500            = "500s"
	ERROR_KEY          = "Errors"
	START_VALUE        = 0
	KEY_LOOKUP_DIVISOR = 100
)

var convert = map[int]string{
	1: "100s",
	2: "200s",
	3: "300s",
	4: "400s",
	5: "500s",
}

type Counter struct {
	sync.RWMutex
	counters map[string]int
}

// Converts an http status code to a string
// for accessing the right counter in the
// counters map. Returns stats.ERROR_KEY if
// status code invalid or key not found in map.
func Key(httpStatusCode int) string {
	key, ok := convert[httpStatusCode/KEY_LOOKUP_DIVISOR]
	if !ok {
		key = ERROR_KEY
	}
	return key
}

// Return pointer to Counter struct whose map is
// initialized with all keys defined as consts
// in stats package. Value of each key set to
// START_VALUE.
func New() (c *Counter) {
	c = &Counter{counters: make(map[string]int)}
	c.counters[TOTAL_KEY] = START_VALUE
	c.counters[KEY_100] = START_VALUE
	c.counters[KEY_200] = START_VALUE
	c.counters[KEY_300] = START_VALUE
	c.counters[KEY_400] = START_VALUE
	c.counters[KEY_500] = START_VALUE
	c.counters[ERROR_KEY] = START_VALUE
	return
}

// Returns copy of counters map with values identical
// to original when method was called.
func (c *Counter) Copy() (copy map[string]int) {
	copy = make(map[string]int)
	c.RLock()
	for statusCode, count := range c.counters {
		copy[statusCode] = count
	}
	c.RUnlock()
	return
}

// Gets value of statistic. Fetching statistic not in map
// will yield result of zero. Callers should use
// stats.Key function if unsure if which statistics are present.
func (c *Counter) Get(statistic string) (count int) {
	c.RLock()
	count = c.counters[statistic]
	c.RUnlock()
	return
}

// Increments statistic by delta, and stats.TOTAL_KEY by same amount.
// Callers should use stats.Key function if unsure of which statistics
// are present.
func (c *Counter) Increment(statistic string, delta int) {
	c.Lock()
	c.counters[TOTAL_KEY] = c.counters[TOTAL_KEY] + delta
	c.counters[statistic] = c.counters[statistic] + delta
	c.Unlock()
}

// Resets all key-value pairs in c.Counter to stats.START_VALUE.
func (c *Counter) Reset() {
	c.Lock()
	for k, _ := range c.counters {
		c.counters[k] = START_VALUE
	}
	c.Unlock()
}

// Prints counters to screen.
func (c *Counter) Print() (output string) {
	copy := c.Copy()
	fmt.Printf("%s:\t %d\n", TOTAL_KEY, copy[TOTAL_KEY])
	fmt.Printf("%s:\t %d\n", KEY_100, copy[KEY_100])
	fmt.Printf("%s:\t %d\n", KEY_200, copy[KEY_200])
	fmt.Printf("%s:\t %d\n", KEY_300, copy[KEY_300])
	fmt.Printf("%s:\t %d\n", KEY_400, copy[KEY_400])
	fmt.Printf("%s:\t %d\n", KEY_500, copy[KEY_500])
	fmt.Printf("%s:\t %d\n", ERROR_KEY, copy[ERROR_KEY])
	return
}

//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015

package counters

import "sync"

const START_VALUE = 0

type Counter struct {
	sync.RWMutex
	counters map[string]int
}

func New(keys []string) (c *Counter) {
	c = &Counter{counters: make(map[string]int)}
	for _, element := range keys {
		c.counters[element] = START_VALUE
	}
	return
}

// Returns copy of counters map with values identical
// to original when method was called.
func (c *Counter) Copy() (copy map[string]int) {
	copy = make(map[string]int)
	c.RLock()
	for item, value := range c.counters {
		copy[item] = value
	}
	c.RUnlock()
	return
}

// Gets value of item. Fetching item not in map
// will yield result of zero. Callers should use
// stats.Key function if unsure if which items are present.
func (c *Counter) Get(item string) (count int) {
	c.RLock()
	count = c.counters[item]
	c.RUnlock()
	return
}

// Increments item by delta, and stats.TOTAL_KEY by same amount.
// Callers should use stats.Key function if unsure of which items
// are present.
func (c *Counter) Increment(item string, delta int) {
	c.Lock()
	c.counters[item] = c.counters[item] + delta
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

//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015

package stats

import (
	"fmt"
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

type Counter struct {
	sync.RWMutex
	counters map[string]int
}

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

func (c *Counter) Copy() (copy map[string]int) {
	copy = make(map[string]int)
	c.RLock()
	for statusCode, count := range c.counters {
		copy[statusCode] = count
	}
	c.RUnlock()
	return
}

func (c *Counter) Get(statistic string) (count int) {
	c.RLock()
	count = c.counters[statistic]
	c.RUnlock()
	return
}

func (c *Counter) Increment(statistic string, delta int) {
	c.Lock()
	c.counters[TOTAL_KEY] = c.counters[TOTAL_KEY] + delta
	c.counters[statistic] = c.counters[statistic] + delta
	c.Unlock()
}

func (c *Counter) Reset(statistic string) {
	c.Lock()
	c.counters[statistic] = START_VALUE
	c.Unlock()
}

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

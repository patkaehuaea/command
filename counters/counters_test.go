//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015
//
// Performs  unit testing of stats package. Tests are performed within
// the context of the load generator application for clarity and ease of use.
// Please note that this package is separate from counters becase of an import
// circle error I was unable to resolve. Shouldn't pose a problem as unit tests
// focus on public interface of counters package.

package counters_test

import (
	"github.com/patkaehuaea/command/counters"
	"sync"
	"testing"
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

var (
	dataSetOne = map[string]int{
		"100s":   1,
		"200s":   22,
		"300s":   3,
		"400s":   15,
		"500s":   5,
		"Errors": 100,
	}

	dataSetTwo = map[string]int{
		"100s":   6,
		"200s":   -1,
		"300s":   10,
		"400s":   0,
		"500s":   5,
		"Errors": -99,
	}

	dataSetThree = map[string]int{
		"100s":   8,
		"200s":   2,
		"300s":   0,
		"400s":   99,
		"500s":   5,
		"Errors": 1,
	}

	dataSetFour = map[string]int{
		"100s":   0,
		"200s":   0,
		"300s":   0,
		"400s":   0,
		"500s":   0,
		"Errors": 0,
	}

	expectedCopy = map[string]int{
		"Total":  146,
		"100s":   1,
		"200s":   22,
		"300s":   3,
		"400s":   15,
		"500s":   5,
		"Errors": 100,
	}

	expectedIncrement = map[string]int{
		"Total":  182,
		"100s":   15,
		"200s":   23,
		"300s":   13,
		"400s":   114,
		"500s":   15,
		"Errors": 2,
	}
)

func increment(c *counters.Counter, data map[string]int) {
	for statistic, delta := range data {
		c.Increment(TOTAL_KEY, delta)
		c.Increment(statistic, delta)
	}
}

func setup() *counters.Counter {
	keys := []string{
		TOTAL_KEY,
		KEY_200,
		KEY_300,
		KEY_400,
		KEY_500,
		ERROR_KEY,
	}
	return counters.New(keys)
}

func TestCopy(t *testing.T) {
	counter := setup()
	increment(counter, dataSetOne)
	copy := counter.Copy()
	for k, v := range expectedCopy {
		if v != copy[k] {
			t.Errorf("copy %s: expected %d, got %d", k, v, copy[k])
		}
	}
}

func TestGet(t *testing.T) {
	counter := setup()
	increment(counter, dataSetOne)
	total := counter.Get(TOTAL_KEY)
	expected := 146
	if total != expected {
		t.Errorf("%s: expected %d, got %d", TOTAL_KEY, expected, total)
	}
}

func TestIncrement(t *testing.T) {
	counter := setup()
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		increment(counter, dataSetOne)
	}()
	go func() {
		defer wg.Done()
		increment(counter, dataSetTwo)
	}()
	go func() {
		defer wg.Done()
		increment(counter, dataSetThree)
	}()
	go func() {
		defer wg.Done()
		increment(counter, dataSetFour)
	}()
	wg.Wait()

	actual := counter.Copy()

	for k, v := range expectedIncrement {
		if v != actual[k] {
			t.Errorf("counter %s: expected %d, got %d", k, v, actual[k])
		}
	}
}

func TestReset(t *testing.T) {
	counter := setup()
	increment(counter, dataSetOne)
	counter.Reset()
	actual := counter.Copy()
	for k, v := range actual {
		if v != START_VALUE {
			t.Errorf("counter %s: expected %d, got %d", k, v, actual[k])
		}
	}
}

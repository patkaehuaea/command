package stats_test

import (
	"github.com/patkaehuaea/command/loadgen/stats"
	"sync"
	"testing"
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

    expectedReset = map[string]int{
        "Total":  67,
        "100s":   6,
        "200s":   21,
        "300s":   13,
        "400s":   15,
        "500s":   10,
        "Errors": 1,
    }
)

func increment(c *stats.Counter, data map[string]int) {
    for statistic, delta := range data {
        c.Increment(statistic, delta)
    }
}

func TestCopy(t *testing.T) {
    counter := stats.New()
    var copy map[string]int
    var wg sync.WaitGroup
    wg.Add(3)
    go func() {
        defer wg.Done()
        increment(counter, dataSetOne)
    }()
    go func() {
        defer wg.Done()
        copy = counter.Copy()
    }()
    go func() {
        defer wg.Done()
        increment(counter, dataSetTwo)
    }()
    wg.Wait()

    for k, v := range expectedCopy {
        if v != copy[k] {
            t.Errorf("copy %s: expected %d, got %d", k, v, copy[k])
        }
    }
}

func TestGet(t *testing.T) {
    counter := stats.New()
    var total int
    var wg sync.WaitGroup
    wg.Add(3)
    go func() {
        defer wg.Done()
        increment(counter, dataSetOne)
    }()
    go func() {
        defer wg.Done()
        total = counter.Get(stats.TOTAL_KEY)
    }()
    go func() {
        defer wg.Done()
        increment(counter, dataSetTwo)
    }()
    wg.Wait()

    expected := 146

    if total != expected {
        t.Errorf("%s: expected %d, got %d", stats.TOTAL_KEY, expected, total)
    }
}

func TestIncrement(t *testing.T) {
	counter := stats.New()
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
    counter := stats.New()
    var wg sync.WaitGroup
    wg.Add(3)
    go func() {
        defer wg.Done()
        increment(counter, dataSetOne)
    }()
    go func() {
        defer wg.Done()
        counter.Reset(stats.KEY_100)
    }()
    go func() {
        defer wg.Done()
        increment(counter, dataSetTwo)
    }()
    wg.Wait()

    actual := counter.Copy()

    for k, v := range expectedReset {
        if v != actual[k] {
            t.Errorf("counter %s: expected %d, got %d", k, v, actual[k])
        }
    }
}

//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015
//
// Implements load generator required in assignment-05. The loadgen
// application uses the adjacent stats pacakge for storage
// of counter data and the config package for parsing of
// command line options. See config pacakge for defaults.

package main

import (
	"fmt"
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/counters"
	"net/http"
	"os"
	"time"
)

const (
	DEFAULT_DELTA          = 1
	UNIT_CONVERSION_FACTOR = 1000000
	TOTAL_KEY              = "Total"
	KEY_100                = "100s"
	KEY_200                = "200s"
	KEY_300                = "300s"
	KEY_400                = "400s"
	KEY_500                = "500s"
	ERROR_KEY              = "Errors"
	START_VALUE            = 0
	KEY_LOOKUP_DIVISOR     = 100
)

var (
	counter *counters.Counter
	period  <-chan time.Time
	stop    <-chan time.Time
	client  *http.Client
	convert = map[int]string{
		1: "100s",
		2: "200s",
		3: "300s",
		4: "400s",
		5: "500s",
	}
)

// Converts an http status code to a string
// for accessing the right counter in the
// counters map. Returns ERROR_KEY if
// status code invalid or key not found in map.
func key(httpStatusCode int) string {
	key, ok := convert[httpStatusCode/KEY_LOOKUP_DIVISOR]
	if !ok {
		key = ERROR_KEY
	}
	return key
}

func load(url string, burst int) {
	for {
		for i := 0; i < burst; i++ {
			go request(url)
		}
		<-period

		select {
		case <-stop:
			return
		default:
		}
	}
}

// Prints counters to screen.
func print() {
	copy := counter.Copy()
	fmt.Printf("%s:\t %d\n", TOTAL_KEY, copy[TOTAL_KEY])
	fmt.Printf("%s:\t %d\n", KEY_100, copy[KEY_100])
	fmt.Printf("%s:\t %d\n", KEY_200, copy[KEY_200])
	fmt.Printf("%s:\t %d\n", KEY_300, copy[KEY_300])
	fmt.Printf("%s:\t %d\n", KEY_400, copy[KEY_400])
	fmt.Printf("%s:\t %d\n", KEY_500, copy[KEY_500])
	fmt.Printf("%s:\t %d\n", ERROR_KEY, copy[ERROR_KEY])
	return
}

func request(url string) {
	counter.Increment(TOTAL_KEY, DEFAULT_DELTA)
	if response, err := client.Get(url); err != nil {
		counter.Increment(ERROR_KEY, DEFAULT_DELTA)
		return
	} else {
		defer response.Body.Close()
		counter.Increment(key(response.StatusCode), DEFAULT_DELTA)
	}
}

func init() {
	keys := []string{TOTAL_KEY, KEY_200, KEY_300, KEY_400, KEY_500, ERROR_KEY}
	counter = counters.New(keys)
	period = time.Tick(time.Duration((*config.Burst*UNIT_CONVERSION_FACTOR) / *config.Rate) * time.Microsecond)
	stop = time.Tick(*config.Runtime + *config.LoadTimeoutMS)
	client = &http.Client{Timeout: *config.LoadTimeoutMS}
}

func main() {

	/*
			Paramters surfaced via config pacakge used in this program:

			*config.Rate: average rate of requests (per second)
		  	*config.Burst: number of concurrent requests to issue
		   	*config.LoadTimeout-ms: max time to wait for response
		   	*config.Runtime: number of seconds to process
		   	*config.url: URL to sample
	*/

	load(*config.URL, *config.Burst)
	print()
	os.Exit(0)

}

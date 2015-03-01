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
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/loadgen/stats"
	"net/http"
	"os"
	"time"
)

const (
	DEFAULT_DELTA          = 1
	UNIT_CONVERSION_FACTOR = 1000000
)

var (
	counter *stats.Counter
	period  <-chan time.Time
	stop    <-chan time.Time
	client  *http.Client
)

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

func request(url string) {
	if response, err := client.Get(url); err != nil {
		counter.Increment(stats.ERROR_KEY, DEFAULT_DELTA)
		return
	} else {
		defer response.Body.Close()
		counter.Increment(stats.Key(response.StatusCode), DEFAULT_DELTA)
	}
}

func init() {
	counter = stats.New()
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
	counter.Print()
	os.Exit(0)

}

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
	"encoding/json"
	"fmt"
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/monitorserver/stats"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	period <-chan time.Time
	stop   <-chan time.Time
	client *http.Client
	data   *stats.MonitorData
)

func monitor(urls []string) {
	for {
		for _, url := range urls {
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

    now := time.Now()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	dict := make(map[string]int)
	if err := json.Unmarshal(body, &dict); err != nil {
		return
	}

	for counter, value := range dict {
		data.Add(url, counter, stats.Sample{Time: now, Value: value})
	}
}

func init() {
	period = time.Tick(*config.MonIntSec)
	stop = time.Tick(*config.MonRunSec)
}

func main() {

	/*
	   Paramters surfaced via config pacakge used in this program:
	   *config.MonIntSec
	   *config.MonRunSec
	   *config.MonTargets
	*/

	targets := strings.Split(*config.MonTargets, ",")
	data = stats.New(targets)
	monitor(targets)
	data.Print()
	os.Exit(0)

}

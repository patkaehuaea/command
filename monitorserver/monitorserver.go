//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, February 2015
//
// Implements monitorserver required in assignment-06. The
// application uses the adjacent metrics pacakge for storage
// of collected monitoring data and the config package for parsing
// of command line options. See config pacakge for defaults.

package main

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/monitorserver/metrics"
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
	data   *metrics.Data
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
		log.Error("monitorserver: " + err.Error())
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("monitorserver: " + err.Error())
		return
	}

	dict := make(map[string]int)
	if err := json.Unmarshal(body, &dict); err != nil {
		log.Error("monitorserver: " + err.Error())
		return
	}

	for counter, value := range dict {
		data.Add(url, counter, metrics.Sample{Time: now, Value: value})
	}
}

func init() {
	log.ReplaceLogger(config.Logger)
	if *config.MonTargets == config.MONITOR_TARGETS {
		log.Critical("monitorserver: No targets specified.")
		os.Exit(1)
	}

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
	data = metrics.New(targets)
	monitor(targets)
	data.Print()
	os.Exit(0)

}

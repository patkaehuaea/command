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
    "github.com/patkaehuaea/command/monitorserver/sequence"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "time"
)

var (
    period  <-chan time.Time
    stop    <-chan time.Time
    client  *http.Client
    data    map[string]*sequence.Sequence
)

type Sample struct {

}

func monitor(urls []string) {
    for {
        for _ , url := range urls {
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
    // copy := counter.Copy()
    // fmt.Printf("%s:\t %d\n", TOTAL_KEY, copy[TOTAL_KEY])
    // fmt.Printf("%s:\t %d\n", KEY_100, copy[KEY_100])
    // fmt.Printf("%s:\t %d\n", KEY_200, copy[KEY_200])
    // fmt.Printf("%s:\t %d\n", KEY_300, copy[KEY_300])
    // fmt.Printf("%s:\t %d\n", KEY_400, copy[KEY_400])
    // fmt.Printf("%s:\t %d\n", KEY_500, copy[KEY_500])
    // fmt.Printf("%s:\t %d\n", ERROR_KEY, copy[ERROR_KEY])
    return
}

func request(url string) {

    response, err := http.Get(url);
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return
    }

    stats := make(map[string]int)
    if err := json.Unmarshal(body, &stats) ; err != nil {
        return
    }

    for counter, value := range stats {
        data[url].Add(counter, sequence.Sample{Time: time.Now(), Value: value})
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

    data = make(map[string]*sequence.Sequence)

    for _ , target := range targets {
        sequence := sequence.New()
        data[target] = sequence
    }

    monitor(targets)

    for target, sequence := range data {
        copy := sequence.Copy()
        data, _ := json.MarshalIndent(copy, "", "   ")
        fmt.Println(target)
        fmt.Println(string(data))
    }
    //fmt.Printf("%s", output)
    //print()
    os.Exit(0)

}

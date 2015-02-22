package main

import (
    "fmt"
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/loadgen/stats"
	"net/http"
    "os"
    "time"
)

var (
    statistics *stats.Statistics
    channel <-chan time.Time
    client   *http.Client
)

func load(url string, burst int) {
    for {
        <-channel
        for i := 0 ; i < burst ; i++ {
            go request(url)
        }
    }
}

func period(burst int, rate int) time.Duration {
    return time.Duration(burst * 1000000 / rate) * time.Microsecond
}

func request(url string) (err error) {
    var response *http.Response
    if response, err = client.Get(url); err != nil {
        statistics.Error()
        return
    }

    defer response.Body.Close()
    statistics.Increment(response.StatusCode)
    return
}

func init() {
    statistics = stats.NewCounters()
    channel = time.Tick(period(*config.Burst, *config.Rate))
    client = &http.Client{Timeout: *config.LoadTimeoutMS}
}

func main() {

    /*
       --rate: average rate of requests (per second)
       --burst: number of concurrent requests to issue
       --timeout-ms: max time to wait for response
       --runtime: number of seconds to process
       --url: URL to sample
    */

    go load(*config.URL, *config.Burst)
    time.Sleep(*config.Runtime)
    fmt.Println(statistics.String())
    os.Exit(0)

}

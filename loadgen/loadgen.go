package main

import (
	"fmt"
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/loadgen/stats"
	"net/http"
	"os"
	"time"
)

const KEY_LOOKUP_DIVISOR = 100

var (
	statistics *stats.Statistics
	period     <-chan time.Time
	stop       <-chan time.Time
	client     *http.Client
	convert    = map[int]string{
		1: "100s",
		2: "200s",
		3: "300s",
		4: "400s",
		5: "500s",
	}
)

func key(httpStatusCode int) string {
	key, ok := convert[httpStatusCode/KEY_LOOKUP_DIVISOR]
	if !ok {
		key = stats.ERROR_KEY
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

func request(url string) (err error) {
	var response *http.Response
	if response, err = client.Get(url); err != nil {
		statistics.Increment(stats.ERROR_KEY, 1)
		return
	}

	defer response.Body.Close()
	statistics.Increment(key(response.StatusCode), 1)
	return
}

func init() {
	statistics = stats.New()
	period = time.Tick(time.Duration((*config.Burst*1000000) / *config.Rate) * time.Microsecond)
	stop = time.Tick(*config.Runtime + *config.LoadTimeoutMS)
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

	load(*config.URL, *config.Burst)
	fmt.Println(statistics.String())
	os.Exit(0)

}

package main

import (
	"github.com/patkaehuaea/command/config"
	"github.com/patkaehuaea/command/loadgen/stats"
	"net/http"
	"os"
	"time"
)

const KEY_LOOKUP_DIVISOR = 100

var (
	counter *stats.Counter
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

func request(url string) {
	if response, err := client.Get(url); err != nil {
		counter.Increment(stats.ERROR_KEY, 1)
		return
	} else {
		defer response.Body.Close()
		counter.Increment(key(response.StatusCode), 1)
	}
}

func init() {
	counter = stats.New()
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
	counter.Print()
	os.Exit(0)

}

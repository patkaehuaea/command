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
	interval   <-chan time.Time
	client     *http.Client
	convert    = map[int]string{
		1: "100s",
		2: "200s",
		3: "300s",
		4: "400s",
		5: "500s",
	}
)

func key(httpStatusCode int) (key string) {
	return convert[httpStatusCode/100]
}

func load(url string, burst int) {
	//timeout := time.Tick(time.Duration(2) * *config.Runtime)
	for {
		for i := 0; i < burst; i++ {
			go request(url)
		}
		<-interval

		// Poll for timeout
		// select {
		// 	case <- interval:
		// 	case <- timeout:
		// 		return
		// 	default:
		// }
	}
}

func period(burst int, rate int) time.Duration {
	return time.Duration(burst*1000000/rate) * time.Microsecond
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

	/*
		c.Incr("total")

		if err != nil {
			c.Incr("errors", 1)
			return
		}

		key, ok := convert[response.StatusCode / 100]
		if !ok {
			key = "errors"
		}

		c.Incr(key, 1)
	*/

}

func init() {
	statistics = stats.NewCounters()
	interval = time.Tick(period(*config.Burst, *config.Rate))
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

	// If waiting -> load is no longer a go routine.
	go load(*config.URL, *config.Burst)
	// Adjust sleep so you wait long enough based on
	// requests in flight.
	time.Sleep(*config.Runtime)
	fmt.Println(statistics.String())
	os.Exit(0)

}

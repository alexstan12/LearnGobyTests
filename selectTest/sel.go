package selectTest

import (
	"net/http"
	"time"
)

func Racer(a,b string) (winner string) {
	durationA := measureResponseTime(a)
	durationB := measureResponseTime(b)

	if durationA < durationB {
		return a
	}
	return b
}

func measureResponseTime(url string) time.Duration {
	startTime := time.Now()
	http.Get(url)
	duration := time.Since(startTime)

	return duration
}
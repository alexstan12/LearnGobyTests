package selectTest

import (
	"net/http"
	"time"
)

func Racer(a,b string) (winner string) {
	startTimeA := time.Now()
	http.Get(a)
	durationA := time.Since(startTimeA)

	startTimeB := time.Now()
	http.Get(b)
	durationB := time.Since(startTimeB)

	if durationA < durationB {
		return a
	}
	return b
}
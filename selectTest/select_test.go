package selectTest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T){
	t.Run("return an error if a server doesn't respond within the specified time", func(t *testing.T) {
		server := makeDelayedServer(25 * time.Millisecond)

		defer server.Close()
		_,err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("Expected an error but didn't get one")
		}

	})

	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T){
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		want :=fastServer.URL
		got,err := Racer(slowServer.URL, fastServer.URL)

		if err!= nil {
			t.Error("Didn't expect an error, but got one")
		}
		if got!=want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func makeDelayedServer(delay time.Duration)  *httptest.Server{
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
}

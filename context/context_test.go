package context

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubStore struct {
	response string
}

func (stubS *StubStore) Fetch() string {
	return stubS.response
}

func TestServer(t *testing.T) {
	data := "hello, world"
	svr := Server(&StubStore{data})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	svr.ServeHTTP(response, request)

	if response.Body.String() != data{
		t.Errorf("got %s, want %s", response.Body.String(), data)
	}
}
package context

import (
	"fmt"
	"net/http"
)

func Server(store Store) http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, store.Fetch())
	}
}

type Store interface {
	Fetch() string
	Cancel()
}

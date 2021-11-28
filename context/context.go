package context

import (
	"fmt"
	"net/http"
)

func Server(store Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		data := make(chan string, 1)

		go func(){
			data <- store.Fetch()
		}()

		select{
			case d:= <-data :
				fmt.Fprint(writer, d)
			case <-ctx.Done():
				store.Cancel()
		}
	}
}

type Store interface {
	Fetch() string
	Cancel()
}

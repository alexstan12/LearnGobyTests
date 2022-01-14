package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
}

type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	player := strings.TrimPrefix(request.URL.String(), "/players/")

	fmt.Fprint(writer, p.store.GetPlayerScore(player))
}

//func GetPlayerScore(player string) string {
//	if player == "Pepper" {
//		return "20"
//	} else if player == "Floyd" {
//		return "10"
//	}
//	return ""
//}

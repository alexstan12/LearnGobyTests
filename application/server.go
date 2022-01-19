package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer{
	p := new(PlayerServer)
	p.store = store
	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	p.Handler.ServeHTTP(writer, request)
}

func (p *PlayerServer) playersHandler(writer http.ResponseWriter, request *http.Request){
		player := strings.TrimPrefix(request.URL.Path, "/players/")

		switch request.Method {
		case http.MethodPost:
			p.processWin(writer, player)
		case http.MethodGet:
			p.showScore(writer, player)
		}
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

}

func (p *PlayerServer) showScore(writer http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)
	if score != 0 {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(writer, score)
}

func (p *PlayerServer) processWin(writer http.ResponseWriter, player string) {
	writer.WriteHeader(http.StatusAccepted)
	p.store.RecordWin(player)
	return
}

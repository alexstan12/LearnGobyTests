package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

const jsonContentType = "application/json"

func NewPlayerServer(store PlayerStore) *PlayerServer{
	p := new(PlayerServer)
	p.store = store
	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}

// this is no longer needed due to the embedding of the handler interface in the
// PlayerServer type, operation that promotes the method ServeHTTP to the level
// of PlayerServers structure
//func (p *PlayerServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//	p.Handler.ServeHTTP(writer, request)
//}

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
	//leagueTable := p.getLeagueTable()
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())

}

func (p *PlayerServer) getLeagueTable() []Player {
	leagueTable := []Player{
		{"Chris", 20},
	}
	return leagueTable
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

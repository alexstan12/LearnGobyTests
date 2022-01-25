package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league []Player
}

func (s *StubPlayerStore) GetLeague() []Player{
	return s.league
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestRecordingAndRetrievingPostgres(t *testing.T) {
	db := connectToDB(t)
	defer db.Close()

	store := &PostgresPlayerStore{
		db,
		sync.Mutex{},
	}
	server := NewPlayerServer(store)

	t.Run("return not existing user score from DB", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest("james"))
		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("add new user to DB or increment existing one", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newPostWinRequest("linux"))
		got := response.Body.String()
		want := ""
		assertStatus(t, response.Code, http.StatusAccepted)
		assertResponseBody(t, got, want)
	})
	t.Run("increment user score in DB with multiple concurrent calls", func(t *testing.T) {
		wantedCount := 1000
		var wg sync.WaitGroup
		var response *httptest.ResponseRecorder
		wg.Add(wantedCount)
		for  i:= 0; i< wantedCount ; i++ {
			go func() {
				response = httptest.NewRecorder()
				server.ServeHTTP(response, newPostWinRequest("linux"))
				wg.Done()
			}()
		}
		wg.Wait()
		response = httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest("linux"))
		assertResponseBody(t, response.Body.String(),strconv.Itoa(wantedCount))
	})

}

func TestLeague(t *testing.T){
	store := &StubPlayerStore{
	}
	log.Printf("The current number of allocated logical CPUs are %d", runtime.NumCPU())
	server := NewPlayerServer(store)
	t.Run("it returns 200 on /league", func(t *testing.T) {
		response := httptest.NewRecorder()
		request,_ := http.NewRequest(http.MethodGet, "/league", nil)

		server.ServeHTTP(response, request)

		var got []Player
		err:= json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("it returns the table league as JSON", func(t *testing.T){
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store.league = wantedLeague

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)
	})
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
	}
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func getLeagueFromResponse(t *testing.T, response *httptest.ResponseRecorder) (league []Player) {
	err := json.NewDecoder(response.Body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
	}
	return league
}

// This function returns the DB object on which any subsequent actions are made.
// Closing of the db connection should be handled outside this function
func connectToDB(t *testing.T) *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s port=%d", user, password, dbname, sslmode, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Successfully connected to the db inside test")
	return db
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		winCalls: nil,
	}
	server := NewPlayerServer(&store)

	t.Run("return Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertResponseBody(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)

	})

	t.Run("return Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

		assertResponseBody(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Aurica")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		t.Logf("score is %s", response.Body.String())
		assertStatus(t, response.Code, http.StatusNotFound)
	})

	//t.Run("it returns accepted on POST", func(t *testing.T) {
	//
	//	request := newPostWinRequest("Pepper")
	//	response := httptest.NewRecorder()
	//
	//	server.ServeHTTP(response, request)
	//
	//	assertStatus(t, response.Code, http.StatusAccepted)
	//})

	t.Run("it records wins when POST", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}
		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner, got  %q want %q", store.winCalls[0], player)
		}
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertResponseBody(t, response.Body.String(), "3")
}

func newPostWinRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func assertResponseBody(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q, want %q", got, want)
	}
}

func assertStatus(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response code is wrong, got %d, want %d", got, want)
	}
}

func assertLeague(t *testing.T, got []Player, wantedLeague []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, wantedLeague) {
		t.Errorf("got %v want %v", got, wantedLeague)
	}
}

func newGetScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

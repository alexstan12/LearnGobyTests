package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
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
	server := PlayerServer{store: store}

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
	server := PlayerServer{store: store}
	t.Run("it returns 200 on /league", func(t *testing.T) {
		response := httptest.NewRecorder()
		request,_ := http.NewRequest(http.MethodGet, "/league", nil)
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
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
	server := &PlayerServer{
		store: &store,
	}

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
	server := PlayerServer{store}
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

func newGetScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

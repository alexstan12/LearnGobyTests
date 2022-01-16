package main

type PostgresPlayerStore struct {

}

func (p *PostgresPlayerStore) GetPlayerScore(name string) int {
	return 5
}

func (p *PostgresPlayerStore) RecordWin(name string) {

}

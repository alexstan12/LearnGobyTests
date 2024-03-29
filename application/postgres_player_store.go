package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

type PostgresPlayerStore struct {
	//store map[string]int
	db *sql.DB
	mu sync.Mutex
}

func (p *PostgresPlayerStore) GetLeague() []Player {
	var league []Player
	league, err := getLeagueFromDB(p.db)
	if err != nil {
		fmt.Printf("could not retrieve league from DB, err: %e", err)
	}
	return league
}

func (p *PostgresPlayerStore) GetPlayerScore(name string) int {
	score, err := getPlayerScoreFromDB(p.db, name)
	if err != nil {
		fmt.Printf("Could not retrieve player from DB, err:%e", err)
	}
	//fmt.Println(score)
	return score
}

func (p *PostgresPlayerStore) RecordWin(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	err := incrementPlayerScoreInDB(p.db, name)
	if err != nil {
		fmt.Printf("Could not retrieve player from DB, err:%e", err)
	}
}

func incrementPlayerScoreInDB(db *sql.DB, name string) error {
	score, err := getPlayerScoreFromDB(db, name)
	if err != nil {
		return err
	}
	insertPlayerStmt, err := db.Prepare("insert into player_store(name, score) VALUES ($1, $2) ON CONFLICT DO NOTHING;")
	if err != nil {
		return err
	}
	updateScoreStmt, err := db.Prepare("update player_store set score=$1 where name=$2")
	if err != nil {
		return err
	}
	if score == 0 {
		_, err = insertPlayerStmt.Exec(name, 1)
		if err != nil {
			return err
		}
	} else {
		score += 1
		_, err = updateScoreStmt.Exec(score, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPlayerScoreFromDB(db *sql.DB, name string) (score int, err error) {
	getPlayerQuery := `select score from player_store where name = $1;`
	rows, err := db.Query(getPlayerQuery, name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&score)
		if err != nil {
			return 0, err
		}
		log.Println(score)
	}
	err = rows.Err()
	if err != nil {
		return 0, err
	}
	return score, nil
}

func getLeagueFromDB(db *sql.DB) ([]Player, error) {
	getLeagueQuery := `select name, score from player_store order by score asc;`
	rows, err := db.Query(getLeagueQuery)
	if err != nil {
		return []Player{}, err
	}
	defer rows.Close()
	var league []Player
	for rows.Next() {
		var name string
		var wins int
		err := rows.Scan(&name, &wins)
		if err != nil {
			return []Player{}, err
		}
		player := Player{name, wins}
		log.Println(player)
		league = append(league, player)
	}
	err = rows.Err()
	if err != nil {
		return []Player{}, err
	}
	return league, nil
}

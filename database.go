package main

import (
	_ "github.com/mattn/go-sqlite3"

	"fmt"

	"database/sql"
)

var Database *sql.DB

// watchlist

var (
	WatchlistADDREPORT   *sql.Stmt
	WatchlistCOUNT       *sql.Stmt
	WatchlistUSERREPORTS *sql.Stmt
)

func initDB() {
	var err error
	Database, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		fmt.Println(err)
	}

	Database.Exec("CREATE TABLE IF NOT EXISTS watchlist (userID TEXT,userTag TEXT,userPFP TEXT, reason TEXT, originGuildID TEXT, originGuildName TEXT)")

	WatchlistADDREPORT, err = Database.Prepare("INSERT INTO watchlist VALUES (?,?,?,?,?,?)")

	WatchlistCOUNT, err = Database.Prepare("SELECT COUNT(DISTINCT userID) FROM watchlist")
	WatchlistUSERREPORTS, err = Database.Prepare(`SELECT (*) FROM watchlist WHERE userID = ?`)

	if err != nil {
		fmt.Println(err)
	}
}

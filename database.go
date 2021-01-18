package main

import (
	_ "github.com/mattn/go-sqlite3"

	"fmt"

	"database/sql"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var Database *sqlx.DB

// watchlist

var (
	WatchlistADDREPORT   *sql.Stmt
	WatchlistCOUNT       *sql.Stmt
	WatchlistUSERREPORTS *sqlx.Stmt
	WatchlistALL         *sqlx.Stmt
)

var (
	PermissionsRetrieve *sqlx.Stmt
	PermissionsSet      *sql.Stmt
)

type WatchlistEntree struct {
	UserID          string `db:"userID"`
	UserTag         string `db:"userTag"`
	UserPFP         string `db:"userPFP"`
	Reason          string `db:"reason"`
	OriginGuildID   string `db:"originGuildID"`
	OriginGuildName string `db:"originGuildName"`
}

type PermissionEntree struct {
	UserID         string `db:"userID"`
	WatchlistAdmin int    `db:"watchlistAdmin"`
	EditPermissions int    `db:"editPermissions"`
}

func initDB() {
	var err error
	Database, err = sqlx.Open("sqlite3", "./data.db")
	if err != nil {
		fmt.Println(err)
	}

	Database.Exec("CREATE TABLE IF NOT EXISTS watchlist (userID TEXT,userTag TEXT,userPFP TEXT, reason TEXT, originGuildID TEXT, originGuildName TEXT)")
	Database.Exec("CREATE TABLE IF NOT EXISTS internalPermissions (userID TEXT UNIQUE, watchlistAdmin INT, editPermissions INT)")

	// watchlist
	WatchlistADDREPORT, err = Database.Prepare("INSERT INTO watchlist VALUES (?,?,?,?,?,?)")
	WatchlistCOUNT, err = Database.Prepare("SELECT COUNT(DISTINCT userID) FROM watchlist")
	WatchlistUSERREPORTS, err = Database.Preparex(`SELECT * FROM watchlist WHERE userID = $1`)
	WatchlistALL, err = Database.Preparex("SELECT * FROM watchlist")
	// permissions
	PermissionsRetrieve, err = Database.Preparex("SELECT * FROM internalPermissions WHERE userID = $1")
	PermissionsSet, err = Database.Prepare("INSERT INTO internalPermissions VALUES (?,?,?)")

	if err != nil {
		fmt.Println(err)
	}
}

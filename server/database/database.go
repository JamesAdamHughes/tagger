package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var DBConn *sql.DB

func init() {
	db, err := sql.Open("sqlite3", "../dev.db")
	DBConn = db
	if err != nil {
		fmt.Println(err)
	}
	return
}

func Exec(query string, args ...interface{}) (result sql.Result, err error) {
	stmt, es := DBConn.Prepare(query)
	if es != nil {
		panic(es.Error())
	}

	result, err = stmt.Exec(args...)
	return
}

func Insert(query string, args ...interface{}) (result *sql.Row, err error) {
	stmt, es := DBConn.Prepare(query)
	if es != nil {
		panic(es.Error())
	}

	result = stmt.QueryRow(args...)
	return
}
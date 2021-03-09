package database

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DBConn *sqlx.DB

func init() {
	db, err := sqlx.Open("sqlite3", "../dev.db")
	DBConn = db
	if err != nil {
		fmt.Println(err)
	}
	return
}

func Exec(query string, args map[string]interface{}) (result sql.Result, err error) {
	result, err = DBConn.NamedExec(query, args)
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

func Select(query string, args map[string]interface{}) (result *sqlx.Rows, err error) {
	result, err = DBConn.NamedQuery(query, args)
	return
}

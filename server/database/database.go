package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DBConn *sqlx.DB

func init() {
	dbName := os.Getenv("DBNAME")
	dbPath := fmt.Sprintf("/go/src/tagger/%s.db", dbName)

	fmt.Printf("init connection to sqlite database %s\n", dbPath)

	db, err := sqlx.Open("sqlite3", dbPath)
	DBConn = db
	if err != nil {
		fmt.Println(err)
	}

	// test running a query on the DB. Sqlite3 will always find a db, doesn't mean there is anything in it
	_, err = Select("select * from tb_user_song_tags limit 1", map[string]interface{}{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database  - error: %s", err.Error()))
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

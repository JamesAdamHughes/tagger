package database

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
)

var db *sql.DB

func init() {
	db, err := sqlx.Open("sqlite3", "../dev.db")
	if err != nil {
		fmt.Println(err)
	}
}

func AddTag(songId string, tagId int64, userId *int64) {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err := db.Exec(
		"INSERT INTO tb_user_song_tags (fk_user_id, fk_song_id, fk_tag_id) VALUES ($1, $2, $3)",
		userId,
		songId,
		tagId)

	if err != nil {
		log.Fatal(err)
	}

	return
}
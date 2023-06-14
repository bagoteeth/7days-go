package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "gee.db")
	defer func() {
		_ = db.Close()
	}()
	_, _ = db.Exec("drop table if exists user;")
	_, _ = db.Exec("create table user(Name text);")
	res, err := db.Exec("insert into user(`Name`) values (?), (?)", "Tom", "Sam")
	if err == nil {
		affected, _ := res.RowsAffected()
		log.Println(affected)
	} else {
		log.Println(err)
	}

	row := db.QueryRow("select Name from user limit 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	} else {
		log.Println(err)
	}
}

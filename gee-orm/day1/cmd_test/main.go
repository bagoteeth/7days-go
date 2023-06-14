package main

import (
	geeorm "7days-go/gee-orm/day1"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	e, _ := geeorm.NewEngine("sqlite3", "gee.db")
	defer e.Close()
	s := e.NewSession()
	_, _ = s.Raw("drop table if exists user;").Exec()
	_, _ = s.Raw("create table user(Name text);").Exec()
	_, _ = s.Raw("create table user(Name text);").Exec()
	res, _ := s.Raw("insert into user(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := res.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}

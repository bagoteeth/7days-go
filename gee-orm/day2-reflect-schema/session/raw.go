package session

import (
	"7days-go/gee-orm/day2-reflect-schema/dialect"
	"7days-go/gee-orm/day2-reflect-schema/log"
	"7days-go/gee-orm/day2-reflect-schema/schema"
	"database/sql"
	"strings"
)

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	sql      strings.Builder
	sqlVars  []interface{}
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (r *Session) Clear() {
	r.sql.Reset()
	r.sqlVars = nil
}

func (r *Session) DB() *sql.DB {
	return r.db
}

func (r *Session) Raw(sql string, values ...interface{}) *Session {
	r.sql.WriteString(sql)
	r.sql.WriteString(" ")
	r.sqlVars = append(r.sqlVars, values...)
	return r
}

func (r *Session) Exec() (res sql.Result, err error) {
	defer r.Clear()
	log.Info(r.sql.String(), r.sqlVars)
	if res, err = r.DB().Exec(r.sql.String(), r.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (r *Session) QueryRow() *sql.Row {
	defer r.Clear()
	log.Info(r.sql.String(), r.sqlVars)
	return r.DB().QueryRow(r.sql.String(), r.sqlVars...)
}

func (r *Session) QueryRows() (rows *sql.Rows, err error) {
	defer r.Clear()
	log.Info(r.sql.String(), r.sqlVars)
	if rows, err = r.DB().Query(r.sql.String(), r.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

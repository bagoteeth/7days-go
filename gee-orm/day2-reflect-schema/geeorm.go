package geeorm

import (
	"7days-go/gee-orm/day2-reflect-schema/dialect"
	"7days-go/gee-orm/day2-reflect-schema/log"
	"7days-go/gee-orm/day2-reflect-schema/session"
	"database/sql"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
	}

	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (r *Engine) Close() {
	if err := r.db.Close(); err != nil {
		log.Error("Failed to close databases")
	}
	log.Info("Close database success")
}

func (r *Engine) NewSession() *session.Session {
	return session.New(r.db, r.dialect)
}

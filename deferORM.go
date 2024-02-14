package deferORM

import (
	"database/sql"
	"deferORM/log"
	"deferORM/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver string, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	e = &Engine{
		db: db,
	}
	log.PrintInfo("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close db")
	}
	log.PrintInfo("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db)
}

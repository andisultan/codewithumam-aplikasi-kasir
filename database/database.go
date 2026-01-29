package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(0)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to Supabase via transaction pooler")
	return db, nil
}

package controllers

import (
	"database/sql"
	"log"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_latihan_pbp_2")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

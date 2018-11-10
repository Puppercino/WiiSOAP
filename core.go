package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func StartSQL() {
	db, err := sql.Open("mysql",
		":"+"@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

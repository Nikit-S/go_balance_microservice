package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Database *sql.DB
	L        *log.Logger
}

var DB Database

func ConnectDB(lo *log.Logger) {
	DB.L = lo
	db, err := sql.Open("mysql", "barcher:1@tcp(mariadb:3306)/avito")
	if err != nil {
		DB.L.Println("Connection", err.Error())
	}
	DB.L.Printf("Successfully connected to avito database aft")
	DB.Database = db
}

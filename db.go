package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/sunshineplan/utils/database"
	"github.com/sunshineplan/utils/database/mysql"
)

var dbConfig database.Database
var db *sql.DB

func initMySQL() error {
	var config mysql.Config
	if err := meta.Get("mystocks_mysql", &config); err != nil {
		return err
	}
	dbConfig = &config
	return nil
}

func getDB() {
	var err error
	db, err = dbConfig.Open()
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

package main

import (
	"database/sql"

	"github.com/sunshineplan/utils/database"
	"github.com/sunshineplan/utils/database/mysql"
)

var dbConfig database.Database

func initMySQL() error {
	var config mysql.Config
	if err := meta.Get("mystocks_mysql", &config); err != nil {
		return err
	}
	dbConfig = &config
	return nil
}

func getDB() (*sql.DB, error) {
	return dbConfig.Open()
}

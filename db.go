package main

import (
	"database/sql"
	"time"

	"github.com/sunshineplan/utils/database"
	"github.com/sunshineplan/utils/database/mysql"
	"github.com/sunshineplan/utils/database/sqlite"
)

var dbConfig database.Database
var db *sql.DB

func initDB() (err error) {
	if local {
		dbConfig = &sqlite.Config{
			Path: joinPath(dir(self), "instance/mystocks.db"),
		}
	} else {
		var config mysql.Config
		if err = meta.Get("mystocks_mysql", &config); err != nil {
			return
		}
		dbConfig = &config
	}

	db, err = dbConfig.Open()
	if err != nil {
		return
	}

	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return
}

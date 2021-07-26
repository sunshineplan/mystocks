package main

import (
	"github.com/sunshineplan/utils"
	"github.com/sunshineplan/utils/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbConfig mongodb.Config
var collAccount *mongo.Collection
var collStock *mongo.Collection
var collFlows *mongo.Collection

func initDB() (err error) {
	if err = utils.Retry(func() error {
		return meta.Get("mystocks_mongo", &dbConfig)
	}, 3, 20); err != nil {
		return
	}

	var client *mongo.Client
	client, err = dbConfig.Open()
	if err != nil {
		return
	}

	database := client.Database(dbConfig.Database)

	collAccount = database.Collection("account")
	collStock = database.Collection("stock")
	collFlows = database.Collection("capitalflows")

	return
}

func test() error {
	if err := meta.Get("mystocks_mongo", &dbConfig); err != nil {
		return err
	}

	_, err := dbConfig.Open()
	return err
}

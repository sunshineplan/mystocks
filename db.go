package main

import (
	"errors"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/database/mongodb/driver"
)

var (
	mongoClient                             driver.Client
	accountClient, stockClient, flowsClient mongodb.Client
)

func initDB() (err error) {
	if mongoClient.Server == "" {
		return errors.New("MongoDB Server Address is required")
	}

	account, stock, flows := mongoClient, mongoClient, mongoClient
	account.Collection = "account"
	stock.Collection = "stock"
	flows.Collection = "capitalflows"
	accountClient, stockClient, flowsClient = &account, &stock, &flows

	if err = accountClient.Connect(); err != nil {
		return
	}
	if err = stockClient.Connect(); err != nil {
		return
	}
	return flowsClient.Connect()
}

func test() error {
	return initDB()
}

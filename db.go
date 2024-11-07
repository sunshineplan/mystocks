package main

import (
	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/database/mongodb/driver"
	"github.com/sunshineplan/utils/retry"
)

var accountClient, stockClient, flowsClient mongodb.Client

func initDB() (err error) {
	var apiClient driver.Client
	if err = retry.Do(func() error {
		return meta.Get("mystocks_mongo", &apiClient)
	}, 3, 20); err != nil {
		return err
	}

	account, stock, flows := apiClient, apiClient, apiClient
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

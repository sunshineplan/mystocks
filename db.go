package main

import (
	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/database/mongodb/api"
	"github.com/sunshineplan/utils"
)

var accountClient, stockClient, flowsClient mongodb.Client

func initDB() (err error) {
	var apiClient api.Client
	if err = utils.Retry(func() error {
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

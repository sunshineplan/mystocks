package main

import (
	"github.com/sunshineplan/database/mongodb/api"
	"github.com/sunshineplan/utils"
)

var accountClient, stockClient, flowsClient api.Client

func initDB() (err error) {
	var mongo api.Client
	if err = utils.Retry(func() error {
		return meta.Get("mystocks_mongo", &mongo)
	}, 3, 20); err != nil {
		return
	}

	accountClient, stockClient, flowsClient = mongo, mongo, mongo
	accountClient.Collection = "account"
	stockClient.Collection = "stock"
	flowsClient.Collection = "capitalflows"

	return
}

func test() error {
	var mongo api.Client
	return meta.Get("mystocks_mongo", &mongo)
}

package main

import (
	"context"
	"log"
	"time"

	"github.com/sunshineplan/stock"
	"github.com/sunshineplan/utils/cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var stockCache = cache.New(true)

func loadStocks(id interface{}, init bool) ([]stock.Stock, error) {
	if id == "" {
		return []stock.Stock{
			stock.Init("SSE", "000001"),
			stock.Init("SZSE", "399001"),
			stock.Init("SZSE", "399106"),
			stock.Init("SZSE", "399006"),
			stock.Init("SZSE", "399005"),
		}, nil
	}

	if !init {
		value, ok := stockCache.Get(id)
		if ok {
			return value.([]stock.Stock), nil
		}
	}

	ss, err := getStocks(id)
	if err != nil {
		return nil, err
	}

	stockCache.Set(id, ss, 1*time.Hour, func() interface{} {
		ss, err := getStocks(id)
		if err != nil {
			log.Print(err)
			return []stock.Stock{}
		}
		return ss
	})

	return ss, nil
}

func getStocks(id interface{}) ([]stock.Stock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collStock.Find(
		ctx, bson.M{"user": id}, options.Find().SetSort(bson.M{"seq": 1}))
	if err != nil {
		log.Println("Failed to query stocks:", err)
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var ss []stock.Stock
	var res []struct{ Index, Code string }
	if err := cursor.All(ctx, &res); err != nil {
		log.Println("Failed to get stocks:", err)
		return nil, err
	}
	for _, i := range res {
		ss = append(ss, stock.Init(i.Index, i.Code))
	}

	return ss, nil
}

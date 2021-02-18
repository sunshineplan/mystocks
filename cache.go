package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/sunshineplan/stock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type stocks struct {
	stocks     []stock.Stock
	expiration int64
}

type cache struct {
	sync.RWMutex
	data map[string]stocks
}

func (c *cache) get(id string) ([]stock.Stock, error) {
	if id == "" {
		return []stock.Stock{
			stock.Init("SSE", "000001"),
			stock.Init("SZSE", "399001"),
			stock.Init("SZSE", "399106"),
			stock.Init("SZSE", "399006"),
			stock.Init("SZSE", "399005"),
		}, nil
	}

	c.RLock()
	stocks, ok := c.data[id]
	c.RUnlock()

	if ok {
		if stocks.expiration != 0 && time.Now().UnixNano() > stocks.expiration {
			stocks.expiration = 0
			defer c.init(id)
		}
		return stocks.stocks, nil
	}

	return c.init(id)
}

func (c *cache) init(id string) ([]stock.Stock, error) {
	c.Lock()
	defer c.Unlock()

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

	c.data[id] = stocks{
		stocks:     ss,
		expiration: time.Now().Add(10 * time.Minute).UnixNano(),
	}

	return ss, nil
}

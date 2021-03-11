package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sunshineplan/gohttp"
	"github.com/sunshineplan/stock"
	"github.com/sunshineplan/stock/capitalflows/sector"
	"github.com/sunshineplan/utils/cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var stockCache = cache.New(false)
var flowsCache = cache.New(true)

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

	stockCache.Set(id, ss, 1*time.Hour, func() (interface{}, error) {
		return getStocks(id)
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

func loadFlows(date string) ([]sector.Chart, error) {
	value, ok := flowsCache.Get(date)
	if ok {
		return value.([]sector.Chart), nil
	}

	flows, err := getFlows(date)
	if err != nil {
		return nil, err
	}

	var duration time.Duration
	if date == "" {
		duration = time.Minute
	} else {
		duration = time.Hour
	}

	flowsCache.Set(date, flows, duration, nil)

	return flows, nil
}

func getFlows(date string) (flows []sector.Chart, err error) {
	if date != "" {
		date = strings.ReplaceAll(date, "-", "/")

		github := "https://raw.githubusercontent.com/sunshineplan/capital-flows-data/main/data/%s.json"
		jsdelivr := "https://cdn.jsdelivr.net/gh/sunshineplan/capital-flows-data/data/%s.json"

		rc := make(chan *gohttp.Response, 1)
		done := make(chan bool, 1)
		get := func(url string) {
			var mustReturn bool
			c := make(chan *gohttp.Response, 1)
			go func() { c <- gohttp.Get(fmt.Sprintf(url, date), nil) }()
			for {
				select {
				case ok := <-done:
					mustReturn = true
					if ok {
						return
					}
				case resp := <-c:
					if resp.Error != nil && !mustReturn {
						done <- false
						return
					}
					rc <- resp
					done <- true
					return
				}
			}
		}

		go get(github)
		go get(jsdelivr)

		resp := <-rc
		if resp.StatusCode == 404 {
			resp.Close()
			return
		}

		var tl []sector.TimeLine
		if err = resp.JSON(&tl); err != nil {
			return nil, err
		}

		for _, i := range tl {
			flows = append(flows, sector.TimeLine2Chart(i))
		}

		return
	}

	t := time.Now().In(time.FixedZone("CST", 8*60*60))
	date = fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())

	flows, err = sector.GetChart(date, collFlows)

	return
}

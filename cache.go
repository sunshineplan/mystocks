package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sunshineplan/database/mongodb/api"
	"github.com/sunshineplan/stock"
	"github.com/sunshineplan/stock/capitalflows/sector"
	"github.com/sunshineplan/utils/cache"
	"github.com/sunshineplan/utils/executor"
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

func getStocks(id interface{}) (ss []stock.Stock, err error) {
	var res []struct{ Index, Code string }
	if err = stockClient.Find(
		api.M{"user": id},
		&api.FindOpt{Sort: api.M{"seq": 1}},
		&res,
	); err != nil {
		log.Println("Failed to get stocks:", err)
		return
	}
	for _, i := range res {
		ss = append(ss, stock.Init(i.Index, i.Code))
	}

	return
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
		var result interface{}
		result, err = executor.ExecuteConcurrentArg(
			[]string{
				"https://raw.githubusercontent.com/sunshineplan/capital-flows-data/main/data/%s.json",
				"https://cdn.jsdelivr.net/gh/sunshineplan/capital-flows-data/data/%s.json",
			},
			func(url interface{}) (interface{}, error) {
				return http.Get(fmt.Sprintf(url.(string), strings.ReplaceAll(date, "-", "/")))
			},
		)
		if err != nil {
			return
		}

		resp, ok := result.(*http.Response)
		if !ok || resp == nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return
		}

		decoder := json.NewDecoder(resp.Body)
		var tl []sector.TimeLine
		if err = decoder.Decode(&tl); err != nil {
			return
		}

		for _, i := range tl {
			flows = append(flows, sector.TimeLine2Chart(i))
		}

		return
	}

	t := time.Now().In(time.FixedZone("CST", 8*60*60))
	date = fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())

	flows, err = sector.GetChart(date, &flowsClient)

	return
}

package main

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/gohttp"
	"github.com/sunshineplan/stock"
	"github.com/sunshineplan/stock/capitalflows/sector"
	"github.com/sunshineplan/utils/cache"
	"github.com/sunshineplan/utils/executor"
	"github.com/sunshineplan/workday"
	"github.com/sunshineplan/workday/apihubs"
	"github.com/sunshineplan/workday/timor"
)

var (
	stockCache    = cache.New(false)
	flowsCache    = cache.New(true)
	isTradingDate atomic.Value
)

func loadStocks(id any, init bool) ([]stock.Stock, error) {
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

	stockCache.Set(id, ss, 1*time.Hour, func() (any, error) {
		return getStocks(id)
	})

	return ss, nil
}

func getStocks(id any) (ss []stock.Stock, err error) {
	var res []struct{ Index, Code string }
	if err = stockClient.Find(
		mongodb.M{"user": id},
		&mongodb.FindOpt{Sort: mongodb.M{"seq": 1}},
		&res,
	); err != nil {
		svc.Println("Failed to get stocks:", err)
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
		var res any
		res, err = executor.ExecuteConcurrentArg(
			[]string{
				"https://raw.githubusercontent.com/sunshineplan/capital-flows-data/main/data/%s.json",
				"https://cdn.jsdelivr.net/gh/sunshineplan/capital-flows-data/data/%s.json",
				"https://fastly.jsdelivr.net/gh/sunshineplan/capital-flows-data/data/%s.json",
			},
			func(url string) (any, error) {
				return gohttp.Get(fmt.Sprintf(url, strings.ReplaceAll(date, "-", "/")), nil)
			},
		)
		if err != nil {
			return
		}

		resp, ok := res.(*gohttp.Response)
		if !ok || resp == nil {
			return
		}
		defer resp.Close()
		if resp.StatusCode == 404 {
			return
		}

		var timeline []sector.TimeLine
		if err = resp.JSON(&timeline); err != nil {
			return
		}
		for _, i := range timeline {
			flows = append(flows, sector.TimeLine2Chart(i))
		}
		return
	}

	t := time.Now().In(time.FixedZone("CST", 8*60*60))
	sectors, err := sector.GetSectors(fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()), flowsClient)
	if err != nil {
		return
	}
	return sectors.Charts(), nil
}

type tradingDate struct {
	date string
	is   bool
}

func loadTradingDate(t time.Time) (bool, error) {
	if res, ok := isTradingDate.Load().(tradingDate); ok && res.date == t.Format("2006-01-02") {
		return res.is, nil
	}
	is, err := getTradingDate(t)
	if err != nil {
		return false, err
	}
	isTradingDate.Store(tradingDate{t.Format("2006-01-02"), is})
	return is, nil
}

func getTradingDate(t time.Time) (bool, error) {
	if *workdayAPI != "" {
		is, err := workday.NewWorkdayAPI(*workdayAPI).IsWorkday(t)
		if err == nil {
			return is, nil
		}
	}
	return workday.IsWorkday(t, apihubs.API, timor.API)
}

package main

import (
	"sync"
	"time"

	"github.com/sunshineplan/stock"
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
	c.RLock()
	defer c.RUnlock()

	stocks, ok := c.data[id]
	if ok {
		if time.Now().UnixNano() > stocks.expiration {
			c.delete(id)
			defer c.init(id)
		}
		return stocks.stocks, nil
	}
	return c.init(id)
}

func (c *cache) init(id string) ([]stock.Stock, error) {
	c.Lock()
	defer c.Unlock()

	rows, err := db.Query(`SELECT idx, code FROM stock JOIN seq ON stock.user_id = seq.user_id AND stock.id = seq.stock_id
WHERE stock.user_id = ? ORDER BY seq`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ss []stock.Stock
	for rows.Next() {
		var index, code string
		if err := rows.Scan(&index, &code); err != nil {
			return nil, err
		}
		ss = append(ss, stock.Init(index, code))
	}

	c.data[id] = stocks{
		stocks:     ss,
		expiration: time.Now().Add(10 * time.Minute).UnixNano(),
	}

	return ss, nil
}

func (c *cache) delete(id string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.data[id]; ok {
		delete(c.data, id)
	}
}

package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/stock"
)

func myStocks(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		userID = 0
	}

	rows, err := db.Query(`SELECT idx, code FROM stock JOIN seq ON stock.user_id = seq.user_id AND stock.id = seq.stock_id
WHERE stock.user_id = ? ORDER BY seq`, userID)
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			restore("")
			c.String(501, "")
			return
		}
		log.Println("Failed to get all stocks:", err)
		c.String(500, "")
		return
	}
	defer rows.Close()
	var stocks []stock.Stock
	for rows.Next() {
		var index, code string
		if err := rows.Scan(&index, &code); err != nil {
			log.Println("Failed to scan all stocks:", err)
			c.String(500, "")
			return
		}
		stocks = append(stocks, stock.Init(index, code))
	}
	c.JSON(200, stock.Realtimes(stocks))
}

func indices(c *gin.Context) {
	indices := stock.Realtimes(
		[]stock.Stock{
			stock.Init("SSE", "000001"),
			stock.Init("SZSE", "399001"),
			stock.Init("SZSE", "399006"),
			stock.Init("SZSE", "399005"),
		})
	c.JSON(200, gin.H{"沪": indices[0], "深": indices[1], "创": indices[2], "中": indices[3]})
}

func getStock(c *gin.Context) {
	var r struct{ Index, Code, Q string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}

	s := stock.Init(r.Index, r.Code)

	if r.Q == "realtime" {
		realtime := s.GetRealtime()
		c.JSON(200, realtime)
		return
	} else if r.Q == "chart" {
		chart := s.GetChart()
		c.JSON(200, chart)
		return
	}
	c.String(400, "")
}

func getSuggest(c *gin.Context) {
	var r struct{ Keyword string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}
	c.JSON(200, stock.Suggests(r.Keyword))
}

func star(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.String(200, "0")
		return
	}
	refer := strings.Split(c.Request.Referer(), "/")
	index := refer[len(refer)-2]
	code := refer[len(refer)-1]

	if userID != nil {
		var exist string
		if err := db.QueryRow("SELECT idx FROM stock WHERE idx = ? AND code = ? AND user_id = ?",
			index, code, userID).Scan(&exist); err == nil {
			c.String(200, "1")
			return
		}
	}
	c.String(200, "0")
}

func doStar(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.String(200, "0")
		return
	}
	refer := strings.Split(c.Request.Referer(), "/")
	index := refer[len(refer)-2]
	code := refer[len(refer)-1]

	var r struct{ Action string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}

	if userID != nil {
		if r.Action == "unstar" {
			if _, err := db.Exec("DELETE FROM stock WHERE idx = ? AND code = ? AND user_id = ?",
				index, code, userID); err != nil {
				log.Println("Failed to unstar stock:", err)
				c.String(500, "")
				return
			}
		} else {
			if _, err := db.Exec("INSERT INTO stock (idx, code, user_id) VALUES (?, ?, ?)",
				index, code, userID); err != nil {
				log.Println("Failed to star stock:", err)
				c.String(500, "")
				return
			}
		}
		c.String(200, "1")
		return
	}
	c.String(200, "0")
}

func reorder(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.String(200, "0")
		return
	}

	var r struct{ New, Old string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}

	orig := strings.Split(r.Old, " ")
	dest := strings.Split(r.New, " ")

	ec := make(chan error, 1)
	var oldID, old, new int
	go func() {
		ec <- db.QueryRow(`SELECT id, seq FROM stock JOIN seq ON stock.user_id = seq.user_id AND stock.id = seq.stock_id
WHERE idx = ? AND code = ? AND stock.user_id = ?`, orig[0], orig[1], userID).Scan(&oldID, &old)
	}()
	if err := db.QueryRow(`SELECT seq FROM stock JOIN seq ON stock.user_id = seq.user_id AND stock.id = seq.stock_id
WHERE idx = ? AND code = ? AND stock.user_id = ?`,
		dest[0], dest[1], userID).Scan(&new); err != nil {
		log.Println("Failed to scan dest seq:", err)
		c.String(500, "")
		return
	}
	if err := <-ec; err != nil {
		log.Println("Failed to scan orig seq:", err)
		c.String(500, "")
		return
	}

	var err error
	if old > new {
		_, err = db.Exec("UPDATE seq SET seq = seq + 1 WHERE seq >= ? AND seq < ? AND user_id = ?",
			new, old, userID)
	} else {
		_, err = db.Exec("UPDATE seq SET seq = seq - 1 WHERE seq > ? AND seq <= ? AND user_id = ?",
			old, new, userID)
	}
	if err != nil {
		log.Println("Failed to update other seq:", err)
		c.String(500, "")
		return
	}
	if _, err := db.Exec("UPDATE seq SET seq = ? WHERE stock_id = ? AND user_id = ?",
		new, oldID, userID); err != nil {
		log.Println("Failed to update orig seq:", err)
		c.String(500, "")
		return
	}
	c.String(200, "1")
}

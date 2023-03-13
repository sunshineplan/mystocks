package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/stock"
)

func myStocks(c *gin.Context) {
	userID, _, err := getUser(c)
	if err != nil {
		svc.Println("Failed to get user:", err)
		c.String(500, "")
		return
	}

	stocks, err := loadStocks(userID, false)
	if err != nil {
		svc.Println("Failed to get all stocks:", err)
		c.String(500, "")
		return
	}

	c.JSON(200, stock.Realtimes(stocks))
}

func capitalFlows(c *gin.Context) {
	date, _ := c.GetQuery("date")

	flows, err := loadFlows(date)
	if err != nil {
		svc.Println("Failed to get flows chart:", err)
		c.String(500, "")
		return
	}

	c.JSON(200, flows)
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

func getRealtime(c *gin.Context) {
	var r struct{ Index, Code string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}
	c.JSON(200, stock.Init(r.Index, r.Code).GetRealtime())
}

func getChart(c *gin.Context) {
	var r struct{ Index, Code string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}
	c.JSON(200, stock.Init(r.Index, r.Code).GetChart())
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
	userID, _, err := getUser(c)
	if err != nil {
		svc.Println("Failed to get user:", err)
		c.String(500, "")
		return
	} else if userID == "" {
		c.String(200, "0")
		return
	}

	refer := strings.Split(c.Request.Referer(), "/")
	index := refer[len(refer)-2]
	code := refer[len(refer)-1]

	if n, err := stockClient.CountDocuments(mongodb.M{"index": index, "code": code, "user": userID}, nil); n > 0 {
		c.String(200, "1")
		return
	} else if err != nil {
		svc.Print(err)
	}
	c.String(200, "0")
}

func doStar(c *gin.Context) {
	userID, _, err := getUser(c)
	if err != nil {
		svc.Println("Failed to get user:", err)
		c.String(500, "")
		return
	} else if userID == "" {
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

	if r.Action == "unstar" {
		var s struct{ Seq int }
		if err := stockClient.FindOneAndDelete(mongodb.M{"index": index, "code": code, "user": userID}, nil, &s); err != nil {
			svc.Println("Failed to unstar stock:", err)
			c.String(500, "")
			return
		}

		if _, err := stockClient.UpdateMany(
			mongodb.M{"user": userID, "seq": mongodb.M{"$gt": s.Seq}},
			mongodb.M{"$inc": mongodb.M{"seq": -1}},
			nil,
		); err != nil {
			svc.Println("Failed to reorder after unstar stock:", err)
			c.String(500, "")
			return
		}
	} else {
		var s []struct{ Seq int }
		if err := stockClient.Find(
			mongodb.M{"user": userID},
			&mongodb.FindOpt{Sort: mongodb.M{"seq": -1}, Limit: 1},
			&s,
		); err != nil {
			svc.Println("Failed to get stocks:", err)
			c.String(500, "")
			return
		}

		var seq int
		if len(s) == 0 {
			seq = 1
		} else {
			seq = s[0].Seq + 1
		}

		res, err := stockClient.UpdateOne(
			struct {
				Index string `json:"index" bson:"index"`
				Code  string `json:"code" bson:"code"`
				User  string `json:"user" bson:"user"`
			}{index, code, userID},
			mongodb.M{"$setOnInsert": mongodb.M{"seq": seq}},
			&mongodb.UpdateOpt{Upsert: true},
		)
		if err != nil {
			svc.Println("Failed to star stock:", err)
			c.String(500, "")
			return
		}

		if res.MatchedCount == 1 {
			svc.Print("Stock already exists")
		}
	}

	loadStocks(userID, true)

	c.String(200, "1")
}

func reorder(c *gin.Context) {
	userID, _, err := getUser(c)
	if err != nil {
		svc.Println("Failed to get user:", err)
		c.String(500, "")
		return
	} else if userID == "" {
		c.String(200, "")
		return
	}

	var r struct{ New, Old string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}

	orig := strings.Split(r.Old, " ")
	dest := strings.Split(r.New, " ")

	code := 200
	if err := reorderStock(userID, orig, dest); err != nil {
		svc.Println("Failed to reorder stock:", err)
		code = 500
	}

	loadStocks(userID, true)

	c.String(code, "")
}

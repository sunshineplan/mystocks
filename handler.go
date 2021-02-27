package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/gohttp"
	"github.com/sunshineplan/stock"
	"github.com/sunshineplan/stock/capitalflows/sector"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func myStocks(c *gin.Context) {
	userID, _, err := getUser(c)
	if err != nil {
		log.Println("Failed to get user:", err)
		c.String(500, "")
		return
	}

	stocks, err := loadStocks(userID, false)
	if err != nil {
		log.Println("Failed to get all stocks:", err)
		c.String(500, "")
		return
	}

	c.JSON(200, stock.Realtimes(stocks))
}

func capitalFlows(c *gin.Context) {
	var date string
	if date, _ = c.GetQuery("date"); date != "" {
		url := fmt.Sprintf(
			"https://cdn.jsdelivr.net/gh/sunshineplan/capital-flows-data/data/%s.json",
			strings.ReplaceAll(date, "-", "/"),
		)

		var tl []sector.TimeLine
		if err := gohttp.Get(url, nil).JSON(&tl); err != nil {
			log.Println("Failed to get flows chart:", err)
			c.String(500, "")
			return
		}

		var flows []sector.Chart
		for _, i := range tl {
			flows = append(flows, sector.TimeLine2Chart(i))
		}

		c.JSON(200, flows)
		return
	}

	tz, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(tz)
	date = fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())

	flows, err := sector.GetChart(date, collFlows)
	if err != nil {
		log.Println("Failed to get flows chart:", err)
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
	userID, _, err := getUser(c)
	if err != nil {
		log.Println("Failed to get user:", err)
		c.String(500, "")
		return
	} else if userID == "" {
		c.String(200, "0")
		return
	}

	refer := strings.Split(c.Request.Referer(), "/")
	index := refer[len(refer)-2]
	code := refer[len(refer)-1]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := collStock.FindOne(
		ctx, bson.M{"index": index, "code": code, "user": userID}).Err(); err == nil {
		c.String(200, "1")
		return
	} else if err != mongo.ErrNoDocuments {
		log.Print(err)
	}
	c.String(200, "0")
}

func doStar(c *gin.Context) {
	userID, _, err := getUser(c)
	if err != nil {
		log.Println("Failed to get user:", err)
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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var s struct{ Seq int }
		if err := collStock.FindOneAndDelete(ctx,
			bson.M{"index": index, "code": code, "user": userID}).Decode(&s); err != nil {
			log.Println("Failed to delete stock:", err)
			c.String(500, "")
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if _, err := collStock.UpdateMany(ctx,
			bson.M{"user": userID, "seq": bson.M{"$gt": s.Seq}},
			bson.M{"$inc": bson.M{"seq": -1}},
		); err != nil {
			log.Println("Failed to reorder after delete stock:", err)
			c.String(500, "")
			return
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := collStock.Find(
			ctx, bson.M{"user": userID}, options.Find().SetSort(bson.M{"seq": -1}).SetLimit(1))
		if err != nil {
			log.Println("Failed to query stocks:", err)
			c.String(500, "")
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var s []struct{ Seq int }
		if err := cursor.All(ctx, &s); err != nil {
			log.Println("Failed to get stocks:", err)
			c.String(500, "")
			return
		}

		var seq int
		if len(s) == 0 {
			seq = 1
		} else {
			seq = s[0].Seq + 1
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := collStock.UpdateOne(
			ctx,
			bson.D{
				{Key: "index", Value: index},
				{Key: "code", Value: code},
				{Key: "user", Value: userID},
			},
			bson.M{"seq": seq},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			log.Println("Failed to add stock:", err)
			c.String(500, "")
			return
		}

		if res.MatchedCount == 1 {
			log.Print("Stock already exists")
		}
	}

	loadStocks(userID, true)

	c.String(200, "1")
}

func reorder(c *gin.Context) {
	userID, _, err := getUser(c)
	if err != nil {
		log.Println("Failed to get user:", err)
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
		log.Println("Failed to reorder stock:", err)
		code = 500
	}

	loadStocks(userID, true)

	c.String(code, "")
}
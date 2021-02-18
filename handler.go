package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/stock"
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

	stocks, err := stockCache.get(userID)
	if err != nil {
		log.Println("Failed to get all stocks:", err)
		c.String(500, "")
		return
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

	stockCache.init(userID)

	c.String(200, "1")
}

func reorder(c *gin.Context) {
	userID, _, err := getUser(c)
	if err != nil {
		log.Println("Failed to get user:", err)
		c.String(500, "")
		return
	} else if userID == "" {
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

	if err := reorderStock(userID, orig, dest); err != nil {
		log.Println("Failed to reorder stock:", err)
		c.String(500, "")
		return
	}

	stockCache.init(userID)

	c.String(200, "1")
}

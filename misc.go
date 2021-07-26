package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func addUser(username string) {
	log.Print("Start!")
	if err := initDB(); err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}

	username = strings.TrimSpace(strings.ToLower(username))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collAccount.InsertOne(ctx, bson.D{
		{Key: "username", Value: username},
		{Key: "password", Value: "123456"},
		{Key: "uid", Value: username},
	})
	if err != nil {
		log.Fatal(err)
	}
	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatal("Failed to get last insert id.")
	}
	if _, err := collStock.InsertMany(ctx, []interface{}{
		bson.D{
			{Key: "index", Value: "SSE"},
			{Key: "code", Value: "000001"},
			{Key: "user", Value: objectID.Hex()},
			{Key: "seq", Value: 1},
		},
		bson.D{
			{Key: "index", Value: "SZSE"},
			{Key: "code", Value: "399001"},
			{Key: "user", Value: objectID.Hex()},
			{Key: "seq", Value: 2},
		},
		bson.D{
			{Key: "index", Value: "SZSE"},
			{Key: "code", Value: "399106"},
			{Key: "user", Value: objectID.Hex()},
			{Key: "seq", Value: 3},
		},
		bson.D{
			{Key: "index", Value: "SZSE"},
			{Key: "code", Value: "399005"},
			{Key: "user", Value: objectID.Hex()},
			{Key: "seq", Value: 4},
		},
		bson.D{
			{Key: "index", Value: "SZSE"},
			{Key: "code", Value: "399006"},
			{Key: "user", Value: objectID.Hex()},
			{Key: "seq", Value: 5},
		},
	}); err != nil {
		log.Fatal(err)
	}
	log.Print("Done!")
}

func deleteUser(username string) {
	log.Print("Start!")
	if err := initDB(); err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}

	username = strings.TrimSpace(strings.ToLower(username))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collAccount.DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		log.Fatalln("Failed to delete user:", err)
	} else if res.DeletedCount == 0 {
		log.Fatalf("User %s does not exist.", username)
	}
	log.Print("Done!")
}

func reorderStock(userID interface{}, orig, dest []string) error {
	var origStock, destStock struct {
		ID  primitive.ObjectID `bson:"_id"`
		Seq int
	}

	c := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		c <- collStock.FindOne(
			ctx, bson.M{"index": orig[0], "code": orig[1], "user": userID}).Decode(&origStock)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := collStock.FindOne(
		ctx, bson.M{"index": dest[0], "code": dest[1], "user": userID}).Decode(&destStock); err != nil {
		return err
	}
	if err := <-c; err != nil {
		return err
	}

	var filter, update bson.M
	if origStock.Seq > destStock.Seq {
		filter = bson.M{"user": userID, "seq": bson.M{"$gte": destStock.Seq, "$lt": origStock.Seq}}
		update = bson.M{"$inc": bson.M{"seq": 1}}
	} else {
		filter = bson.M{"user": userID, "seq": bson.M{"$gt": origStock.Seq, "$lte": destStock.Seq}}
		update = bson.M{"$inc": bson.M{"seq": -1}}
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := collStock.UpdateMany(ctx, filter, update); err != nil {
		log.Println("Failed to reorder stock:", err)
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := collStock.UpdateOne(
		ctx, bson.M{"_id": origStock.ID}, bson.M{"$set": bson.M{"seq": destStock.Seq}}); err != nil {
		log.Println("Failed to reorder stock:", err)
		return err
	}

	return nil
}

func backup(file string) {
	log.Print("Start!")
	if err := initDB(); err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}

	if err := dbConfig.Backup(file); err != nil {
		log.Fatal(err)
	}
	log.Print("Done!")
}

func restore(file string) {
	log.Print("Start!")
	if _, err := os.Stat(file); err != nil {
		log.Fatalln("File not found:", err)
	}

	if err := initDB(); err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}

	if err := dbConfig.Restore(file); err != nil {
		log.Fatal(err)
	}
	log.Print("Done!")
}

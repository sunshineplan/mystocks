package main

import (
	"log"
	"strings"

	"github.com/sunshineplan/database/mongodb/api"
)

func addUser(username string) {
	log.Print("Start!")
	if err := initDB(); err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}

	username = strings.TrimSpace(strings.ToLower(username))

	insertedID, err := accountClient.InsertOne(
		struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Uid      string `json:"uid"`
		}{username, "123456", username},
	)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := stockClient.InsertMany(
		[]struct {
			Index string `json:"index"`
			Code  string `json:"code"`
			User  string `json:"user"`
			Seq   int    `json:"seq"`
		}{
			{"SSE", "000001", insertedID, 1},
			{"SZSE", "399001", insertedID, 2},
			{"SZSE", "399106", insertedID, 3},
			{"SZSE", "399005", insertedID, 4},
			{"SZSE", "399006", insertedID, 5},
		},
	); err != nil {
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

	deletedCount, err := accountClient.DeleteOne(api.M{"username": username})
	if err != nil {
		log.Fatalln("Failed to delete user:", err)
	} else if deletedCount == 0 {
		log.Fatalf("User %s does not exist.", username)
	}
	log.Print("Done!")
}

func reorderStock(userID interface{}, orig, dest []string) error {
	var origStock, destStock struct {
		ID  string `json:"_id"`
		Seq int
	}

	c := make(chan error, 1)
	go func() {
		c <- stockClient.FindOne(api.M{"index": orig[0], "code": orig[1], "user": userID}, nil, &origStock)
	}()
	if err := stockClient.FindOne(api.M{"index": dest[0], "code": dest[1], "user": userID}, nil, &destStock); err != nil {
		return err
	}
	if err := <-c; err != nil {
		return err
	}

	var filter, update api.M
	if origStock.Seq > destStock.Seq {
		filter = api.M{"user": userID, "seq": api.M{"$gte": destStock.Seq, "$lt": origStock.Seq}}
		update = api.M{"$inc": api.M{"seq": 1}}
	} else {
		filter = api.M{"user": userID, "seq": api.M{"$gt": origStock.Seq, "$lte": destStock.Seq}}
		update = api.M{"$inc": api.M{"seq": -1}}
	}

	if _, err := stockClient.UpdateMany(filter, update, nil); err != nil {
		log.Println("Failed to reorder stock:", err)
		return err
	}

	if _, err := stockClient.UpdateOne(api.M{"_id": origStock.ID}, api.M{"$set": api.M{"seq": destStock.Seq}}, nil); err != nil {
		log.Println("Failed to reorder stock:", err)
		return err
	}

	return nil
}

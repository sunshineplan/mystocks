package main

import (
	"fmt"
	"strings"

	"github.com/sunshineplan/database/mongodb"
)

func addUser(username string) error {
	svc.Print("Start!")
	if err := initDB(); err != nil {
		return err
	}

	username = strings.TrimSpace(strings.ToLower(username))

	insertedID, err := accountClient.InsertOne(
		struct {
			Username string `json:"username" bson:"username"`
			Password string `json:"password" bson:"password"`
			Uid      string `json:"uid" bson:"uid"`
		}{username, "123456", username},
	)
	if err != nil {
		return err
	}

	type stock struct {
		Index string `json:"index" bson:"index"`
		Code  string `json:"code" bson:"code"`
		User  string `json:"user" bson:"user"`
		Seq   int    `json:"seq" bson:"seq"`
	}
	if _, err := stockClient.InsertMany(
		[]stock{
			{"SSE", "000001", insertedID.(mongodb.ObjectID).Hex(), 1},
			{"SZSE", "399001", insertedID.(mongodb.ObjectID).Hex(), 2},
			{"SZSE", "399106", insertedID.(mongodb.ObjectID).Hex(), 3},
			{"SZSE", "399005", insertedID.(mongodb.ObjectID).Hex(), 4},
			{"SZSE", "399006", insertedID.(mongodb.ObjectID).Hex(), 5},
		},
	); err != nil {
		return err
	}
	svc.Print("Done!")
	return nil
}

func deleteUser(username string) error {
	svc.Print("Start!")
	if err := initDB(); err != nil {
		return err
	}

	username = strings.TrimSpace(strings.ToLower(username))

	deletedCount, err := accountClient.DeleteOne(mongodb.M{"username": username})
	if err != nil {
		return err
	} else if deletedCount == 0 {
		return fmt.Errorf("user %s does not exist", username)
	}
	svc.Print("Done!")
	return nil
}

func reorderStock(userID string, orig, dest []string) error {
	var origStock, destStock struct {
		ID  mongodb.OID `json:"_id" bson:"_id"`
		Seq int
	}

	c := make(chan error, 1)
	go func() {
		c <- stockClient.FindOne(mongodb.M{"index": orig[0], "code": orig[1], "user": userID}, nil, &origStock)
	}()
	if err := stockClient.FindOne(mongodb.M{"index": dest[0], "code": dest[1], "user": userID}, nil, &destStock); err != nil {
		return err
	}
	if err := <-c; err != nil {
		return err
	}

	var filter, update mongodb.M
	if origStock.Seq > destStock.Seq {
		filter = mongodb.M{"user": userID, "seq": mongodb.M{"$gte": destStock.Seq, "$lt": origStock.Seq}}
		update = mongodb.M{"$inc": mongodb.M{"seq": 1}}
	} else {
		filter = mongodb.M{"user": userID, "seq": mongodb.M{"$gt": origStock.Seq, "$lte": destStock.Seq}}
		update = mongodb.M{"$inc": mongodb.M{"seq": -1}}
	}

	if _, err := stockClient.UpdateMany(filter, update, nil); err != nil {
		svc.Println("Failed to reorder stock:", err)
		return err
	}

	id, _ := stockClient.ObjectID(origStock.ID.Hex())
	if _, err := stockClient.UpdateOne(
		mongodb.M{"_id": id},
		mongodb.M{"$set": mongodb.M{"seq": destStock.Seq}},
		nil,
	); err != nil {
		svc.Println("Failed to reorder stock:", err)
		return err
	}

	return nil
}

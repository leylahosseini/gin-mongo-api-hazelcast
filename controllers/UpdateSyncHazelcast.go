package controllers

import (
	"context"
	"fmt"
	"gin-mongo-api-hazelcast/configs"
	hazelcastMe "gin-mongo-api-hazelcast/hazelcast"
	"time"

	"github.com/hazelcast/hazelcast-go-client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func UpdateSyncHazelcast() gin.HandlerFunc {
// 	return func(c *gin.Context) {
func UpdateSyncHazelcast() {
	// Connect to Hazelcast
	var client_hazel *hazelcast.Client = hazelcastMe.Hazelcast_Connect()

	// Connect to MongoDB
	var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
	defer client_hazel.Shutdown(context.TODO())

	matchStage := bson.D{{"$match", bson.D{{"operationType", "insert"}}}}
	//	matchStage2 := bson.D{{"$match", bson.D{{"fullDocument","fullDocument"}}}}
	opts := options.ChangeStream().SetMaxAwaitTime(1 * time.Second)

	changeStream, err := userCollection.Watch(context.Background(), mongo.Pipeline{matchStage}, opts)

	//changeStream, err := client.Database(db.EntityDbName).Collection(constants.AdSaleModelDbCollectionName).Watch(context.Background(), mongo.Pipeline{})

	if err != nil {
		fmt.Println("Watch hazelcast ERROR .....", err)
	}

	for changeStream.Next(context.Background()) {
		event := changeStream.Current
		//fmt.Println("EVEEEEEEENT ................", event.Lookup("fullDocument"))

		data := event.Lookup("fullDocument")
		data_id := data.Document().Lookup("_id")

		fmt.Println("PRINT FULLDOCUMENT..........................", data)
		fmt.Println("PRINT FULLDOCUMENT..........................", data_id)

		mapName := "mapUsers"
		myMap, err := client_hazel.GetMap(context.Background(), mapName)
		if err != nil {
			fmt.Println("ERRORRRRRRRRRR in Get hazzelcast map ", err)
		}
		err = myMap.Set(context.Background(), data_id, data)
		if err != nil {
			fmt.Println("ERRROR", err)
		}
		fmt.Println("Added record with ID %v to Hazelcast\n", data_id)
		// //}
	}

	//	c.JSON(200, "OKKKKKKKKKKK")
}

//		fmt.Println("FULLLLLLLLLLLLLLL .......", Full)

//event.FullDocument["fullDocument"].(string)

//var FullDocument bson.Raw = event
// FullDocument = changeStream.Current
//fmt.Println(FullDocument.Values())

// fmt.Println(FullDocument["fullDocument"])
//event["fullDocument"]
//	fmt.Println(changeStream.Next(context.Background()))
//	var event bson.M
//operationType := event.operationType
//if OperationType == "insert" {

//	record := event.Lookup("fullDocument")
//fmt.Println("SHOWWWWWWWWWWWWWWWWWWWWWWW RECORD", record)

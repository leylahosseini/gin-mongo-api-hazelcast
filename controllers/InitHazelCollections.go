package controllers

import (
	"context"
	"fmt"
	"gin-mongo-api-hazelcast/configs"
	hazelcastMe "gin-mongo-api-hazelcast/hazelcast"
	"log"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/serialization"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//var client, errSession = connectionManager.GetSessionDb()

//func InitHazelcast_AdSales_advertiser_id() gin.HandlerFunc {

func InitHazelcast() {
	//	return func(c *gin.Context) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//	var users []models.User
	var documents []bson.M
	//defer cancel()
	var client_hazel *hazelcast.Client = hazelcastMe.Hazelcast_Connect()
	var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
	defer client_hazel.Shutdown(context.TODO())
	//Collection := client.Database(golangAPI).Collection(users)

	results, err := userCollection.Find(context.TODO(), bson.M{})

	//fmt.Println(re)
	if err != nil {

		fmt.Println("Error Find Collection ...")
	}

	if err = results.All(context.TODO(), &documents); err != nil {
		log.Fatal(err)
	}

	mapUsers, err := client_hazel.GetMap(context.TODO(), "mapUsers")

	if err != nil {
		//log.Fatal(err)
		fmt.Println("ERRORRRRRRRRRRRRR     could not Create  mapUsers in Hazelcast  ")
	}

	for _, doc := range documents {

		json, _ := bson.MarshalExtJSON(doc, true, false)
		mapUsers.Put(context.TODO(), doc["_id"], serialization.JSON(json))

		//	fmt.Println("kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk", json)

	}

	// // load all data hazelcast and show
	re2, err := mapUsers.GetEntrySet(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("EChOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO INITilize MAP is corrct ........", re2[len(re2)-1])

}

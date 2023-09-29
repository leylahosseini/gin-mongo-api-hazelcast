package controllers

import (
	"context"
	"fmt"
	"gin-mongo-api-hazelcast/configs"
	hazelcastMe "gin-mongo-api-hazelcast/hazelcast"
	"gin-mongo-api-hazelcast/models"
	"gin-mongo-api-hazelcast/responses"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hazelcast/hazelcast-go-client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Name:     user.Name,
			Location: user.Location,
			Title:    user.Title,
		}

		result, err := userCollection.InsertOne(ctx, newUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		fmt.Println(userId)
		//	var user models.User
		defer cancel()

		// hazelcast connect client to cluster

		var client_hazel *hazelcast.Client = hazelcastMe.Hazelcast_Connect()

		//config.Logger.Level = logger.InfoLevel

		defer client_hazel.Shutdown(ctx)

		// fmt.Println(client)
		myMap, _ := client_hazel.GetMap(ctx, "mapUsers")

		//	fmt.Println(myMap.LoadAllWithoutReplacing(ctx, "64cbce8cd7d359532f88c52f"))
		re2, err := myMap.Get(ctx, userId)

		fmt.Println(re2)

		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, re2)
		/////////////////////////

		// objId, _ := primitive.ObjectIDFromHex(userId)

		// err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		// 	return
		// }

		// c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}

func EditAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

		result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated user details
		var updatedUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}})
	}
}

func DeleteAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}},
		)
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//	var users []models.User
		var documents []bson.M
		defer cancel()

		results, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if err = results.All(ctx, &documents); err != nil {
			log.Fatal(err)
		}

		// // hazelcast connect client to cluster

		var client_hazel *hazelcast.Client = hazelcastMe.Hazelcast_Connect()
		if err != nil {
			log.Fatal(err)
		}
		defer client_hazel.Shutdown(ctx)

		//fmt.Println(client)
		mapusers, err := client_hazel.GetMap(ctx, "mapUsers")
		if err != nil {
			log.Fatal(err)
		}

		// load all data hazelcast and show
		re2, err := mapusers.GetEntrySet(ctx)

		//fmt.Println(mapusers.Get(ctx, "64cbce8cd7d359532f88c52f"))
		//	fmt.Println(mapusers.LoadAllWithoutReplacing(ctx, doc["_id"])
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(200, re2)

		//reading from the db in an optimal way
		defer results.Close(ctx)

	}
}

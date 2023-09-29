package main

import (
	"fmt"
	"gin-mongo-api-hazelcast/configs"
	"gin-mongo-api-hazelcast/controllers"
	"gin-mongo-api-hazelcast/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(router)

	// go rutin hazelcast init and map and AutoSync
	fmt.Println("TEST GO RUTIN ................................")
	go controllers.InitHazelcast()
	go controllers.UpdateSyncHazelcast()

	router.Run("localhost:8000")
}

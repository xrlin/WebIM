package main

import (
	"github.com/xrlin/WebIM/server/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/user/token", controllers.UserToken)
		// Register user
		api.POST("/users", controllers.CreateUser)
	}
	router.Run(":8080")
}
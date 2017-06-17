package main

import (
	"github.com/xrlin/WebIM/server/controllers"
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/middlewares"
)

func main() {
	router := gin.Default()
	router.Use(middlewares.CORS())
	api := router.Group("/api")
	{
		api.POST("/user/token", controllers.UserToken)
		// Register user
		api.POST("/users", controllers.CreateUser)
	}
	router.Run(":8080")
}
package routes

import (
	"github.com/xrlin/WebIM/server/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/controllers"
)

func RouterEngin() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORS())
	api := router.Group("/api")
	{
		api.POST("/user/token", controllers.UserToken)
		// Register user
		api.POST("/users", controllers.CreateUser)
	}
	return router
}

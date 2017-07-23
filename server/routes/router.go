package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/controllers"
	"github.com/xrlin/WebIM/server/middlewares"
)

func RouterEngin() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORS())
	api := router.Group("/api")
	{
		api.POST("/user/token", middlewares.Auth(), controllers.UserToken)
		// Register user
		api.POST("/users", controllers.CreateUser)
	}
	ws := router.Group("/ws")
	//ws.Use(middlewares.Auth())
	ws.GET("/chat", controllers.Chat)
	return router
}

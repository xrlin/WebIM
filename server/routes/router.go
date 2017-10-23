package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/controllers"
	"github.com/xrlin/WebIM/server/middlewares"
)

func RouterEngine() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORS())
	api := router.Group("/api")
	{
		api.POST("/user/token", middlewares.Auth(), controllers.UserToken)
		api.POST("/user/info", middlewares.Auth(), controllers.GetUserInfo)
		api.GET("/user/rooms", middlewares.Auth(), controllers.GetRecentRooms)
		api.POST("/user/rooms", middlewares.Auth(), controllers.CreateRoom)
		api.PUT("/user/avatar", middlewares.Auth(), controllers.UpdateAvatar)
		// Register user
		api.POST("/users", controllers.CreateUser)
		api.POST("/friends", middlewares.Auth(), controllers.AddFriend)
		api.GET("/friends", middlewares.Auth(), controllers.GetFriends)
		api.GET("/users/search", controllers.SearchUsers)
		api.DELETE("/rooms/:roomID/leave", middlewares.Auth(), controllers.LeaveRoom)
		api.GET("/messages/unread", middlewares.Auth(), controllers.GetUnreadOfflineMessages)
		api.DELETE("/messages/ack", middlewares.Auth(), controllers.AckReceive)
		api.POST("/friendship/apply", middlewares.Auth(), controllers.FriendApplication)
		api.POST("/friendship/check", middlewares.Auth(), controllers.CheckFriendApplication)
		api.POST("/notifications/read", middlewares.Auth(), controllers.AckReadFriendApplications)

		// qiniu
		api.POST("/qiniu/uptoken", middlewares.Auth(), controllers.UploadToken)
	}
	ws := router.Group("/ws")
	ws.Use(middlewares.Auth())
	ws.GET("/chat", controllers.Chat)
	return router
}

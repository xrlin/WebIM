package main

import (
	"github.com/gin-gonic/gin"
)

func RouterEngine() *gin.Engine {
	router := gin.Default()
	router.Use(CORS())
	api := router.Group("/api")
	{
		api.POST("/user/token", Auth(), UserToken)
		api.POST("/user/info", Auth(), GetUserInfo)
		api.GET("/user/rooms", Auth(), GetRecentRooms)
		api.POST("/user/rooms", Auth(), CreateRoom)
		api.PUT("/user/avatar", Auth(), UpdateAvatar)
		// Register user
		api.POST("/users", CreateUserHandler)
		api.POST("/friends", Auth(), AddFriend)
		api.GET("/friends", Auth(), GetFriends)
		api.GET("/users/search", SearchUsers)
		api.DELETE("/rooms/:roomID/leave", Auth(), LeaveRoom)
		api.GET("/messages/unread", Auth(), GetUnreadOfflineMessages)
		api.DELETE("/messages/ack", Auth(), AckReceive)
		api.POST("/friendship/apply", Auth(), FriendApplication)
		api.POST("/friendship/check", Auth(), CheckFriendApplication)
		api.POST("/notifications/read", Auth(), AckReadFriendApplications)

		// qiniu
		api.POST("/qiniu/uptoken", Auth(), UploadToken)
	}
	ws := router.Group("/ws")
	ws.Use(Auth())
	ws.GET("/chat", Chat)
	return router
}

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LeaveRoom(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	var roomID int
	var err error
	if roomID, err = strconv.Atoi(c.Param("roomID")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	user := userObj.(*User)
	if room, err := LeaveRoomService(hub, user, roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "leave room" + room.RoomName() + "success."})
	}
}

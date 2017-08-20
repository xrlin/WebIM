package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/models"
	"github.com/xrlin/WebIM/server/services"
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
	user := userObj.(*models.User)
	if room, err := services.LeaveRoom(hub, user, roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "leave room" + room.RoomName() + "success."})
	}
}

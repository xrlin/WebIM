package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateRoom(c *gin.Context) {
	var params CreateRoomParam
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	user, _ := c.Get("user")
	if room, err := CreateRoomService(hub, user.(*User), params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"room": room})
	}
}

type updateRoomParam struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func UpdateRoom(c *gin.Context) {
	var param updateRoomParam
	c.BindJSON(&param)
	user, _ := c.Get("user")
	roomID, _ := strconv.Atoi(c.Param("roomID"))
	room, err := UpdateRoomService(roomID, *user.(*User), param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"room": room})
}

func LeaveRoom(c *gin.Context) {
	user, _ := c.Get("user")
	roomID, _ := strconv.Atoi(c.Param("roomID"))
	if room, err := LeaveRoomService(hub, user.(*User), roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "leave room" + room.RoomName() + "success."})
	}
}

func JoinRoom(c *gin.Context) {
	user, _ := c.Get("user")
	roomId, _ := strconv.Atoi(c.Param("roomID"))
	room, err := JoinRoomService(user.(*User), roomId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"room": room})
}

func SearchRooms(c *gin.Context) {
	name := c.Query("name")
	rooms, err := SearchRoomService(name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

type inviteParam struct {
	UserIds []int `json:"userIds" binding:"required"`
}

func InviteUserToRoom(c *gin.Context) {
	user, _ := c.Get("user")
	roomID, _ := strconv.Atoi(c.Param("roomID"))
	var param inviteParam
	c.BindJSON(&param)
	err := InviteUserToRoomService(roomID, param.UserIds, *user.(*User))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func GetRoomInfo(c *gin.Context) {
	roomID, _ := strconv.Atoi(c.Param("roomID"))
	room, err := GetRoomInfoService(roomID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"room": room})
}

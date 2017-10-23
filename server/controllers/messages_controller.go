package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/models"
	"github.com/xrlin/WebIM/server/services"
	"net/http"
)

func GetUnreadOfflineMessages(c *gin.Context) {
	obj, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := obj.(*models.User)
	if messages, err := services.GetUnreadMessages(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		msgDetails := make([]models.MessageDetail, 0)
		for _, msg := range messages {
			msgDetails = append(msgDetails, msg.GetDetails())
		}
		c.JSON(http.StatusOK, gin.H{"messages": msgDetails})
	}
}

type Ack struct {
	MessageIds []uint `json:"message_ids"`
}

// When user read the messages, call this to delete the message(off-line) stored in server
func AckReceive(c *gin.Context) {
	obj, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var ack Ack
	if err := c.BindJSON(&ack); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user := obj.(*models.User)
	if err := services.DeleteUnreadMessages(user, ack.MessageIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "Successfully ack to server"})
	}
}

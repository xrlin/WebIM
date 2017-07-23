package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/database"
	"github.com/xrlin/WebIM/server/models"
	"github.com/xrlin/WebIM/server/services"
	"log"
)

var hub *services.Hub

func Chat(c *gin.Context) {
	u := &models.User{}
	database.DBConn.First(u)
	c.Set("user", u)
	userObj, ok := c.Get("user")
	if !ok {
		log.Fatalln("No user exist in context(Chat controller).")
		return
	}
	user := userObj.(*models.User)
	client, err := services.NewClient(hub, user, c.Writer, c.Request)
	if err != nil {
		log.Fatalln("Create client failed!(Chat controller).")
		return
	}
	hub.Register <- client

	go client.Read()
	go client.Write()

}

func init() {
	hub = services.NewHub()
	go hub.Run()
	go services.MonitorAndDeliverMessages(hub)
}

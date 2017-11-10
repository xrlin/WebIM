package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

var hub *Hub

func Chat(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		log.Fatalln("No user exist in context(Chat controller).")
		return
	}
	user := userObj.(*User)
	client, err := NewClient(hub, user, c.Writer, c.Request)
	if err != nil {
		log.Fatalln("Create client failed!(Chat controller).")
		return
	}
	hub.Register <- client

	go client.Read()
	go client.Write()

}

func init() {
	hub = NewHub()
	go hub.Run()
	go MonitorAndDeliverMessages(hub)
}

package services

import (
	"github.com/xrlin/WebIM/server/models"
	"log"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Rooms map[string][]*Client

	// Inbound messages from the clients.
	Messages chan models.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string][]*Client),
		Messages:   make(chan models.Message, 512),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.addClient(client)
		case client := <-hub.Unregister:
			hub.removeClient(client)
		case message := <-hub.Messages:
			hub.deliver(message)
		}
	}
}

func (hub *Hub) addClient(client *Client) {
	// Single user has a single room for itself
	userRoomName := client.user.UserRoomName()
	hub.Rooms[userRoomName] = []*Client{client}
	for _, name := range client.user.RoomNames() {
		hub.Rooms[name] = append(hub.Rooms[name], client)
	}
}

func (hub *Hub) removeClient(client *Client) {
	client.Close()
	delete(hub.Rooms, client.user.UserRoomName())
	for _, name := range client.user.RoomNames() {
		clientIdx := len(hub.Rooms[name])
		for idx, client := range hub.Rooms[name] {
			if client == client {
				clientIdx = idx
				break
			}
		}
		hub.Rooms[name] = append(hub.Rooms[name][0:clientIdx-1], hub.Rooms[name][clientIdx+1:]...)
	}
}

func (hub *Hub) deliver(message models.Message) {
	switch message.MsgType {
	case models.SingleMessage:
		if message.UserId != 0 {
			hub.deliverMsgToRoom(message.UserRoomName(), message)
		}
	case models.RoomMessage:
		hub.deliverMsgToRoom(message.RoomName(), message)
	}
}

func (hub *Hub) deliverMsgToRoom(room string, message models.Message) {
	for _, client := range hub.Rooms[room] {
		log.Println("Send messages to client %v", client)
		client.send <- message
	}
}

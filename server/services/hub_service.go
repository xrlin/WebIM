package services

import (
	"fmt"
	"github.com/xrlin/WebIM/server/database"
	"github.com/xrlin/WebIM/server/models"
	"github.com/xrlin/WebIM/server/utils"
	"log"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Store WebSocket clients
	clients []*Client

	// Registered clients.
	Rooms map[string][]*Client

	// Inbound messages from the clients.
	Messages chan models.MessageDetail

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Update/create room request, add more clients to room
	UpdateRoom chan *models.Room

	// User leave chat room forever
	LeaveRoom chan *leaveRoomParam
}

type leaveRoomParam struct {
	user *models.User
	room *models.Room
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		Rooms:      make(map[string][]*Client),
		Messages:   make(chan models.MessageDetail, 512),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		UpdateRoom: make(chan *models.Room, 128),
		LeaveRoom:  make(chan *leaveRoomParam),
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
		case room := <-hub.UpdateRoom:
			hub.updateRoom(room)
		case leaveRoomParam := <-hub.LeaveRoom:
			hub.leaveRoom(leaveRoomParam.room, leaveRoomParam.user)
		}
	}
}

// User leave existing room, remove user from room's clients,
// if room has no clients after removed the user then remove it from queue.
func (hub *Hub) leaveRoom(room *models.Room, user *models.User) {
	clients := hub.Rooms[room.RoomName()]
	for idx, client := range clients {
		if client.user.ID == user.ID {
			if idx+1 < len(clients) {
				clients = append(clients[:idx], clients[idx+1:]...)
			} else {
				clients = clients[:idx]
			}
			break
		}
	}
	if len(clients) == 0 {
		delete(hub.Rooms, room.RoomName())
		return
	}
	hub.Rooms[room.RoomName()] = clients
	uuid, _ := utils.GenerateUUID()
	var msgContent string
	msgContent = fmt.Sprintf("用户%v离开了房间", user.Name)
	msg := models.Message{UUID: uuid, Content: msgContent, MsgType: models.SystemMessage, RoomId: room.ID}
	hub.Messages <- msg.GetDetails()
}

// Add new users' clients to room. To prevent redundant messages send to user,
// the room passed to this function hte Users must only contains the new users of room.
func (hub *Hub) updateRoom(room *models.Room) {
	for _, user := range room.Users {
		if client, ok := getUserClient(&user, hub.clients); ok {
			room_clients := hub.Rooms[room.RoomName()]
			hub.Rooms[room.RoomName()] = append(room_clients, client)
		}
	}
}

func getUserClient(user *models.User, clients []*Client) (*Client, bool) {
	for _, client := range clients {
		if client.user.ID == user.ID {
			return client, true
		}
	}
	return nil, false
}

func removeDuplicatedClients(clients []*Client) []*Client {
	remark := make(map[int]bool)
	results := []*Client{}
	var id int
	for _, client := range clients {
		id = int(client.user.ID)
		if remark[id] {
			continue
		}
		results = append(results, client)
		remark[id] = true
	}
	return results
}

func (hub *Hub) addClient(client *Client) {
	// Single user has a single room for itself
	hub.clients = append(hub.clients, client)
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
		rooms := hub.Rooms[name]
		if len(rooms) == 0 {
			return
		}
		if clientIdx+1 < len(rooms) {
			copy(rooms[clientIdx:], rooms[clientIdx+1:])
		}
		hub.Rooms[name] = rooms[:len(hub.Rooms[name])-1]
	}
}

func (hub *Hub) deliver(messageDetail models.MessageDetail) {
	message := messageDetail.Message
	fmt.Println(message.MsgType)
	fmt.Println(models.SingleMessage)
	fmt.Println("Rooms", hub.Rooms)
	hub.deliverMsgToRoom(message.RoomName(), message)
	switch message.MsgType {
	case models.SingleMessage, models.FriendshipMessage:
		hub.deliverMsgToUser(message.UserId, message)
	case models.RoomMessage:
		fmt.Println("Rooms", hub.Rooms)
		if message.UserId != 0 {
			hub.deliverMsgToUser(message.UserId, message)
			return
		}
		hub.deliverMsgToRoom(message.RoomName(), message)
	}
}

func (hub *Hub) deliverMsgToUser(userID uint, message models.Message) {
	client := selectClientByUserID(hub.clients, userID)
	if client == nil {
		log.Printf("No client with userID %v now! will save offline message.", userID)
		log.Printf("Check to save offline message %#v.", message)
		err := database.DBConn.FirstOrCreate(&message, message).Error
		if err != nil {
			log.Fatalf("Check exist and save message %#v failed with error %v", message, err.Error())
		}
		return
	}
	log.Printf("Will send message to user.Message details: %#v", message.GetDetails())
	client.send <- message.GetDetails()
}

func selectClientByUserID(clients []*Client, id uint) *Client {
	for _, client := range clients {
		if client.user.ID == id {
			return client
		}
	}
	return nil
}

func (hub *Hub) deliverMsgToRoom(room string, message models.Message) {
	log.Println("Deleiver message to room: ", room)
	log.Println("Now the room has clients", hub.Rooms[room])
	sendUsers := make(map[uint]bool)

	var users []models.User
	var r models.Room
	database.DBConn.Where("id = ?", message.RoomId).Find(&r)
	database.DBConn.Model(&r).Related(&users, "Users")

	for _, client := range hub.Rooms[room] {
		log.Println("Send messages to client", client.user.Name)
		client.send <- message.GetDetails()
		sendUsers[client.user.ID] = true
	}
	fmt.Printf("room_id: %v, room: %v, users: %v", message.RoomId, r, users)
	for _, user := range users {
		if sendUsers[user.ID] {
			continue
		}
		// If user off-line save the message
		message.UserId = user.ID
		SaveOfflineMessage(message)
	}
}

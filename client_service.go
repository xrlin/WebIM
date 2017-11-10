package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub  *Hub
	user *User
	conn *websocket.Conn
	send chan MessageDetail
}

func (client *Client) Read() {
	defer func() {
		client.hub.Unregister <- client
		client.conn.Close()
	}()
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msg := new(Message)
		err := client.conn.ReadJSON(msg)
		log.Printf("message with room_id %v, will be sended to clident %v, user: %v", msg.RoomId, client, client.user.ID)
		if err != nil {
			log.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		PushMessage(*msg)
	}
}

func (client *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case msgDetail, ok := <-client.send:
			log.Printf("Message to client %#v", msgDetail)
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.conn.WriteJSON(msgDetail); err != nil {
				msgDetail.UserId = client.user.ID
				SaveOfflineMessage(msgDetail.Message)
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (client *Client) Close() {
	client.conn.Close()
}

func NewClient(hub *Hub, user *User, w http.ResponseWriter, r *http.Request) (*Client, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	client := &Client{hub: hub, user: user, conn: conn, send: make(chan MessageDetail, 256)}
	return client, nil
}

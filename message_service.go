package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin/json"
	"log"
	"qiniupkg.com/x/errors.v7"
)

const MessageList string = "messages"
const PopTimeout int = 3

func closeConn(conn redis.Conn) {
	conn.Close()
}

func DeliverMessage(msg Message, topic string, user *User) error {
	if err := checkMessage(msg, topic, user); err != nil {
		return err
	}
	client, err := GetMQTTClient()
	if err != nil {
		return err
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Printf("connect status %s", client.IsConnected())
	log.Printf("payload %#v\n", payload)
	token := client.Publish(topic, 1, false, payload)
	//token := client.Subscribe(topic, byte(0), nil )
	token.Wait()
	log.Printf("connect %s", client.IsConnected())
	log.Printf("connect %s", token.Error())
	client.Disconnect(250)
	return token.Error()
}

func checkMessage(msg Message, topic string, user *User) error {
	if msg.From != user.ID {
		return errors.New("forbidden")
	}
	if !hasPrivilegeToSend(topic, user) {
		return errors.New("forbidden")
	}
	return nil
}

// TODO
func hasPrivilegeToSend(topic string, user *User) bool {
	return true
}

func PushMessage(msg Message) error {
	fmt.Println(fmt.Sprintf("Message push type %v", msg))
	redisConn := RedisPool.Get()

	defer closeConn(redisConn)

	msgText, err := MarshalMessage(msg)
	if err != nil {
		return err
	}
	reply, err := redis.Int(redisConn.Do("LPUSH", MessageList, msgText))
	log.Printf("Reply with %d", reply)
	return err
}

func BRPopMessage() (Message, error) {
	var msg Message
	redisConn := RedisPool.Get()

	defer closeConn(redisConn)

	reply, err := redis.Values(redisConn.Do("BRPOP", MessageList, PopTimeout))
	if err != nil {
		return msg, err
	}
	var key, msgText string
	_, err = redis.Scan(reply, &key, &msgText)
	if err != nil {
		return msg, err
	}
	msg, err = UnMarshalMessage(msgText)
	return msg, err
}

func SaveOfflineMessage(msg Message) error {
	log.Printf("SaveOfflineMessage")
	return CreateMessage(&msg)
}

func MonitorAndDeliverMessages(hub *Hub) {
	for {
		message, err := BRPopMessage()
		if err == nil {
			hub.Messages <- message.GetDetails()
		}
	}
}

func GetUnreadMessages(user *User) ([]*Message, error) {
	if len(user.Rooms) == 0 {
		db.Model(user).Related(&user.Rooms, "Rooms")
	}
	var room_ids []uint
	for _, room := range user.Rooms {
		room_ids = append(room_ids, room.ID)
	}
	messages := make([]*Message, 0)
	err := db.Where("user_id = ? AND read <> true", user.ID).Find(&messages).Error
	return messages, err
}

func DeleteUnreadMessages(user *User, message_ids []uint) error {
	var messages []Message
	db.Where("id IN (?)", message_ids).Find(&messages)
	// check all messages belong to user
	for _, message := range messages {
		if message.UserId != user.ID {
			return fmt.Errorf("User %v can only modify messages of self", user.ID)
		}
	}
	return db.Where("id IN (?)", message_ids).Delete(Message{}).Error
}

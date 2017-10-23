package services

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/xrlin/WebIM/server/database"
	"github.com/xrlin/WebIM/server/models"
	"github.com/xrlin/WebIM/server/utils"
	"log"
)

const MessageList string = "messages"
const PopTimeout int = 3

func closeConn(conn redis.Conn) {
	conn.Close()
}

func PushMessage(msg models.Message) error {
	fmt.Println(fmt.Sprintf("Message push type %v", msg))
	redisConn := database.RedisPool.Get()

	defer closeConn(redisConn)

	msgText, err := utils.MarshalMessage(msg)
	if err != nil {
		return err
	}
	reply, err := redis.Int(redisConn.Do("LPUSH", MessageList, msgText))
	log.Printf("Reply with %d", reply)
	return err
}

func BRPopMessage() (models.Message, error) {
	var msg models.Message
	redisConn := database.RedisPool.Get()

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
	msg, err = utils.UnMarshalMessage(msgText)
	return msg, err
}

func SaveOfflineMessage(msg models.Message) error {
	log.Printf("SaveOfflineMessage")
	return models.CreateMessage(&msg)
}

func MonitorAndDeliverMessages(hub *Hub) {
	for {
		message, err := BRPopMessage()
		if err == nil {
			hub.Messages <- message.GetDetails()
		}
	}
}

func GetUnreadMessages(user *models.User) ([]*models.Message, error) {
	if len(user.Rooms) == 0 {
		database.DBConn.Model(user).Related(&user.Rooms, "Rooms")
	}
	var room_ids []uint
	for _, room := range user.Rooms {
		room_ids = append(room_ids, room.ID)
	}
	messages := make([]*models.Message, 0)
	err := database.DBConn.Where("user_id = ? AND read <> true", user.ID).Find(&messages).Error
	return messages, err
}

func DeleteUnreadMessages(user *models.User, message_ids []uint) error {
	var messages []models.Message
	database.DBConn.Where("id IN (?)", message_ids).Find(&messages)
	// check all messages belong to user
	for _, message := range messages {
		if message.UserId != user.ID {
			return fmt.Errorf("User %v can only modify messages of self", user.ID)
		}
	}
	return database.DBConn.Where("id IN (?)", message_ids).Delete(models.Message{}).Error
}

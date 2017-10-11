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

func SaveOfflineMessage(msg models.Message) {
	models.CreateMessage(&msg)
}

func MonitorAndDeliverMessages(hub *Hub) {
	for {
		message, err := BRPopMessage()
		if err == nil {
			hub.Messages <- message
		}
	}
}

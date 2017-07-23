package services

import (
	"testing"
	"github.com/xrlin/WebIM/server/models"
	"math/rand"
	"github.com/xrlin/WebIM/server/database"
	"os"
	"fmt"
	"time"
)

func setup() {
	redisConn := database.RedisPool.Get()
	defer func() {
		redisConn.Close()
	}()
	redisConn.Do("DEL", MessageList)
}

func TestPushMessageSuccess(t *testing.T) {
	msg := models.Message{RoomId: rand.Int()}
	if err := PushMessage(msg); err != nil {
		t.Errorf("Push message to list failed with error: %s", err)
	}
}

func TestBRPopMessageSuccess(t *testing.T) {
	msg := models.Message{UUID: fmt.Sprintf("%v-%v", time.Now().Unix(), 1), RoomId: rand.Int()}
	PushMessage(msg)
	replyMsg, err := BRPopMessage()
	if err != nil {
		t.Errorf("Pop message to list failed with error: %s", err)
	}
	if replyMsg.UUID != msg.UUID {
		t.Errorf("Poped message incorrect. Expected %v but %v instead!", msg, replyMsg)
	}
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

package models

import (
	"fmt"
	"github.com/xrlin/WebIM/server/database"
	"testing"
	"time"
)

func TestMessage_RoomName(t *testing.T) {
	msg := Message{RoomId: 3}
	if msg.RoomName() != "room_3" {
		t.Fail()
	}
}

func TestMessage_UserRoomName(t *testing.T) {
	msg := Message{UserId: 3}
	fmt.Println(msg.User)
	if msg.UserRoomName() != "room_user_3" {
		t.Fail()
	}
}

func TestCreateMessage(t *testing.T) {
	// Should success
	msg := Message{UUID: fmt.Sprintf("%v_%v", time.Now().Unix(), 1), RoomId: 1, UserId: 1, MsgType: SingleMessage, Content: "test"}
	defer func() {
		database.DBConn.Delete(&msg)
	}()
	if err := database.DBConn.Create(&msg).Error; err != nil {
		t.Errorf("Create msg %v failed! error: %s", msg, err)
	}
	// Should failed
	messagesFailed := []Message{{UUID: fmt.Sprintf("%v_%v", time.Now().Unix(), 2), RoomId: 1, UserId: 1, Content: "test"},
		{UUID: fmt.Sprintf("%v_%v", time.Now().Unix(), 3), RoomId: 1, UserId: 1, MsgType: SingleMessage},
		{UUID: fmt.Sprintf("%v_%v", time.Now().Unix(), 4), UserId: 1, MsgType: SingleMessage, Content: "test"},
		{UUID: fmt.Sprintf("%v_%v", time.Now().Unix(), 5), RoomId: 1, MsgType: SingleMessage, Content: "test"},
		{RoomId: 1, UserId: 1, MsgType: SingleMessage, Content: "test"}}
	for _, msgFailed := range messagesFailed {
		if err := database.DBConn.Create(&msgFailed).Error; err == nil {
			t.Errorf("Message with %v should failed", msgFailed)
		}
	}
	defer func() {
		for _, msgFailed := range messagesFailed {
			database.DBConn.Delete(&msgFailed)
		}
	}()
}

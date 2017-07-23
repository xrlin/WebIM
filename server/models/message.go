package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/xrlin/WebIM/server/database"
)

const (
	_ = iota
	SingleMessage
	RoomMessage
)

type Message struct {
	gorm.Model
	UUID    string `gorm:"not null;unique;column:uuid" json:"uuid"`
	RoomId  int    `gorm:"not null;index" json:"room_id" binding:"required"`
	Room    Room   `json:"-"`
	UserId  int    `gorm:"not null;index" json:"user_id"`
	User    User   `json:"-"`
	MsgType int    `json:"msg_type"`
	Content string `json:"content"`
}

func (msg *Message) RoomName() string {
	return fmt.Sprintf("room_%v", msg.RoomId)
}

func (msg *Message) UserRoomName() string {
	return fmt.Sprintf("room_user_%v", msg.UserId)
}

func CreateMessage(msg *Message) error {
	return database.DBConn.Create(msg).Error
}

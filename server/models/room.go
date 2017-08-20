package models

import (
	"fmt"
	"time"
)

const (
	MultiRoom = iota
	FriendRoom
)

type Room struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Name      string     `json:"name"`
	// 主要用于离线消息查询
	Messages []Message

	RoomType int `gorm:"not null;default:0" json:"room_type"`

	Users []User `gorm:"many2many:user_rooms;" json:"users"`
}

func (room *Room) RoomName() string {
	roomName := fmt.Sprintf("room_%v", room.ID)
	return roomName
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin/json"
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
	Name      string     `json:"name" gorm:"not null;default:''"`
	Avatar    string     `json:"avatar"`
	// 主要用于离线消息查询
	Messages []Message

	RoomType int `gorm:"not null;default:0" json:"room_type"`

	Users []User `gorm:"many2many:user_rooms;" json:"users"`
}

func (room *Room) AvatarURL() string {
	if room.Avatar == "" {
		return QiniuCfg.FileDomain + "/" + "default_group.jpg"
	}
	return QiniuCfg.FileDomain + "/" + room.Avatar
}

func (room *Room) RoomName() string {
	roomName := fmt.Sprintf("room_%v", room.ID)
	return roomName
}

func (room Room) MarshalJSON() ([]byte, error) {
	users, _ := GetRoomUsers(room)
	room.Users = users
	type RoomDetail Room
	detail := RoomDetail(room)
	detail.Avatar = room.AvatarURL()
	return json.Marshal(detail)
}

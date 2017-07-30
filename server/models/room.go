package models

import "github.com/jinzhu/gorm"

type Room struct {
	gorm.Model
	Name string
	// 主要用于离线消息查询
	Messages []Message

	Users []User `gorm:"many2many:user_rooms;"`
}

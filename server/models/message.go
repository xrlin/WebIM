package models

import (
	"errors"
	"fmt"
	"github.com/xrlin/WebIM/server/database"
	"reflect"
	"time"
)

const (
	_ = iota
	SingleMessage
	RoomMessage
)

type Message struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UUID      string    `gorm:"not null;unique;column:uuid" json:"uuid"`
	RoomId    int       `gorm:"not null;index" json:"room_id" binding:"required"`
	Room      Room      `json:"-"`
	UserId    int       `gorm:"not null;index" json:"user_id"`
	User      User      `json:"-"`
	MsgType   int       `gorm:"not null" json:"msg_type"`
	Content   string    `gorm:"not null" json:"content"`
}

func (msg *Message) BeforeSave() (err error) {
	err = checkBlank(msg, "UUID")
	fields := []string{"UUID", "Content", "UserId", "RoomId", "MsgType"}
	for _, field := range fields {
		err = checkBlank(msg, field)
		if err != nil {
			break
		}
	}
	return
}

func checkBlank(obj interface{}, field string) (err error) {
	v := reflect.ValueOf(obj).Elem()
	v = v.FieldByName(field)
	if v.IsValid() == false {
		err = errors.New(fmt.Sprintf("Field %s doesn't exist!", field))
		return
	}
	switch v.Kind() {
	case reflect.Int:
		if v.Int() == reflect.Zero(v.Type()).Int() {
			err = errors.New(fmt.Sprintf("Column %s could not be zero value", field))
		}
	case reflect.String:
		if v.String() == reflect.Zero(v.Type()).String() {
			err = errors.New(fmt.Sprintf("Column %s could not be zero value", field))
		}
	}
	return
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

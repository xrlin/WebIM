package models

import (
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
	RoomId    int       `gorm:"index" json:"room_id"`
	Room      Room      `json:"-"`
	// Id of the user the message will send to, if zero the message if not for a certain user.
	// When save offline message, user_id is required.
	UserId int  `gorm:"not null;index" json:"user_id"`
	User   User `json:"-"`
	// Id of the user that sends the message
	FromUser int    `gorm:"not null;index" json:"from_user" binding:"required"`
	MsgType  int    `gorm:"not null" json:"msg_type" binding:"required"`
	Content  string `gorm:"not null" json:"content"`
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
		err = fmt.Errorf("Field %s doesn't exist", field)
		return
	}
	switch v.Kind() {
	case reflect.Int:
		if v.Int() == reflect.Zero(v.Type()).Int() {
			err = fmt.Errorf("Column %s could not be zero value", field)
		}
	case reflect.String:
		if v.String() == reflect.Zero(v.Type()).String() {
			err = fmt.Errorf("Column %s could not be zero value", field)
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

package main

import (
	"fmt"
	"reflect"
	"time"
)

const (
	_ = iota
	SingleMessage
	RoomMessage
	SingleImageMessage
	RoomImageMessage
	SingleMusicMessage
	RoomMusicMessage
	FriendshipMessage
	SystemMessage
)

type Message struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UUID      string    `gorm:"not null;unique;column:uuid" json:"uuid"`
	RoomId    uint      `gorm:"index" json:"room_id"`
	Room      Room      `json:"-"`
	// Id of the user the message will send to, if zero the message if not for a certain user.
	// When save offline message, user_id is required.
	UserId uint `gorm:"not null;index" json:"user_id"`
	User   User `json:"-"`
	// Id of the user that sends the message
	FromUser uint   `gorm:"not null;index" json:"from_user"`
	MsgType  int    `gorm:"not null" json:"msg_type"`
	Content  string `gorm:"not null" json:"content"`
	Checked  bool   `sql:"DEFAULT:false" json:"checked"`
	Read     bool   `sql:"DEFAULT:false" json:"read"`
	Topic    string `json:"topic" binding:"required"`
	From     uint   `json:"from" binding:"required"`
	Payload  string `json:"payload" binding:"required"`
}

type MessageDetail struct {
	Message
	SourceUser User `json:"source_user"`
	TargetUser User `json:"target_user"`
}

// Callback before save the record. Most of the case, should not call the method manually.
func (msg *Message) BeforeSave() (err error) {
	fields := []string{"UUID", "UserId", "MsgType"}
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
	case reflect.Uint:
		if v.Uint() == reflect.Zero((v.Type())).Uint() {
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

func (msg Message) GetDetails() MessageDetail {
	var sourceUser, targetUser User
	db.Where("id = ?", msg.FromUser).Find(&sourceUser)
	db.Where("id = ?", msg.UserId).Find(&targetUser)
	return MessageDetail{msg, sourceUser, targetUser}
}

func CreateMessage(msg *Message) error {
	return db.Create(msg).Error
}

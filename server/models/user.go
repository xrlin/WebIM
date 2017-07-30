package models

import (
	"errors"
	"fmt"
	"github.com/xrlin/WebIM/server/database"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string     `gorm:"not null;unique_index:idx_name_deleted_at;primary_key"`
	DeletedAt    *time.Time `sql:"index" gorm:"unique_index:idx_name_deleted_at"`
	Password     string     `gorm:"-"`
	PasswordHash string     `gorm:"not null"`
	// 用于查询离线消息
	Messages []Message

	Rooms []Room `gorm:"many2many:user_rooms"`
}

// Room name just for user itself
func (u *User) UserRoomName() string {
	userId := u.ID
	// Single user has a single room for itself
	return fmt.Sprintf("room_user_%v", userId)
}

// Room names to specify the chat rooms shared with other users
func (u *User) RoomNames() []string {
	// Find rooms if haven't done before
	if len(u.Rooms) == 0 {
		database.DBConn.Model(u).Related(&u.Rooms, "Rooms")
	}
	roomNames := []string{}
	for _, room := range u.Rooms {
		roomName := fmt.Sprintf("room_%v", room.ID)
		roomNames = append(roomNames, roomName)
	}
	return roomNames
}

func FindUserByName(name string) *User {
	var user User
	database.DBConn.Where("name = ?", name).First(&user)
	// No record found
	if user.ID == 0 {
		return nil
	}
	return &user
}

/* Create user from a pointer to user, if successed update the pointer to the User struct with id and other fields form database*/
func CreateUser(user *User) error {
	if !database.DBConn.NewRecord(user) {
		return errors.New(fmt.Sprintf("User with name %s has exist!", user.Name))
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)
	return database.DBConn.Create(user).Error
}

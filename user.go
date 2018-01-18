package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/json"
	"time"
)

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// Hash(file name) of avatar in cdn
	Avatar       string     `json:"avatar"`
	Name         string     `gorm:"not null;unique_index:idx_name_deleted_at;" json:"name"`
	DeletedAt    *time.Time `sql:"index" gorm:"unique_index:idx_name_deleted_at" json:"-"`
	Password     string     `gorm:"-" json:"-"`
	PasswordHash string     `gorm:"not null" json:"-"`
	// 用于查询离线消息
	Messages []Message `json:"-"`
	Friends  []User    `gorm:"many2many:friendship" json:"-"`

	Rooms []Room `gorm:"many2many:user_rooms" json:"-"`
}

func (u *User) AvatarUrl() string {
	if u.Avatar == "" {
		return "https://xrlin.github.io/assets/img/crown-logo.png"
	}
	return QiniuCfg.FileDomain + "/" + u.Avatar
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
		db.Model(u).Related(&u.Rooms, "Rooms")
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
	db.Where("name = ?", name).First(&user)
	// No record found
	if user.ID == 0 {
		return nil
	}
	return &user
}

func FindUserById(id uint) *User {
	var user User
	db.Where("id = ?", id).First(&user)
	// No record found
	if user.ID == 0 {
		return nil
	}
	return &user
}

/* Create user from a pointer to user, if successed update the pointer to the User struct with id and other fields form database*/
func CreateUser(user *User) error {
	if !db.NewRecord(user) {
		return errors.New(fmt.Sprintf("User with name %s has exist!", user.Name))
	}
	passwordHash, err := getPasswordHash(user.Password)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash
	return db.Create(user).Error
}

func (u User) MarshalJSON() ([]byte, error) {
	u.Avatar = u.AvatarUrl()
	type Detail User
	detail := (Detail)(u)
	return json.Marshal(detail)
}

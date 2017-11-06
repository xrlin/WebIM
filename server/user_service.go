package main

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

/* Check if user info is valid with username and password */
func ValidateUser(userName, password string) (*User, bool) {
	user := FindUserByName(userName)
	if user == nil {
		return nil, false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, false
	}
	return user, true
}

/* Register and create a new user */
func RegisterUser(user *User) error {
	return CreateUser(user)
}

func UpdateAvatarService(user *User, avatar string) error {
	return db.Model(user).Update("avatar", avatar).Error
}

type room struct {
	Room
	UniqueName string `json:"unique_name"`
}

func GetUserRecentRooms(user *User) []room {
	rooms := make([]Room, 0)
	db.Preload("Users").Model(user).Related(&rooms, "Rooms")
	result := make([]room, 0)
	for _, r := range rooms {
		uniqueName := fmt.Sprintf("%s_%d", r.Name, r.ID)
		fmt.Println(r.Users)
		if r.RoomType == FriendRoom {
			// Set room name with the friend's name
			var friend User
			for _, u := range r.Users {
				if u.ID != user.ID {
					friend = u
				}
			}
			r.Users = []User{friend}
			r.Name = friend.Name
		}
		result = append(result, room{Room: r, UniqueName: uniqueName})
	}
	fmt.Printf("%#v", result)
	return result
}

func AddFriendService(hub *Hub, caller *User, friendID uint) (*User, *Room, error) {
	if caller.ID == friendID {
		return nil, nil, errors.New("Cannot add friend of yourself")
	}
	friends := GetUserFriends(*caller)
	for _, user := range friends {
		if user.ID == friendID {
			return nil, nil, errors.New("Already has friendship")
		}
	}
	friend := &User{}
	if err := db.First(friend, friendID).Error; err != nil {
		return nil, nil, err
	}
	room := Room{RoomType: FriendRoom, Users: []User{*caller, *friend}}
	if err := db.Create(&room).Error; err != nil {
		return nil, nil, err
	}
	room.Name = friend.Name
	hub.UpdateRoom <- &room
	return friend, &room, nil
}

func GetUserFriends(user User) []*User {
	friends := make([]*User, 0)

	rawSql := "SELECT * FROM users INNER JOIN userRooms ON userRooms.user_id = users.id WHERE users.id!=?  AND users.deleted_at IS NULL AND (userRooms.room_id IN (SELECT id FROM rooms INNER JOIN userRooms ON userRooms.room_id = rooms.id AND userRooms.user_id=? WHERE rooms.room_type = ? AND rooms.deleted_at IS NULL));"
	db.Raw(rawSql, user.ID, user.ID, FriendRoom).Scan(&friends)
	log.Printf("%#v", friends)
	return friends
}

func SearchUsersByName(name string) []*User {
	users := make([]*User, 0)
	db.Where("name LIKE ?", "%"+name+"%").Find(&users)
	log.Println("Search users with name ", name)
	log.Println(users)
	return users
}

// An application to ask for a friendship relation with another user
func AddFriendApplication(hub *Hub, fromUser User, toUserID uint) error {
	uuid, err := GenerateUUID()
	if err != nil {
		return err
	}
	toUser := User{ID: toUserID}
	err = db.First(&toUser).Error
	if err != nil {
		return err
	}
	msg := Message{UUID: uuid, FromUser: fromUser.ID, UserId: toUserID, MsgType: FriendshipMessage}
	log.Printf("Add friend msg %#v", msg)
	if err = SaveOfflineMessage(msg); err != nil {
		return err
	}
	hub.Messages <- msg.GetDetails()
	return nil
}

func PassFriendApplication(hub *Hub, applicationMsgUUID string) (*Room, error) {
	var msg Message
	if err := db.Where("uuid = ?", applicationMsgUUID).Find(&msg).Error; err != nil {
		return nil, err
	}
	fromUser := User{ID: msg.FromUser}
	db.First(&fromUser)
	_, room, err := AddFriendService(hub, &fromUser, msg.UserId)
	uuid, _ := GenerateUUID()
	hub.Messages <- Message{MsgType: SingleMessage, UUID: uuid, FromUser: msg.UserId, Content: "现在我们是朋友了，可以开始聊天了。", UserId: fromUser.ID, RoomId: room.ID}.GetDetails()
	err = checkedApplicationMessage(msg)
	return room, err
}

func RejectFriendApplication(applicationMsgUUID string) error {
	var msg Message
	if err := db.Where("uuid = ?", applicationMsgUUID).Find(&msg).Error; err != nil {
		return err
	}
	return checkedApplicationMessage(msg)
}

func SetApplicationRead(applicationMsgUUIDs []string) error {
	err := db.Model(Message{}).Where("uuid IN (?)", applicationMsgUUIDs).Updates(map[string]interface{}{"read": true}).Error
	return err
}

func checkedApplicationMessage(msg Message) error {
	return db.Model(&msg).Updates(map[string]interface{}{"checked": true, "read": true}).Error
}

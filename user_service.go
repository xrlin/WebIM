package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"log"
	"reflect"
	"strconv"
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

func AddFriendService(caller *User, friendID uint) (User, error) {
	friend := User{}
	//if caller.ID == friendID {
	//	return friend, errors.New("cannot add friend of yourself")
	//}
	for CheckFriendship(caller.ID, friendID) {
		return friend, errors.New("already has friendship")
	}
	if err := db.First(&friend, friendID).Error; err != nil {
		return friend, err
	}
	if err := db.Exec("INSERT INTO friendship(user_id, friend_id) VALUES (?, ?)", caller.ID, friend.ID).Error; err != nil {
		return friend, err
	}
	return friend, nil
}

// TODO remove friend from caller's hub after removed
func RemoveFriend(hub *Hub, user User, friendID uint) error {
	return RemoveFriendFromDB(user.ID, friendID)
}

func RemoveFriendFromDB(userID, friendID uint) error {
	if !CheckFriendship(userID, friendID) {
		return errors.New("the user is not your friend")
	}
	return db.Exec("DELETE FROM friendship WHERE user_id = ? AND friend_id = ?", userID, friendID).Error
}

// Check user2(represented with userId2) if a friend of user1(represented with userId2)
func CheckFriendship(userId1, userId2 uint) bool {
	type result struct {
		Count uint
	}
	var r result
	db.Raw("SELECT COUNT(*) FROM friendship WHERE user_id = ? AND friend_id = ?", userId1, userId2).Scan(&r)
	log.Printf("result in check friendship %#v", r)
	return r.Count > 0
}

func GetUserFriends(user User) []*User {
	friends := make([]*User, 0)

	rawSql := "SELECT * FROM users as f INNER JOIN friendship AS r ON r.friend_id = f.id AND r.user_id = ?"
	db.Raw(rawSql, user.ID).Scan(&friends)
	log.Printf("%#v", friends)
	return friends
}

func wrapToUserDetail(user User) UserDetail {
	return UserDetail{user, user.AvatarUrl()}
}

func wrapToUserDetailArray(users []User) []UserDetail {
	details := make([]UserDetail, len(users))
	for idx, user := range users {
		details[idx] = UserDetail{user, user.AvatarUrl()}
	}
	return details
}

func GetUserRooms(user User) []*Room {
	rooms := make([]*Room, 0)
	db.Model(user).Related(&rooms, "Rooms")
	return rooms
}

func SearchUsersByName(name string) []*User {
	users := make([]*User, 0)
	db.Where("name LIKE ?", "%"+name+"%").Find(&users)
	log.Println("Search users with name ", name)
	log.Println(users)
	return users
}

// An application to ask for a friendship relation with another user
func AddFriendApplication(fromUser User, toUserID uint, extraMessage string) error {
	uuid, err := GenerateUUID()
	if err != nil {
		return err
	}
	toUser := User{ID: toUserID}
	err = db.First(&toUser).Error
	if err != nil {
		return err
	}
	msg := Message{UUID: uuid, From: fromUser.ID, To: toUserID, MsgType: FriendshipMessage, Payload: extraMessage}
	log.Printf("Add friend msg %#v", msg)
	if err = SaveInboxMessage(msg); err != nil {
		return err
	}
	return nil
}

func SaveInboxMessage(message Message) error {
	var rawMsg []byte
	rawMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	inboxMsg := Inbox{Message: postgres.Jsonb{RawMessage: rawMsg}}
	return db.Create(inboxMsg).Error
}

func PassFriendApplication(user *User, applicationMsgUUID string) error {
	msg, err := retrieveMessageWithUUID(applicationMsgUUID)
	if user.ID != msg.To {
		return errors.New("you have no right to accept the application")
	}
	err = checkedApplicationMessage(msg)
	if err != nil {
		return err
	}
	_, err = AddFriendService(user, msg.To)
	if err != nil {
		return err
	}
	welcomeMessage := WelcomeMessage(user.ID, msg.From)
	DeliverMessage(*welcomeMessage, welcomeMessage.Topic, user)
	return err
}

// 添加好友后的欢迎语
func WelcomeMessage(from uint, to uint) *Message {
	topic := "IMS/P2P/" + strconv.Itoa(int(to)) + "/" + strconv.Itoa(int(from))
	message := Message{From: from, To: to, Topic: topic, MsgType: SingleMessage}
	return &message
}

func retrieveMessageWithUUID(uuid string) (Message, error) {
	inbox := Inbox{}
	var message Message
	err := db.Exec(`select message from inboxex where message::jsonb @> {"uuid": %s}::jsonb`, uuid).First(&inbox).Error
	if err != nil {
		return message, err
	}
	err = json.Unmarshal(inbox.Message.RawMessage, &message)
	return message, err
}

func RejectFriendApplication(applicationMsgUUID string) error {
	msg, err := retrieveMessageWithUUID(applicationMsgUUID)
	if err != nil {
		return err
	}
	if err := checkedApplicationMessage(msg); err != nil {
		return err
	}
	rejectMessage := &Message{From: msg.To, Payload: "好友申请被拒绝"}
	DeliverMessageToInbox(msg.From, rejectMessage)
	return nil
}

func SetApplicationRead(applicationMsgUUIDs []string) error {
	err := db.Model(Message{}).Where("uuid IN (?)", applicationMsgUUIDs).Updates(map[string]interface{}{"read": true}).Error
	return err
}

func UpdateProfileService(user *User, profile Profile) error {
	profileType := reflect.TypeOf(profile)
	log.Print(*user)
	finalUserValue := reflect.ValueOf(user).Elem()
	for i := 0; i < profileType.NumField(); i++ {
		name := profileType.Field(i).Name
		value := reflect.ValueOf(profile).FieldByName(name)
		if value.IsValid() {
			finalUserValue.FieldByName(name).Set(value)
		}
	}
	return db.Save(user).Error
}

func UpdatePasswordService(user *User, oldPassword, newPassword string) error {
	if _, ok := ValidateUser(user.Name, oldPassword); !ok {
		return errors.New("old password is incorrect")
	}
	newPasswordHsah, _ := getPasswordHash(newPassword)
	return db.Model(user).Update("passwordHash", newPasswordHsah).Error
}

func checkedApplicationMessage(msg Message) error {
	return db.Model(&msg).Updates(map[string]interface{}{"checked": true, "read": true}).Error
}

// Check all friend identified with id has a friend of user
func checkFriendship(user User, friendIds []int) bool {
	log.Printf("checkFriendship with user_id: %d, friendIds: %#v", user.ID, friendIds)
	type Result struct {
		Count int
	}
	if len(friendIds) == 0 {
		return false
	}
	rows, err := db.Raw("SELECT user_id FROM friendship WHERE friend_id = ? ", user.ID).Rows()
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		rows.Scan(&userId)
		if !checkIntInSlice(userId, friendIds) {
			return false
		}
	}
	return true
}

func checkIntInSlice(value int, slice []int) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}
	return false
}

func checkUserInRoom(userID, roomID int) bool {
	var count int
	db.Raw("SELECT COUNT(*) as count FROM user_rooms WHERE user_id = ? AND room_id = ?", userID, roomID).Row().Scan(&count)
	return count > 0
}

package services

import (
	"errors"
	"fmt"
	"github.com/xrlin/WebIM/server/database"
	"github.com/xrlin/WebIM/server/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

/* Check if user info is valid with username and password */
func ValidateUser(userName, password string) (*models.User, bool) {
	user := models.FindUserByName(userName)
	if user == nil {
		return nil, false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, false
	}
	return user, true
}

/* Register and create a new user */
func RegisterUser(user *models.User) error {
	return models.CreateUser(user)
}

func UpdateAvatar(user *models.User, avatar string) error {
	return database.DBConn.Model(user).Update("avatar", avatar).Error
}

type room struct {
	models.Room
	UniqueName string `json:"unique_name"`
}

func GetUserRecentRooms(user *models.User) []room {
	rooms := make([]models.Room, 0)
	database.DBConn.Preload("Users").Model(user).Related(&rooms, "Rooms")
	result := make([]room, 0)
	for _, r := range rooms {
		uniqueName := fmt.Sprintf("%s_%d", r.Name, r.ID)
		fmt.Println(r.Users)
		if r.RoomType == models.FriendRoom {
			// Set room name with the friend's name
			var friend models.User
			for _, u := range r.Users {
				if u.ID != user.ID {
					friend = u
				}
			}
			r.Users = []models.User{friend}
			r.Name = friend.Name
		}
		result = append(result, room{Room: r, UniqueName: uniqueName})
	}
	fmt.Printf("%#v", result)
	return result
}

func AddFriend(hub *Hub, caller *models.User, friendID uint) (*models.User, *models.Room, error) {
	if caller.ID == friendID {
		return nil, nil, errors.New("Cannot add friend of yourself")
	}
	friends := GetUserFriends(*caller)
	for _, user := range friends {
		if user.ID == friendID {
			return nil, nil, errors.New("Already has friendship")
		}
	}
	friend := &models.User{}
	if err := database.DBConn.First(friend, friendID).Error; err != nil {
		return nil, nil, err
	}
	room := models.Room{RoomType: models.FriendRoom, Users: []models.User{*caller, *friend}}
	if err := database.DBConn.Create(&room).Error; err != nil {
		return nil, nil, err
	}
	room.Name = friend.Name
	hub.UpdateRoom <- &room
	return friend, &room, nil
}

func GetUserFriends(user models.User) []*models.User {
	friends := make([]*models.User, 0)

	rawSql := "SELECT * FROM users INNER JOIN user_rooms ON user_rooms.user_id = users.id WHERE users.id!=?  AND users.deleted_at IS NULL AND (user_rooms.room_id IN (SELECT id FROM rooms INNER JOIN user_rooms ON user_rooms.room_id = rooms.id AND user_rooms.user_id=? WHERE rooms.room_type = ? AND rooms.deleted_at IS NULL));"
	database.DBConn.Raw(rawSql, user.ID, user.ID, models.FriendRoom).Scan(&friends)
	log.Printf("%#v", friends)
	return friends
}

func SearchUsersByName(name string) []*models.User {
	users := make([]*models.User, 0)
	database.DBConn.Where("name LIKE ?", "%"+name+"%").Find(&users)
	log.Println("Search users with name ", name)
	log.Println(users)
	return users
}

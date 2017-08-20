package services

import (
	"bytes"
	"errors"
	"github.com/xrlin/WebIM/server/database"
	"github.com/xrlin/WebIM/server/models"
)

func CreateRoom(hub *Hub, creator *models.User, userIds []int) (*models.Room, error) {
	users := make([]models.User, 0)
	if err := database.DBConn.Where("id in (?)", userIds).Find(&users).Error; err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("Must provide at least one another user before create room")
	}
	users = append(users, *creator)
	strBuf := bytes.Buffer{}
	for _, u := range users {
		strBuf.WriteString(u.Name)
		strBuf.WriteString(" ")
	}
	room := models.Room{Name: strBuf.String(), Users: users}
	err := database.DBConn.Create(&room).Error
	hub.UpdateRoom <- &room
	return &room, err
}

// User leave the chat room forever and will also leave the room in hub
func LeaveRoom(hub *Hub, user *models.User, room_id int) (*models.Room, error) {
	var room models.Room
	if err := database.DBConn.Where("id = ?", room_id).Find(&room).Error; err != nil {
		return nil, err
	}
	params := leaveRoomParam{user: user, room: &room}
	hub.LeaveRoom <- &params
	if err := deleteUserFromRoom(user, &room); err != nil {
		return nil, err
	}
	return &room, nil
}

func deleteUserFromRoom(user *models.User, room *models.Room) error {
	err := database.DBConn.Model(&room).Association("users").Delete(*user).Error
	//if err !=nil {
	//	return err
	//}
	//if room.RoomType == models.FriendRoom {
	//	return database.DBConn.Delete(room).Error
	//}
	return err
}

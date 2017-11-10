package main

import (
	"bytes"
	"errors"
)

func CreateRoomService(hub *Hub, creator *User, userIds []int) (*Room, error) {
	users := make([]User, 0)
	if err := db.Where("id in (?)", userIds).Find(&users).Error; err != nil {
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
	room := Room{Name: strBuf.String(), Users: users}
	err := db.Create(&room).Error
	hub.UpdateRoom <- &room
	return &room, err
}

// User leave the chat room forever and will also leave the room in hub
func LeaveRoomService(hub *Hub, user *User, room_id int) (*Room, error) {
	var room Room
	if err := db.Where("id = ?", room_id).Find(&room).Error; err != nil {
		return nil, err
	}
	params := leaveRoomParam{user: user, room: &room}
	hub.LeaveRoom <- &params
	if err := deleteUserFromRoom(user, &room); err != nil {
		return nil, err
	}
	return &room, nil
}

func deleteUserFromRoom(user *User, room *Room) error {
	err := db.Model(&room).Association("users").Delete(*user).Error
	//if err !=nil {
	//	return err
	//}
	//if room.RoomType == FriendRoom {
	//	return database.db.Delete(room).Error
	//}
	return err
}

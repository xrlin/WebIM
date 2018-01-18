package main

import (
	"errors"
	"fmt"
	"log"
)

type CreateRoomParam struct {
	Name   string `json:"name" binding:"required"`
	Avatar string `json:"avatar"`
}

func CreateRoomService(hub *Hub, creator *User, params CreateRoomParam) (Room, error) {
	room := Room{Name: params.Name, Users: []User{*creator}, Avatar: params.Avatar}
	err := db.Create(&room).Error
	hub.UpdateRoom <- &room
	return room, err
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

func JoinRoomService(user *User, roomId int) (room Room, err error) {
	room, err = findRoomById(roomId)
	if err != nil {
		return
	}
	sql := "INSERT INTO user_rooms(user_id, room_id) VALUES (?, ?)"
	err = db.Exec(sql, user.ID, roomId).Error
	return
}

func GetRoomUsers(room Room) (users []User, err error) {
	err = db.Model(room).Related(&room.Users, "Users").Error
	if err != nil {
		users = make([]User, 0)
		return
	}
	users = room.Users
	return
}

func SearchRoomService(roomName string) (rooms []Room, err error) {
	if roomName == "" {
		err = errors.New("name cannot be empty")
		return
	}
	sql := "SELECT id, name FROM rooms WHERE name LIKE ?"
	pattern := "%" + roomName + "%"
	err = db.Raw(sql, pattern).Scan(&rooms).Error
	return
}

func InviteUserToRoomService(roomID int, userIds []int, caller User) (err error) {
	if len(userIds) == 0 {
		err = errors.New("users to invite is blank")
		return
	}
	_, err = findRoomById(roomID)
	if err != nil {
		return
	}
	if !checkUserInRoom(int(caller.ID), roomID) {
		err = errors.New("only the members of room can invite others")
		return
	}
	if !checkFriendship(caller, userIds) {
		err = errors.New("not allow to invite some friends")
		return
	}
	sql := "INSERT INTO user_rooms(user_id, room_id) VALUES "
	values := make([]interface{}, 0)
	// prepare the sql
	for pos, id := range userIds {
		if id <= 0 {
			err = fmt.Errorf("unpermitted user id with: %d", id)
			return
		}
		sql += "(?, ?)"
		if pos < len(userIds)-1 {
			sql += ","
		}
		values = append(values, id, roomID)
	}
	log.Printf("Invite friends to group %d with sql process: %s", roomID, sql)
	log.Printf("values: %#v", values)
	err = db.Exec(sql, values...).Error
	return
}

func findRoomById(roomId int) (room Room, err error) {
	err = db.Where("id = ?", roomId).Find(&room).Error
	return
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

func GetRoomInfoService(roomID int) (Room, error) {
	room, err := findRoomById(roomID)
	return room, err
}

func UpdateRoomService(roomID int, caller User, param updateRoomParam) (room Room, err error) {
	room, err = findRoomById(roomID)
	if err != nil {
		return
	}
	if !checkUserInRoom(int(caller.ID), roomID) {
		err = fmt.Errorf("user with id %d not in room %d", caller.ID, roomID)
		return
	}
	err = db.Model(&room).Updates(param).Error
	return
}

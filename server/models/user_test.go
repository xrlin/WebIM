package models

import (
	"fmt"
	"github.com/xrlin/WebIM/server/database"
	"reflect"
	"testing"
	"time"
)

func TestFindUserByName(t *testing.T) {
	u := User{Name: fmt.Sprintf("User%d", time.Now().Unix()), PasswordHash: "just for test"}
	if FindUserByName(u.Name) != nil {
		t.Error(fmt.Sprintf("Usre %s found before created!", u.Name))
	}
	database.DBConn.Create(&u)
	if FindUserByName(u.Name) == nil {
		t.Error(fmt.Sprintf("Usre %s not found after run created!", u.Name))
	} else {
		t.Log("TestFindUserByName successfully!")
	}
}

func TestCreateUser(t *testing.T) {
	u := &User{Name: fmt.Sprintf("User%d", time.Now().Unix()), Password: "test"}
	if err := CreateUser(u); err == nil && !database.DBConn.NewRecord(u) {
		t.Log(fmt.Sprintf("CreteUser with name %s and password %s successfully", u.Name, u.Password))
	} else {
		t.Error(fmt.Sprintf("CreteUser with name %s and password %s failed", u.Name, u.Password))
		t.Error(err)
	}
}

func TestUser_RoomNames(t *testing.T) {
	// Init data
	rooms := []Room{{Name: "test"}}
	user := User{Name: "Hello", PasswordHash: "ssdfsdfsdfsdf", Rooms: rooms}
	database.DBConn.Create(&user)
	defer func() {
		database.DBConn.Where("id = ?", user.ID).Delete(&User{})
		database.DBConn.Model(&user).Association("Rooms").Clear()
		for _, room := range rooms {
			database.DBConn.Where("id = ?", room.ID).Delete(&room)
		}
	}()

	user.Rooms = []Room{}
	roomNames := user.RoomNames()
	expected := make([]string, len(rooms))
	for idx, room := range rooms {
		expected[idx] = fmt.Sprintf("room_%v", room.ID)
	}
	if !reflect.DeepEqual(roomNames, expected) {
		t.Fail()
	}
}

func TestUser_UserRoomName(t *testing.T) {
	u := User{ID: 3}
	if u.UserRoomName() != "room_user_3" {
		t.Fail()
	}
}

package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestValidateUser(t *testing.T) {
	u := &User{Name: fmt.Sprintf("User%d", time.Now().Unix()), Password: "test"}
	CreateUser(u)
	if _, ok := ValidateUser(u.Name, u.Password); ok {
		t.Log(fmt.Sprintf("Valivate User with name %s and password %s successfully", u.Name, u.Password))
	} else {
		t.Error(fmt.Sprintf("Valivate User with name %s and password %s failed", u.Name, u.Password))
	}
	if _, ok := ValidateUser(u.Name, "invalid password"); ok {
		t.Error(fmt.Sprintf("Valivate User with name %s and password %s should failed", u.Name, "invalid password"))
	} else {
		t.Log("Success")
	}
}

func TestRegisterUser(t *testing.T) {
	u := User{Name: fmt.Sprintf("User%d", time.Now().Unix()), Password: "test"}
	if err := RegisterUser(&u); err == nil {
		t.Log(fmt.Sprintf("Register with user %+v successfully", u))
	} else {
		t.Error(fmt.Sprintf("Register with user %+v failed", u))
		t.Error(err)
	}
}

func getRandomID() uint {
	return uint(rand.Int31())
}

func TestCheckFriendship(t *testing.T) {
	userId1 := getRandomID()
	userId2 := getRandomID()
	if CheckFriendship(uint(userId1), uint(userId2)) {
		t.Fail()
	}
	defer func() {
		db.Exec("DELETE FROM friendship WHERE user_id = ? AND friend_id = ?", userId1, userId2)
	}()
	db.Exec("INSERT INTO friendship(user_id, friend_id) VALUES(?, ?)", userId1, userId2)
	if !CheckFriendship(userId1, userId2) {
		t.Fail()
	}
}

func TestRemoveFriendFromDB(t *testing.T) {
	userId1 := getRandomID()
	userId2 := getRandomID()
	if err := RemoveFriendFromDB(userId1, userId2); err == nil {
		t.Fail()
	}
	defer func() {
		db.Exec("DELETE FROM friendship WHERE user_id = ? AND friend_id = ?", userId1, userId2)
	}()
	db.Exec("INSERT INTO friendship(user_id, friend_id) VALUES(?, ?)", userId1, userId2)
	if err := RemoveFriendFromDB(userId1, userId2); err != nil {
		t.Fail()
	}

}

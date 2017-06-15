package models

import (
	"testing"
	"time"
	"fmt"
	"github.com/xrlin/WebIM/server/database"
)

func TestFindUserByName(t *testing.T) {
	u := User{Name: fmt.Sprintf("User%d", time.Now().Unix()), PasswordHash: "just for test"}
	if FindUserByName(u.Name) != nil {
		t.Error(fmt.Sprintf("Usre %s found before created!", u.Name))
	}
	database.DBConnection.Create(&u)
	if FindUserByName(u.Name) == nil {
		t.Error(fmt.Sprintf("Usre %s not found after run created!", u.Name))
	} else {
		t.Log("TestFindUserByName successfully!")
	}
}

func TestCreateUser(t *testing.T) {
	u := &User{Name: fmt.Sprintf("User%d", time.Now().Unix()), Password: "test"}
	if err := CreateUser(u); err == nil && !database.DBConnection.NewRecord(u) {
		t.Log(fmt.Sprintf("CreteUser with name %s and password %s successfully", u.Name, u.Password))
	} else {
		t.Error(fmt.Sprintf("CreteUser with name %s and password %s failed", u.Name, u.Password))
		t.Error(err)
	}
}

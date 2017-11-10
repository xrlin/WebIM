package main

import (
	"fmt"
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

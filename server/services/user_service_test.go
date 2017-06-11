package services

import (
	"testing"
	"github.com/xrlin/WebIM/server/models"
	"fmt"
	"time"
)

func TestValidateUser(t *testing.T) {
	u := &models.User{Name: fmt.Sprintf("User%d", time.Now().Unix()), Password: "test"}
	models.CreateUser(u)
	if ValidateUser(u.Name, u.Password) {
		t.Log(fmt.Sprintf("Valivate User with name %s and password %s successfully", u.Name, u.Password))
	} else {
		t.Error(fmt.Sprintf("Valivate User with name %s and password %s failed", u.Name, u.Password))
	}
	if ValidateUser(u.Name, "invalid password") {
		t.Error(fmt.Sprintf("Valivate User with name %s and password %s should failed", u.Name, "invalid password"))
	} else {
		t.Log("Success")
	}
}

func TestRegisterUser(t *testing.T) {
	u := models.User{Name: fmt.Sprintf("User%d", time.Now().Unix()), Password: "test"}
	if err := RegisterUser(&u); err == nil {
		t.Log(fmt.Sprintf("Register with user %+v successfully", u))
	} else {
		t.Error(fmt.Sprintf("Register with user %+v failed", u))
		t.Error(err)
	}
}

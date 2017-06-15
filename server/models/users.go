package models

import (
	"github.com/jinzhu/gorm"
	"github.com/xrlin/WebIM/server/database"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"fmt"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null;unique;primary_key"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

func FindUserByName(name string) *User {
	var user User
	database.DBConnection.Where("name = ?", name).First(&user)
	// No record found
	if user.ID == 0 {
		return nil
	}
	return &user
}

/* Create user from a pointer to user, if successed update the pointer to the User struct with id and other fields form database*/
func CreateUser(user *User) error {
	if !database.DBConnection.NewRecord(user) {
		return errors.New(fmt.Sprintf("User with name %s has exist!", user.Name))
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)
	return database.DBConnection.Create(user).Error
}

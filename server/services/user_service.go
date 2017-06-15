package services

import (
	"github.com/xrlin/WebIM/server/models"
	"golang.org/x/crypto/bcrypt"
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

package controllers

import (
	"errors"
	"strings"
	"net/http"
	"github.com/xrlin/WebIM/server/services"
	"time"
	"fmt"
	"github.com/xrlin/WebIM/server/models"
	"github.com/gin-gonic/gin"
)

type Login struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	Login
}

func UserToken(c *gin.Context) {
	var login Login
	if err := c.BindJSON(&login); err != nil {
		fmt.Println(login)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	if user, validated := services.ValidateUser(login.UserName, login.Password); validated {
		// TODO config SignedKey
		tokenSetvice := services.TokenService{time.Hour * 3, "test"}
		token, err := tokenSetvice.Generate(int(user.ID), user.Name)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"token": token})
		} else {
			c.JSON(http.StatusOK, gin.H{"errors": err.Error()})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Username or password is invalid!"})
	}
}

func CreateUser(c *gin.Context) {
	var registerInfo Register
	if err := c.BindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	user := &models.User{Name: registerInfo.UserName, Password: registerInfo.Password}
	if err := services.RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s created successfully.", user.Name)})
	}
}

// Get token from Authorization header
// Token in header is in format:
//		Authorization: Bearer yJhbGciOiJIUzI1NiIsInR5...JIUz
func getTokenFromContext(c *gin.Context) (string, error) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return token, errors.New("Token is not found in header.")
	}
	return strings.TrimPrefix(token, "Bearer "), nil
}

// Check the context if have all the requiredParams and return them
// return presentParams, absentParams and checked result
func checkRequiredParams(c *gin.Context, requiredParams []string) ([]string, []string, bool) {
	params := c.Params
	results := []string{}
	absentParams := []string{}
	for _, v := range requiredParams {
		if value, ok := params.Get(v); ok {
			results = append(results, value)
		} else {
			absentParams = append(absentParams, v)
		}
	}
	return results, absentParams, len(absentParams) == 0
}

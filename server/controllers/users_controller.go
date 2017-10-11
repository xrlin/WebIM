package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/models"
	"github.com/xrlin/WebIM/server/services"
	"net/http"
	"strings"
)

type Login struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	Login
}

func UserToken(c *gin.Context) {
	if userObj, ok := c.Get("user"); ok {
		user := userObj.(*models.User)
		// TODO config SignedKey
		tokenService := services.GetTokenService()
		token, err := tokenService.Generate(int(user.ID), user.Name)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"token": token})
		} else {
			c.JSON(http.StatusOK, gin.H{"errors": err.Error()})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "User is nil"})
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

func GetUserInfo(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	user := userObj.(*models.User)
	c.JSON(http.StatusOK, user)
}

func GetRecentRooms(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	user := userObj.(*models.User)
	rooms := services.GetUserRecentRooms(user)
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func SearchUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Must provide name."})
		return
	}
	users := services.SearchUsersByName(name)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

type friendInfo struct {
	FriendID uint `json:"friend_id" binding: "required"`
}

func AddFriend(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	user := userObj.(*models.User)
	friend := friendInfo{}
	if err := c.BindJSON(&friend); err != nil || friend.FriendID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Must provide the user id of friend."})
		return
	}
	_, room, err := services.AddFriend(hub, user, friend.FriendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Add friend failed." + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"room": *room})
}

func GetFriends(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	user := userObj.(*models.User)
	friends := services.GetUserFriends(*user)
	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

type addRoomParams struct {
	UserIds []int `json:"user_ids" binding:"required"`
}

func CreateRoom(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	var params addRoomParams
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	user := userObj.(*models.User)
	if room, err := services.CreateRoom(hub, user, params.UserIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"room": *room})
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

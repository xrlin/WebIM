package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
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
		user := userObj.(*User)
		// TODO config SignedKey
		tokenService := GetTokenService()
		token, err := tokenService.Generate(user.ID, user.Name)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"token": token})
		} else {
			c.JSON(http.StatusOK, gin.H{"errors": err.Error()})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "User is nil"})
	}
}

func CreateUserHandler(c *gin.Context) {
	var registerInfo Register
	if err := c.BindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	user := &User{Name: registerInfo.UserName, Password: registerInfo.Password}
	if err := RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s created successfully.", user.Name)})
	}
}

type avatarInfo struct {
	Avatar string `json:"avatar" binding:"required"`
}

type UserDetail struct {
	User
	AvatarUrl string `json:"avatarURL"`
}

func UpdateAvatar(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	var avatarInfo avatarInfo
	if err := c.BindJSON(&avatarInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}
	user := userObj.(*User)
	if err := UpdateAvatarService(user, avatarInfo.Avatar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
	}
	userDetail := UserDetail{User: *user, AvatarUrl: user.AvatarUrl()}
	c.JSON(http.StatusOK, gin.H{"user": userDetail})
}

func GetUserInfo(c *gin.Context) {
	userID := c.Query("id")
	var user User
	var err error
	if userID != "" {
		err = db.Where("id = ?", userID).First(&user).Error
	} else {
		userObj, _ := c.Get("user")
		user = *userObj.(*User)
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No such user!"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserRoomsController(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	user := userObj.(*User)
	rooms := GetUserRecentRooms(user)
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func SearchUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Must provide name."})
		return
	}
	users := SearchUsersByName(name)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

type friendInfo struct {
	FriendID uint `json:"friendID" binding:"required"`
}

func AddFriend(c *gin.Context) {
	user, _ := c.Get("user")
	param := friendInfo{}
	c.BindJSON(&param)
	friend, err := AddFriendService(hub, user.(*User), param.FriendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Add friend failed." + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"friend": friend})
}

func DeleteFriend(c *gin.Context) {
	user, _ := c.Get("user")
	friendID, _ := strconv.ParseUint(c.Param("friendID"), 10, 64)
	if err := RemoveFriend(hub, *user.(*User), uint(friendID)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{})
}

func GetFriends(c *gin.Context) {
	user, _ := c.Get("user")
	friends := GetUserFriends(*user.(*User))
	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

func GetRooms(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	user := userObj.(*User)
	rooms := GetUserRooms(*user)
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func FriendApplication(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	var params struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := userObj.(*User)
	err := AddFriendApplication(hub, *user, params.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully apply a friendship application."})
}

func CheckFriendApplication(c *gin.Context) {
	_, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	var reqParams struct {
		Action string `json:"action" binding:"required"`
		UUID   string `json:"uuid" binding:"required"`
	}
	if err := c.BindJSON(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if reqParams.Action == "pass" {
		err := PassFriendApplication(hub, reqParams.UUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatus(http.StatusCreated)
		return
	} else {
		err := RejectFriendApplication(reqParams.UUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"messsage": "Reject friend application successfully"})
	}
}

func AckReadFriendApplications(c *gin.Context) {
	_, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No such user!"})
		return
	}
	var messageParams struct {
		UUIDArray []string `binding:"required" json:"uuid_array"`
	}
	c.BindJSON(&messageParams)
	err := SetApplicationRead(messageParams.UUIDArray)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

type Profile struct {
	Name string
}

func UpdateProfile(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Authorize failed!"})
		return
	}
	user := userObj.(*User)
	var profile Profile
	if c.ShouldBind(&profile) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Params invalid"})
		return
	}
	if UpdateProfileService(user, profile) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messsage": "Update Successfully"})
}

type passwordParams struct {
	OldPassword string `binding:"required" json:"oldPassword"`
	NewPassword string `binding:"required" json:"newPassword"`
}

func UpdatePassword(c *gin.Context) {
	userObj, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Authorize failed!"})
		return
	}
	user := userObj.(*User)
	var params passwordParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := UpdatePasswordService(user, params.OldPassword, params.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update password success."})
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

func getPasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

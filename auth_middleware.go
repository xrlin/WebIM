package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := getToken(c); token != "" {
			log.Print("Will go to authByToken")
			authByToken(token, c)
		} else {
			log.Print("Will go to authByUsername")
			authByUsernameAndPassword(c)
		}
	}
}

func authByUsernameAndPassword(c *gin.Context) {
	var login Login
	if err := c.BindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.Abort()
		return
	}
	if user, validated := ValidateUser(login.UserName, login.Password); validated != true {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized!"})
		c.Abort()
		return
	} else {
		c.Set("user", user)
		c.Next()
	}
}

func authByToken(token string, c *gin.Context) {
	fmt.Print("Go into authByToken")
	tokenService := GetTokenService()
	tokenInfo, err := tokenService.Parse(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
		c.Abort()
		return
	}
	user := FindUserById(tokenInfo.UserId)
	if tokenInfo.ExpiresAt-time.Now().Unix() < 900 {
		// Add new token in header within 30 minutes in response header
		tokenService.Duration = time.Minute * 30
		if newToken, err := tokenService.Generate(tokenInfo.UserId, tokenInfo.UserName); err == nil {
			c.Writer.Header().Set("Token", newToken)
		}
	}
	c.Set("user", user)
	c.Next()
}

func getToken(c *gin.Context) string {
	token := ""
	if token = c.Request.Header.Get("Authorization"); token != "" {
		return strings.TrimPrefix(token, "Bearer ")
	}
	if token = c.Request.URL.Query().Get("token"); token != "" {
		return token
	}
	return ""
}

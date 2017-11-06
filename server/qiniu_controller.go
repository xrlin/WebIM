package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response the upload key to client
func UploadToken(c *gin.Context) {
	_, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	upToken := GenerateUploadToken()
	c.JSON(http.StatusOK, gin.H{"uptoken": upToken})
}

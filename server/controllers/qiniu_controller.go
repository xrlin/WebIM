package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/services"
	"net/http"
)

// Response the upload key to client
func UploadToken(c *gin.Context) {
	_, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	upToken := services.GenerateUploadToken()
	c.JSON(http.StatusOK, gin.H{"uptoken": upToken})
}

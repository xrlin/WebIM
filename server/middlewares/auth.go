package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/controllers"
	"github.com/xrlin/WebIM/server/services"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login controllers.Login
		if err := c.BindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			c.Abort()
			return
		}
		if user, validated := services.ValidateUser(login.UserName, login.Password); validated != true {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized!"})
			c.Abort()
			return
		} else {
			c.Set("user", user)
			c.Next()
		}
	}
}

package main

import (
	"github.com/xrlin/WebIM/server/database"
	"github.com/xrlin/WebIM/server/models"
)

func main() {
	database.DBConnection.AutoMigrate(&models.User{})
}

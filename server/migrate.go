package main

import (
	"github.com/xrlin/WebIM/server/database"
	"github.com/xrlin/WebIM/server/models"
)

func main() {
	database.DBConn.AutoMigrate(&models.User{}, &models.Message{}, &models.Room{})
	type user_rooms struct{}
	database.DBConn.Model(&user_rooms{}).AddUniqueIndex("idx_user_room_id", "user_id", "room_id")
}

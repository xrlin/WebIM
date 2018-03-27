package main

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Inbox struct {
	gorm.Model
	Message postgres.Jsonb `json:"message"`
}

package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xrlin/WebIM/server/config"
)

var (
	DBConnection *gorm.DB
	DBErr        error
)

func init() {
	DBConnection, DBErr = gorm.Open(config.DatabaseCfg.Type, config.DatabaseCfg.DBInfoString())
	if DBErr != nil {
		panic(DBErr)
	}
	DBConnection.DB().SetMaxIdleConns(50)
}

package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xrlin/WebIM/server/config"
)

var (
	DBConn *gorm.DB
	DBErr  error
)

func init() {
	DBConn, DBErr = gorm.Open(config.DatabaseCfg.Type, config.DatabaseCfg.DBInfoString())
	if DBErr != nil {
		panic(DBErr)
	}
	DBConn.DB().SetMaxIdleConns(50)
}

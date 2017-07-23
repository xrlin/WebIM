package database

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xrlin/WebIM/server/config"
	"time"
)

var (
	DBConn    *gorm.DB
	DBErr     error
	RedisPool *redis.Pool
)

func newRedisPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func init() {
	DBConn, DBErr = gorm.Open(config.DatabaseCfg.Type, config.DatabaseCfg.DBInfoString())
	if DBErr != nil {
		panic(DBErr)
	}
	DBConn.DB().SetMaxIdleConns(50)

	RedisPool = newRedisPool("localhost:6379")
}

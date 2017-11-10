package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

var (
	db        *gorm.DB
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
	db, DBErr = gorm.Open(DatabaseCfg.Type, DatabaseCfg.DBInfoString())
	if DBErr != nil {
		panic(DBErr)
	}
	db.DB().SetMaxIdleConns(50)
	db.LogMode(true)
	db.BlockGlobalUpdate(true)

	RedisPool = newRedisPool("localhost:6379")
}

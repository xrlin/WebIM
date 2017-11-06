package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Database struct {
	Type     string
	Host     string
	User     string
	DBName   string
	SSLMode  string
	Password string
}

type Qiniu struct {
	SecretKey, AccessKey, Bucket, FileDomain string
}

func (cfg *Database) DBInfoString() string {
	return fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.User, cfg.DBName, cfg.SSLMode, cfg.Password,
	)
}

var (
	DatabaseCfg *Database
	QiniuCfg    *Qiniu
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	absPath := filepath.Join(filepath.Dir(file), "./database_config.json")
	f, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(f)
	DatabaseCfg = new(Database)
	if err := decoder.Decode(DatabaseCfg); err != nil {
		panic(err)
	}

	QiniuCfg = new(Qiniu)
	qiniuCfgPath := filepath.Join(filepath.Dir(file), "./qiniu_config.json")
	f, err = os.Open(qiniuCfgPath)
	if err != nil {
		panic(err)
	}
	decoder = json.NewDecoder(f)
	if err := decoder.Decode(QiniuCfg); err != nil {
		panic(err)
	}
}

package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"os"
	"encoding/json"
)

type Database struct {
	Type     string
	Host     string
	User     string
	DBName   string
	SSLMode  string
	Password string
}

func (cfg *Database) DBInfoString() string {
	return fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.User, cfg.DBName, cfg.SSLMode, cfg.Password,
	)
}

var (
	DatabaseCfg *Database
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	absPath := filepath.Join(filepath.Dir(file), "./database.json")
	f, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(f)
	DatabaseCfg = new(Database)
	if err := decoder.Decode(DatabaseCfg); err != nil {
		panic(err)
	}
}

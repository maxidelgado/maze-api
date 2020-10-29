package config

import (
	"fmt"
	"os"
)

func init() {
	dbUser := getEnv("DB_USER", "root")
	dbPwd := getEnv("DB_PWD", "example")
	dbHost := getEnv("DB_HOST", "localhost:27017")

	Router = RouterCfg{
		Host:     getEnv("ROUTER_HOST", ":3000"),
		BasePath: getEnv("BASE_PATH", "/api/v1"),
	}

	DB = MongoDB{
		Uri:        fmt.Sprintf(mgoUriPattern, dbUser, dbPwd, dbHost),
		Database:   getEnv("DB_NAME", "maze"),
		Collection: getEnv("DB_COL", "mazes"),
	}
}

const (
	mgoUriPattern = "mongodb://%s:%s@%s"
)

var (
	DB     MongoDB
	Router RouterCfg
)

type MongoDB struct {
	Uri        string
	Database   string
	Collection string
}

type RouterCfg struct {
	Host     string
	BasePath string
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	return v
}

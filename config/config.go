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
		Uri:            fmt.Sprintf(mgoUriPattern, dbUser, dbPwd, dbHost),
		Database:       getEnv("DB_NAME", "maze"),
		MazeCollection: getEnv("DB_MAZE_COL", "mazes"),
		GameCollection: getEnv("DB_GAME_COL", "games"),
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
	Uri            string
	Database       string
	MazeCollection string
	GameCollection string
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

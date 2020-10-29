package config

var (
	DB = MongoDB{
		Uri:        "mongodb://root:example@localhost:27017",
		Database:   "maze",
		Collection: "mazes",
	}

	Host = ":3000"

	BasePath = "/api/v1"
)

type MongoDB struct {
	Uri        string
	Database   string
	Collection string
}

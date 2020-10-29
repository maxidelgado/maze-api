package maze

import "context"

type Service interface {
	GetMaze(context.Context, string) (Maze, error)
	CreateMaze(context.Context, Coordinate) (string, error)
	UpdateQuadrants(context.Context, string, Coordinate) error
	PutSpots(context.Context, string, []Spot) error
	PutPaths(context.Context, string, []Path) error
}

type DataBase interface {
	Get(context.Context, string) (Maze, error)
	Put(context.Context, Maze) error
	Update(context.Context, Maze) error
}

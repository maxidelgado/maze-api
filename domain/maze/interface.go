package maze

import (
	"context"
)

type Service interface {
	Get(context.Context, string) (Maze, error)
	Create(context.Context, Coordinates, []Spot, []Path) (string, error)
	Update(context.Context, string, Coordinates, []Spot, []Path) error
	Delete(context.Context, string) error

	DeleteSpot(context.Context, string, Coordinates) error
	DeletePath(context.Context, string, Path) error
}

type DataBase interface {
	GetMaze(context.Context, string) (Maze, error)
	PutMaze(context.Context, Maze) error
	UpdateMaze(context.Context, Maze) error
	DeleteMaze(context.Context, string) error
}

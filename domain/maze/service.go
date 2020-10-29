package maze

import (
	"context"

	"github.com/google/uuid"
)

func New(db DataBase) Service {
	return service{db: db}
}

type service struct {
	db DataBase
}

func (s service) GetMaze(ctx context.Context, mazeId string) (Maze, error) {
	return s.db.Get(ctx, mazeId)
}

func (s service) CreateMaze(ctx context.Context, center Coordinate) (string, error) {
	maze := Maze{
		Id: uuid.New().String(),
	}

	maze.setQuadrants(center.X(), center.Y())

	if err := s.db.Put(ctx, maze); err != nil {
		return "", err
	}

	return maze.Id, nil
}

func (s service) UpdateQuadrants(ctx context.Context, mazeId string, coordinate Coordinate) error {
	m, err := s.GetMaze(ctx, mazeId)
	if err != nil {
		return err
	}

	newMaze := m.moveAxes(coordinate.X(), coordinate.Y())

	return s.db.Update(ctx, newMaze)
}

func (s service) PutSpots(ctx context.Context, mazeId string, spots []Spot) error {
	maze, err := s.GetMaze(ctx, mazeId)
	if err != nil {
		return err
	}

	for _, spot := range spots {
		maze.addSpot(spot)
	}

	return s.db.Update(ctx, maze)
}

func (s service) PutPaths(ctx context.Context, mazeId string, paths []Path) error {
	maze, err := s.GetMaze(ctx, mazeId)
	if err != nil {
		return err
	}

	for _, path := range paths {
		if maze.findSpot(path.Source) && maze.findSpot(path.Target) {
			maze.addPath(path)
		}
	}

	return s.db.Update(ctx, maze)
}

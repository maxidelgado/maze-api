package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/maxidelgado/maze-api/domain/maze"
)

func NewMaze(db maze.DataBase) maze.Service {
	return mazeSvc{db: db}
}

type mazeSvc struct {
	db maze.DataBase
}

func (s mazeSvc) Create(ctx context.Context, name string, center maze.Coordinates, spots []maze.Spot, paths []maze.Path) (string, error) {
	if name == "" {
		return "", errors.New("name is required")
	}

	m := maze.Maze{
		Id:    uuid.New().String(),
		Name:  name,
		Paths: map[string]map[string]float64{},
	}

	// Coordinates is a wrapper of [2]int64, it will create a default quadrant with center in (0, 0) if center is not specified
	m.SetQuadrants(center.X(), center.Y())

	// Add spots taking care of the corresponding quadrant for every spot
	for _, spot := range spots {
		if err := m.AddSpot(spot); err != nil {
			return "", err
		}
	}

	// Add paths to the maze taking care about source/target spots exist (fail if try to create orphan path)
	for _, path := range paths {
		if ok := m.AddPath(path.Origin, path.Destiny); !ok {
			return "", errors.New("could not add path, spot not found")
		}
	}

	// Save maze to database
	if err := s.db.PutMaze(ctx, m); err != nil {
		return "", err
	}

	return m.Id, nil
}

func (s mazeSvc) Update(ctx context.Context, mazeId string, center maze.Coordinates, spots []maze.Spot, paths []maze.Path) error {
	m, err := s.Get(ctx, mazeId)
	if err != nil {
		return err
	}

	// check if center should be moved
	if x, y := m.GetCenter(); center.X() != x || center.Y() != y {
		m = m.MoveAxes(center.X(), center.Y())
	}

	// check if should add new spots
	if len(spots) != 0 {
		for _, spot := range spots {
			if err := m.AddSpot(spot); err != nil {
				return err
			}
		}
	}

	// check if should add new paths
	// fail if try to create an orphan path
	if len(paths) != 0 {
		for _, path := range paths {
			if ok := m.AddPath(path.Origin, path.Destiny); !ok {
				return errors.New("could not add path, spot not found")
			}
		}
	}

	return s.db.UpdateMaze(ctx, m)
}

func (s mazeSvc) Delete(ctx context.Context, mazeId string) error {
	return s.db.DeleteMaze(ctx, mazeId)
}

func (s mazeSvc) DeleteSpot(ctx context.Context, mazeId string, coordinate maze.Coordinates) error {
	m, err := s.Get(ctx, mazeId)
	if err != nil {
		return err
	}

	if _, ok := m.FindSpot(coordinate.Key()); !ok {
		return errors.New("not found")
	}

	// deletes the spot and all the related paths, so it will not allow orphan paths
	m.DeleteSpot(coordinate)

	return s.db.UpdateMaze(ctx, m)
}

func (s mazeSvc) DeletePath(ctx context.Context, mazeId string, path maze.Path) error {
	m, err := s.Get(ctx, mazeId)
	if err != nil {
		return err
	}

	// deletes the path and the corresponding reverse path
	m.Paths.DeletePath(path.Origin, path.Destiny)

	return s.db.UpdateMaze(ctx, m)
}

func (s mazeSvc) Query(ctx context.Context, name string) ([]maze.Maze, error) {
	return s.db.QueryMaze(ctx, name)
}

func (s mazeSvc) Get(ctx context.Context, mazeId string) (maze.Maze, error) {
	return s.db.GetMaze(ctx, mazeId)
}

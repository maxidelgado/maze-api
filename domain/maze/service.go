package maze

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func NewService(db DataBase) Service {
	return service{db: db}
}

type service struct {
	db DataBase
}

func (s service) Get(ctx context.Context, mazeId string) (Maze, error) {
	return s.db.GetMaze(ctx, mazeId)
}

func (s service) Create(ctx context.Context, center Coordinates, spots []Spot, paths []Path) (string, error) {
	maze := Maze{
		Id:    uuid.New().String(),
		Paths: map[string]map[string]float64{},
	}

	// Coordinates is a wrapper of [2]int64, it will create a default quadrant with center in (0, 0) if center is not specified
	maze.setQuadrants(center.X(), center.Y())

	// Add spots taking care of the corresponding quadrant for every spot
	for _, spot := range spots {
		maze.addSpot(spot)
	}

	// Add paths to the maze taking care about source/target spots exist (fail if try to create orphan path)
	for _, path := range paths {
		if ok := maze.addPath(path.Origin, path.Destiny); !ok {
			return "", errors.New("could not add path, spot not found")
		}
	}

	// Save maze to database
	if err := s.db.PutMaze(ctx, maze); err != nil {
		return "", err
	}

	return maze.Id, nil
}

func (s service) Update(ctx context.Context, mazeId string, center Coordinates, spots []Spot, paths []Path) error {
	maze, err := s.Get(ctx, mazeId)
	if err != nil {
		return err
	}

	// check if center should be moved
	if x, y := maze.getCenter(); center.X() != x || center.Y() != y {
		maze = maze.moveAxes(center.X(), center.Y())
	}

	// check if should add new spots
	if len(spots) != 0 {
		for _, spot := range spots {
			maze.addSpot(spot)
		}
	}

	// check if should add new paths
	// fail if try to create an orphan path
	if len(paths) != 0 {
		for _, path := range paths {
			if ok := maze.addPath(path.Origin, path.Destiny); !ok {
				return errors.New("could not add path, spot not found")
			}
		}
	}

	return s.db.UpdateMaze(ctx, maze)
}

func (s service) Delete(ctx context.Context, mazeId string) error {
	return s.db.DeleteMaze(ctx, mazeId)
}

func (s service) DeleteSpot(ctx context.Context, mazeId string, coordinate Coordinates) error {
	maze, err := s.Get(ctx, mazeId)
	if err != nil {
		return err
	}

	if _, ok := maze.FindSpot(coordinate.Key()); !ok {
		return errors.New("not found")
	}

	// deletes the spot and all the related paths, so it will not allow orphan paths
	maze.deleteSpot(coordinate)

	return s.db.UpdateMaze(ctx, maze)
}

func (s service) DeletePath(ctx context.Context, mazeId string, path Path) error {
	maze, err := s.Get(ctx, mazeId)
	if err != nil {
		return err
	}

	// deletes the path and the corresponding reverse path
	maze.Paths.deletePath(path.Origin, path.Destiny)

	return s.db.UpdateMaze(ctx, maze)
}

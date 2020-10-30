package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/maxidelgado/maze-api/domain/game"
	"github.com/maxidelgado/maze-api/domain/maze"
)

func NewGame(mazeSvc maze.Service, db game.DataBase) game.Service {
	return &gameSvc{mazeSvc: mazeSvc, db: db}
}

type gameSvc struct {
	mazeSvc maze.Service
	db      game.DataBase
}

func (s gameSvc) Start(ctx context.Context, mazeId string) (game.Game, error) {
	// get the maze
	m, err := s.mazeSvc.Get(ctx, mazeId)
	if err != nil {
		return game.Game{}, err
	}

	// validate if the maze is able to be played
	valid, distance, entrance, exit := validateMaze(m)
	if !valid {
		return game.Game{}, errors.New("the selected Maze is not ready to be played")
	}

	// get the spots connected to the entrance spot
	neighbours := m.GetNeighbours(entrance.Key())
	var allowedMovements []string
	for key := range neighbours {
		allowedMovements = append(allowedMovements, key)
	}

	g := game.Game{
		Id:              uuid.New().String(),
		Entrance:        entrance.Key(),
		Exit:            exit.Key(),
		MinimumDistance: distance,
		Maze:            m,
		StartDate:       time.Now(),
		PlayerStats: game.PlayerStats{
			CurrentSpot:      entrance.Key(),
			AllowedMovements: allowedMovements,
		},
	}

	// persist the game
	err = s.db.PutGame(ctx, g)
	if err != nil {
		return game.Game{}, err
	}

	return g, nil
}

func (s gameSvc) Get(ctx context.Context, gameId string) (game.Game, error) {
	return s.db.GetGame(ctx, gameId)
}

func (s gameSvc) Move(ctx context.Context, gameId string, nextSpot string) (game.Game, error) {
	// get the current game
	g, err := s.db.GetGame(ctx, gameId)
	if err != nil {
		return game.Game{}, err
	}

	// check if the game is already finished
	if !g.EndDate.IsZero() {
		return g, nil
	}

	// check if the selected spot is connected to the current one
	var canMove bool
	for _, allowedMovement := range g.PlayerStats.AllowedMovements {
		if nextSpot == allowedMovement {
			canMove = true
			break
		}
	}

	switch {
	case !canMove:
		// if the selected spot is not connected to the current, return an error
		return game.Game{}, errors.New("could not move to the selected spot")
	case nextSpot == g.Exit:
		// if the selected spot is the exit spot
		g.EndDate = time.Now()
		g.SetAllowedMovements(nil)
	default:
		// if the selected spot is not final and is valid
		neighbours := g.Maze.GetNeighbours(nextSpot)
		var movements []string
		for key := range neighbours {
			movements = append(movements, key)
		}
		g.SetAllowedMovements(movements)
	}

	// ensure that we add gold only the first time
	if !g.HasVisited(nextSpot) {
		g.AddGold(nextSpot)
	}

	g.Move(nextSpot)
	g.AddDistance(nextSpot)
	g.SetCurrentSpot(nextSpot)

	return g, s.db.UpdateGame(ctx, g)
}

func (s gameSvc) Delete(ctx context.Context, gameId string) error {
	return s.db.DeleteGame(ctx, gameId)
}

// validate if the maze is well-formed:
// 	- has entrance and exit spots
//	- both entrance and exit are connected
//	- return the minimum distance required to end the game
func validateMaze(m maze.Maze) (valid bool, minDistance float64, entrance maze.Coordinates, exit maze.Coordinates) {
	for _, quadrant := range m.Quadrants {
		for _, spot := range quadrant.Spots {
			switch spot.Name {
			case game.EntranceSpot:
				entrance = spot.Coordinate
			case game.ExitSpot:
				exit = spot.Coordinate
			}
		}
	}

	minDistance, _ = m.GetPath(entrance.Key(), exit.Key())
	valid = minDistance > 0
	return
}

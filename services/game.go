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

func (s gameSvc) Query(ctx context.Context, name string) ([]game.Game, error) {
	return s.db.QueryGames(ctx, name)
}

func (s gameSvc) Start(ctx context.Context, mazeId, name string) (game.Game, error) {
	if name == "" {
		return game.Game{}, errors.New("game name must be provided")
	}

	// get the maze
	m, err := s.mazeSvc.Get(ctx, mazeId)
	if err != nil {
		return game.Game{}, err
	}

	// validate if the maze is able to be played
	valid, distance, entrance, _ := validateMaze(m)
	if !valid {
		return game.Game{}, errors.New("the selected Maze is not ready to be played")
	}

	// get the spots connected to the entrance spot
	allowedMovements := m.GetAllowedMovements(entrance)

	g := game.Game{
		Id:              uuid.New().String(),
		Name:            name,
		MinimumDistance: distance,
		Maze:            m,
		StartDate:       time.Now(),
		PlayerStats: game.PlayerStats{
			CurrentSpot:      entrance,
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
		if nextSpot == allowedMovement.Key {
			canMove = true
			break
		}
	}

	switch {
	case !canMove:
		// if the selected spot is not connected to the current, return an error
		return game.Game{}, errors.New("could not move to the selected spot")
	case nextSpot == g.Maze.Exit:
		// if the selected spot is the exit spot
		g.EndDate = time.Now()
		g.SetAllowedMovements(nil)
		_, g.OptimumPath = g.Maze.GetPath(g.Maze.Entrance, g.Maze.Exit)
	default:
		// if the selected spot is not final and is valid
		allowedMovements := g.Maze.GetAllowedMovements(nextSpot)
		g.SetAllowedMovements(allowedMovements)
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
func validateMaze(m maze.Maze) (valid bool, minDistance float64, entrance, exit string) {
	if m.Entrance == "" || m.Exit == "" {
		return
	}

	entranceSpot, found := m.FindSpot(m.Entrance)
	if !found {
		return
	}

	exitSpot, found := m.FindSpot(m.Exit)
	if !found {
		return
	}

	minDistance, _ = m.GetPath(m.Entrance, m.Exit)
	valid = minDistance > 0
	entrance = entranceSpot.Coordinate.Key()
	exit = exitSpot.Coordinate.Key()
	return
}

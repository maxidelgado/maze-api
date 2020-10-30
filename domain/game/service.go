package game

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/maxidelgado/maze-api/domain/maze"
	"time"
)

func NewService(mazeSvc maze.Service, db DataBase) Service {
	return &service{mazeSvc: mazeSvc, db: db}
}

type service struct {
	mazeSvc maze.Service
	db      DataBase
}

func (s service) Start(ctx context.Context, mazeId string) (Game, error) {
	// get the maze
	m, err := s.mazeSvc.Get(ctx, mazeId)
	if err != nil {
		return Game{}, err
	}

	// validate if the maze is able to be played
	valid, distance, entrance, exit := validateMaze(m)
	if !valid {
		return Game{}, errors.New("the selected Maze is not ready to be played")
	}

	// get the spots connected to the entrance spot
	neighbours := m.GetNeighbours(entrance.Key())
	var allowedMovements []string
	for key := range neighbours {
		allowedMovements = append(allowedMovements, key)
	}

	game := Game{
		Id:              uuid.New().String(),
		Entrance:        entrance.Key(),
		Exit:            exit.Key(),
		MinimumDistance: distance,
		Maze:            m,
		StartDate:       time.Now(),
		PlayerStats: PlayerStats{
			CurrentSpot:      entrance.Key(),
			AllowedMovements: allowedMovements,
		},
	}

	// persist the game
	err = s.db.PutGame(ctx, game)
	if err != nil {
		return Game{}, err
	}

	return game, nil
}

func (s service) Get(ctx context.Context, gameId string) (Game, error) {
	return s.db.GetGame(ctx, gameId)
}

func (s service) Move(ctx context.Context, gameId string, nextSpot string) (Game, error) {
	// get the current game
	game, err := s.db.GetGame(ctx, gameId)
	if err != nil {
		return Game{}, err
	}

	// check if the game is already finished
	if !game.EndDate.IsZero() {
		return game, nil
	}

	// check if the selected spot is connected to the current one
	var canMove bool
	for _, allowedMovement := range game.PlayerStats.AllowedMovements {
		if nextSpot == allowedMovement {
			canMove = true
			break
		}
	}

	switch {
	case !canMove:
		// if the selected spot is not connected to the current, return an error
		return Game{}, errors.New("could not move to the selected spot")
	case nextSpot == game.Exit:
		// if the selected spot is the exit spot
		game.EndDate = time.Now()
		game.setAllowedMovements(nil)
	default:
		// if the selected spot is not final and is valid
		neighbours := game.Maze.GetNeighbours(nextSpot)
		var movements []string
		for key := range neighbours {
			movements = append(movements, key)
		}
		game.setAllowedMovements(movements)
	}

	// ensure that we add gold only the first time
	if !game.hasVisited(nextSpot) {
		game.addGold(nextSpot)
	}

	game.move(nextSpot)
	game.addDistance(nextSpot)
	game.setCurrentSpot(nextSpot)

	return game, s.db.UpdateGame(ctx, game)
}

func (s service) Delete(ctx context.Context, gameId string) error {
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
			case EntranceSpot:
				entrance = spot.Coordinate
			case ExitSpot:
				exit = spot.Coordinate
			}
		}
	}

	minDistance, _ = m.GetPath(entrance.Key(), exit.Key())
	valid = minDistance > 0
	return
}

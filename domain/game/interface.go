package game

import (
	"context"
)

type Service interface {
	Start(context.Context, string, string) (Game, error)
	Get(context.Context, string) (Game, error)
	Move(context.Context, string, string) (Game, error)
	Delete(context.Context, string) error
	Query(ctx context.Context, name string) ([]Game, error)
}

type DataBase interface {
	GetGame(context.Context, string) (Game, error)
	PutGame(context.Context, Game) error
	UpdateGame(context.Context, Game) error
	DeleteGame(context.Context, string) error
	QueryGames(context.Context, string) ([]Game, error)
}

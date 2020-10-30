package game

import (
	"context"
)

type Service interface {
	Start(context.Context, string) (Game, error)
	Get(context.Context, string) (Game, error)
	Move(context.Context, string, string) (Game, error)
	Delete(context.Context, string) error
}

type DataBase interface {
	GetGame(context.Context, string) (Game, error)
	PutGame(context.Context, Game) error
	UpdateGame(context.Context, Game) error
	DeleteGame(context.Context, string) error
}
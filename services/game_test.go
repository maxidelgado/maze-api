package services

import (
	"context"
	"errors"
	"github.com/maxidelgado/maze-api/domain/game"
	"github.com/maxidelgado/maze-api/domain/maze"
	"testing"
	"time"
)

/*
IMPORTANT: Just added some different kind of test in order to demonstrate the usage of Unit Test
		   The purpose is NOT to achieve a high coverage percentage.

		   In this case the important part is about the mocking of different layers of the project.
			For example the data access layer.
*/
type mazeMock struct {
	maze.Service
	get func(ctx context.Context, id string) (maze.Maze, error)
}

func (m mazeMock) Get(ctx context.Context, id string) (maze.Maze, error) { return m.get(ctx, id) }

type dbMock struct {
	get    func(context.Context, string) (game.Game, error)
	put    func(context.Context, game.Game) error
	update func(context.Context, game.Game) error
	delete func(context.Context, string) error
	query  func(context.Context, string) ([]game.Game, error)
}

func (d dbMock) GetGame(ctx context.Context, id string) (game.Game, error) { return d.get(ctx, id) }
func (d dbMock) PutGame(ctx context.Context, g game.Game) error            { return d.put(ctx, g) }
func (d dbMock) UpdateGame(ctx context.Context, g game.Game) error         { return d.update(ctx, g) }
func (d dbMock) DeleteGame(ctx context.Context, id string) error           { return d.delete(ctx, id) }
func (d dbMock) QueryGames(ctx context.Context, name string) ([]game.Game, error) {
	return d.query(ctx, name)
}

func Test_service_Move(t *testing.T) {
	type fields struct {
		mazeSvc maze.Service
		db      game.DataBase
	}
	type args struct {
		ctx      context.Context
		gameId   string
		nextSpot string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				mazeSvc: mazeMock{
					get: func(ctx context.Context, id string) (maze.Maze, error) {
						return maze.Maze{}, nil
					},
				},
				db: dbMock{
					get: func(ctx context.Context, s string) (game.Game, error) {
						return game.Game{PlayerStats: game.PlayerStats{AllowedMovements: []maze.Neighbour{{Key: "[1,1]"}}}}, nil
					},
					update: func(ctx context.Context, g game.Game) error {
						return nil
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				gameId:   "id",
				nextSpot: "[1,1]",
			},
			wantErr: false,
		},
		{
			name: "success: game is already finished",
			fields: fields{
				mazeSvc: mazeMock{
					get: func(ctx context.Context, id string) (maze.Maze, error) {
						return maze.Maze{}, nil
					},
				},
				db: dbMock{
					get: func(ctx context.Context, s string) (game.Game, error) {
						return game.Game{EndDate: time.Now()}, nil
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				gameId:   "id",
				nextSpot: "[1,1]",
			},
			wantErr: false,
		},
		{
			name: "success: exit",
			fields: fields{
				mazeSvc: mazeMock{
					get: func(ctx context.Context, id string) (maze.Maze, error) {
						return maze.Maze{}, nil
					},
				},
				db: dbMock{
					get: func(ctx context.Context, s string) (game.Game, error) {
						return game.Game{PlayerStats: game.PlayerStats{AllowedMovements: []maze.Neighbour{{Key: "[1,1]"}}}}, nil
					},
					update: func(ctx context.Context, g game.Game) error {
						return nil
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				gameId:   "id",
				nextSpot: "[1,1]",
			},
			wantErr: false,
		},
		{
			name: "fail: get game error",
			fields: fields{
				mazeSvc: mazeMock{
					get: func(ctx context.Context, id string) (maze.Maze, error) {
						return maze.Maze{}, nil
					},
				},
				db: dbMock{
					get: func(ctx context.Context, s string) (game.Game, error) {
						return game.Game{}, errors.New("error")
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				gameId:   "id",
				nextSpot: "[1,1]",
			},
			wantErr: true,
		},
		{
			name: "fail: movement not allowed",
			fields: fields{
				mazeSvc: mazeMock{
					get: func(ctx context.Context, id string) (maze.Maze, error) {
						return maze.Maze{}, nil
					},
				},
				db: dbMock{
					get: func(ctx context.Context, s string) (game.Game, error) {
						return game.Game{}, nil
					},
					update: func(ctx context.Context, g game.Game) error {
						return nil
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				gameId:   "id",
				nextSpot: "[1,1]",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := gameSvc{
				mazeSvc: tt.fields.mazeSvc,
				db:      tt.fields.db,
			}
			_, err := s.Move(tt.args.ctx, tt.args.gameId, tt.args.nextSpot)
			if (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

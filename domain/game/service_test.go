package game

import (
	"context"
	"errors"
	"github.com/maxidelgado/maze-api/domain/maze"
	"testing"
	"time"
)

type mazeMock struct {
	maze.Service
	get func(ctx context.Context, id string) (maze.Maze, error)
}

func (m mazeMock) Get(ctx context.Context, id string) (maze.Maze, error) { return m.get(ctx, id) }

type dbMock struct {
	get    func(context.Context, string) (Game, error)
	put    func(context.Context, Game) error
	update func(context.Context, Game) error
	delete func(context.Context, string) error
}

func (d dbMock) GetGame(ctx context.Context, id string) (Game, error) { return d.get(ctx, id) }
func (d dbMock) PutGame(ctx context.Context, game Game) error         { return d.put(ctx, game) }
func (d dbMock) UpdateGame(ctx context.Context, game Game) error      { return d.update(ctx, game) }
func (d dbMock) DeleteGame(ctx context.Context, id string) error      { return d.delete(ctx, id) }

func Test_service_Move(t *testing.T) {
	type fields struct {
		mazeSvc maze.Service
		db      DataBase
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
					get: func(ctx context.Context, s string) (Game, error) {
						return Game{PlayerStats: PlayerStats{AllowedMovements: []string{"[1,1]"}}}, nil
					},
					update: func(ctx context.Context, game Game) error {
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
					get: func(ctx context.Context, s string) (Game, error) {
						return Game{EndDate: time.Now()}, nil
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
					get: func(ctx context.Context, s string) (Game, error) {
						return Game{Exit: "[1,1]", PlayerStats: PlayerStats{AllowedMovements: []string{"[1,1]"}}}, nil
					},
					update: func(ctx context.Context, game Game) error {
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
					get: func(ctx context.Context, s string) (Game, error) {
						return Game{}, errors.New("error")
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
					get: func(ctx context.Context, s string) (Game, error) {
						return Game{}, nil
					},
					update: func(ctx context.Context, game Game) error {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
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

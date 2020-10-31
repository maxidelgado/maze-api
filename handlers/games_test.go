package handlers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/maxidelgado/maze-api/domain/game"
	"io"
	"net/http"
	"strings"
	"testing"
)

/*
IMPORTANT: Just added some different kind of test in order to demonstrate the usage of Unit Test
		   The purpose is NOT to achieve a high coverage percentage.

		   In this case the important part is about the usage of the package http and fiber
			to simulate HTTP calls to the API
*/
func Test_gamesHandler_postGame(t *testing.T) {
	type fields struct {
		svc game.Service
	}
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "success: game created",
			fields: fields{
				svc: gamesSvcMock{start: func(context.Context, string, string) (game.Game, error) {
					return game.Game{}, nil
				}},
			},
			args: args{
				raw: `{"maze_id":"id"}`,
			},
			want:    http.StatusOK,
			wantErr: false,
		},
		{
			name:   "fail: wrong body",
			fields: fields{},
			args: args{
				raw: `not valid body`,
			},
			want:    http.StatusBadRequest,
			wantErr: false,
		},
		{
			name: "fail: service error",
			fields: fields{
				svc: gamesSvcMock{
					start: func(context.Context, string, string) (game.Game, error) {
						return game.Game{}, errors.New("error")
					},
				},
			},
			args: args{
				raw: `{"maze_id":"id"}`,
			},
			want:    http.StatusInternalServerError,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := doRequest("/games", http.MethodPost, tt.fields.svc, strings.NewReader(tt.args.raw))
			if (err != nil) != tt.wantErr {
				t.Errorf("postSales() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != nil && got.StatusCode != tt.want {
				t.Errorf("postSales() got = %v, want %v", err, tt.want)
			}
		})
	}
}

func doRequest(url, method string, svc game.Service, reader io.Reader) (*http.Response, error) {
	app := fiber.New()
	NewGames(app, svc)

	req, _ := http.NewRequest(
		method,
		url,
		reader,
	)
	req.Header.Add("Content-Type", "application/json")
	return app.Test(req, -1)
}

type gamesSvcMock struct {
	game.Service
	start func(ctx context.Context, mazeId, name string) (game.Game, error)
}

func (s gamesSvcMock) Start(ctx context.Context, mazeId, name string) (game.Game, error) {
	return s.start(ctx, mazeId, name)
}

package database

import (
	"context"
	"time"

	"github.com/maxidelgado/maze-api/config"
	"github.com/maxidelgado/maze-api/database/mgo"
	"github.com/maxidelgado/maze-api/domain/game"
	"github.com/maxidelgado/maze-api/domain/maze"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongodb = mgo.WithContext
)

type Repository interface {
	maze.DataBase
	game.DataBase
}

func New() Repository {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DB.Uri))
	if err != nil {
		panic(err)
	}

	mazeColl := client.Database(config.DB.Database).Collection(config.DB.MazeCollection)
	gameColl := client.Database(config.DB.Database).Collection(config.DB.GameCollection)

	return database{
		mazeColl: mazeColl,
		gameColl: gameColl,
	}
}

type database struct {
	mazeColl *mongo.Collection
	gameColl *mongo.Collection
}

func (d database) GetGame(ctx context.Context, id string) (game.Game, error) {
	var result game.Game
	err := mongodb(ctx).Get(d.gameColl, id, &result)
	return result, err
}

func (d database) PutGame(ctx context.Context, game game.Game) error {
	return mongodb(ctx).Put(d.gameColl, game)
}

func (d database) UpdateGame(ctx context.Context, game game.Game) error {
	return mongodb(ctx).Update(d.gameColl, game.Id, game)
}

func (d database) DeleteGame(ctx context.Context, id string) error {
	return mongodb(ctx).DeleteDocument(d.gameColl, id)
}

func (d database) PutMaze(ctx context.Context, maze maze.Maze) error {
	return mongodb(ctx).Put(d.mazeColl, maze)
}

func (d database) UpdateMaze(ctx context.Context, maze maze.Maze) error {
	return mongodb(ctx).Update(d.mazeColl, maze.Id, maze)
}

func (d database) GetMaze(ctx context.Context, id string) (maze.Maze, error) {
	var result maze.Maze
	err := mongodb(ctx).Get(d.mazeColl, id, &result)
	return result, err
}

func (d database) DeleteMaze(ctx context.Context, id string) error {
	return mongodb(ctx).DeleteDocument(d.mazeColl, id)
}

package database

import (
	"context"
	"time"

	"github.com/maxidelgado/maze-api/config"
	"github.com/maxidelgado/maze-api/database/mgo"
	"github.com/maxidelgado/maze-api/domain/game"
	"github.com/maxidelgado/maze-api/domain/maze"
	"go.mongodb.org/mongo-driver/bson"
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

	_, err = gameColl.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"name", "text"}}})
	if err != nil {
		panic(err)
	}

	_, err = mazeColl.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"name", "text"}}})
	if err != nil {
		panic(err)
	}

	return database{
		mazeColl: mazeColl,
		gameColl: gameColl,
	}
}

type database struct {
	mazeColl *mongo.Collection
	gameColl *mongo.Collection
}

func (d database) QueryMaze(ctx context.Context, name string) ([]maze.Maze, error) {
	var result []maze.Maze
	cursor, err := mongodb(ctx).Find(d.mazeColl, name)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var m maze.Maze
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, err
}

func (d database) QueryGames(ctx context.Context, name string) ([]game.Game, error) {
	var result []game.Game
	cursor, err := mongodb(ctx).Find(d.gameColl, name)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var g game.Game
		if err := cursor.Decode(&g); err != nil {
			return nil, err
		}
		result = append(result, g)
	}

	return result, err
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

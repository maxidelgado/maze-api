package database

import (
	"context"
	"github.com/maxidelgado/maze-api/domain/game"
	"time"

	"github.com/maxidelgado/maze-api/config"
	"github.com/maxidelgado/maze-api/domain/maze"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		client:   client,
		mazeColl: mazeColl,
		gameColl: gameColl,
	}
}

type database struct {
	client   *mongo.Client
	mazeColl *mongo.Collection
	gameColl *mongo.Collection
}

func (d database) GetGame(ctx context.Context, id string) (game.Game, error) {
	var result game.Game
	out := d.gameColl.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	if out.Err() != nil {
		return game.Game{}, out.Err()
	}

	err := out.Decode(&result)
	if err != nil {
		return game.Game{}, err
	}

	return result, err
}

func (d database) PutGame(ctx context.Context, game game.Game) error {
	_, err := d.gameColl.InsertOne(ctx, game)
	return err
}

func (d database) UpdateGame(ctx context.Context, game game.Game) error {
	_, err := d.gameColl.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: game.Id}},
		bson.D{{Key: "$set", Value: game}},
	)
	return err
}

func (d database) DeleteGame(ctx context.Context, id string) error {
	_, err := d.gameColl.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}

func (d database) PutMaze(ctx context.Context, maze maze.Maze) error {
	_, err := d.mazeColl.InsertOne(ctx, maze)
	return err
}

func (d database) UpdateMaze(ctx context.Context, maze maze.Maze) error {
	_, err := d.mazeColl.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: maze.Id}},
		bson.D{{Key: "$set", Value: maze}},
	)
	return err
}

func (d database) GetMaze(ctx context.Context, id string) (maze.Maze, error) {
	var result maze.Maze
	out := d.mazeColl.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	if out.Err() != nil {
		return maze.Maze{}, out.Err()
	}

	err := out.Decode(&result)
	if err != nil {
		return maze.Maze{}, err
	}

	return result, err
}

func (d database) DeleteMaze(ctx context.Context, id string) error {
	_, err := d.mazeColl.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}

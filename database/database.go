package database

import (
	"context"
	"fmt"
	"time"

	"github.com/maxidelgado/maze-api/config"
	"github.com/maxidelgado/maze-api/domain/maze"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	maze.DataBase
}

func New() Repository {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DB.Uri))
	if err != nil {
		panic(err)
	}

	collection := client.Database(config.DB.Database).Collection(config.DB.Collection)

	return database{
		client:     client,
		collection: collection,
	}
}

type database struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (d database) PutMaze(ctx context.Context, maze maze.Maze) error {
	_, err := d.collection.InsertOne(ctx, maze)
	return err
}

func (d database) UpdateMaze(ctx context.Context, maze maze.Maze) error {
	_, err := d.collection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: maze.Id}},
		bson.D{{Key: "$set", Value: maze}},
	)
	return err
}

func (d database) GetMaze(ctx context.Context, id string) (maze.Maze, error) {
	var result maze.Maze
	out := d.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
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
	_, err := d.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}

func (d database) DeleteSpot(ctx context.Context, mazeId, quadrantId string, coordinate maze.Coordinate) error {
	key := fmt.Sprintf("quadrants.spots.%s", coordinate.Key())
	_, err := d.collection.UpdateOne(ctx,
		bson.D{
			{Key: "_id", Value: mazeId},
		},
		bson.D{
			{Key: "$pull", Value: bson.E{Key: key, Value: ""}},
		},
	)

	return err
}

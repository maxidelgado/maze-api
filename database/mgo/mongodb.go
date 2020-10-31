package mgo

// isolated this package to facilitate unit testing at repository level

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Get(coll *mongo.Collection, id string, out interface{}) error
	DeleteDocument(coll *mongo.Collection, id string) error
	Update(coll *mongo.Collection, id string, obj interface{}) error
	Put(coll *mongo.Collection, obj interface{}) error
	Find(coll *mongo.Collection, value string) (Cursor, error)
}

type Cursor interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
}

func WithContext(ctx context.Context) Client {
	return mongodb{ctx: ctx}
}

type mongodb struct {
	ctx context.Context
}

func (db mongodb) Get(coll *mongo.Collection, id string, out interface{}) error {
	res := coll.FindOne(db.ctx, bson.D{{Key: "_id", Value: id}})
	err := res.Decode(out)
	if err != nil {
		return err
	}

	return err
}

func (db mongodb) Find(coll *mongo.Collection, value string) (Cursor, error) {
	return coll.Find(db.ctx, bson.D{{Key: "$text", Value: bson.D{{"$search", value}}}})
}

func (db mongodb) DeleteDocument(coll *mongo.Collection, id string) error {
	_, err := coll.DeleteOne(db.ctx, bson.D{{Key: "_id", Value: id}})
	return err
}

func (db mongodb) Update(coll *mongo.Collection, id string, obj interface{}) error {
	_, err := coll.UpdateOne(
		db.ctx,
		bson.D{{Key: "_id", Value: id}},
		bson.D{{Key: "$set", Value: obj}},
	)
	return err
}

func (db mongodb) Put(coll *mongo.Collection, obj interface{}) error {
	_, err := coll.InsertOne(db.ctx, obj)
	return err
}

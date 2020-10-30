package mgo

import "go.mongodb.org/mongo-driver/mongo"

type Mock struct {
	GetFunc    func(coll *mongo.Collection, id string, out interface{}) error
	DeleteFunc func(coll *mongo.Collection, id string) error
	UpdateFunc func(coll *mongo.Collection, id string, obj interface{}) error
	PutFunc    func(coll *mongo.Collection, obj interface{}) error
}

func (m Mock) Get(coll *mongo.Collection, id string, out interface{}) error {
	if m.GetFunc == nil {
		return nil
	}

	return m.GetFunc(coll, id, out)
}

func (m Mock) DeleteDocument(coll *mongo.Collection, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}

	return m.DeleteFunc(coll, id)
}

func (m Mock) Update(coll *mongo.Collection, id string, obj interface{}) error {
	if m.UpdateFunc == nil {
		return nil
	}

	return m.UpdateFunc(coll, id, obj)
}

func (m Mock) Put(coll *mongo.Collection, obj interface{}) error {
	if m.PutFunc == nil {
		return nil
	}

	return m.PutFunc(coll, obj)
}

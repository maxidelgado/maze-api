package database

import (
	"context"
	"errors"
	"github.com/maxidelgado/maze-api/database/mgo"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func createMock(mock mgo.Mock) func(context.Context) mgo.Client {
	return func(context.Context) mgo.Client {
		return mock
	}
}

func Test_database_DeleteGame(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		mgoMock mgo.Mock
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			mgoMock: mgo.Mock{},
			args: args{
				ctx: context.Background(),
				id:  "id",
			},
			wantErr: false,
		},
		{
			name: "fail",
			mgoMock: mgo.Mock{DeleteFunc: func(coll *mongo.Collection, id string) error {
				return errors.New("error")
			}},
			args: args{
				ctx: context.Background(),
				id:  "id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := database{}
			mongodb = createMock(tt.mgoMock)
			if err := d.DeleteGame(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteGame() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

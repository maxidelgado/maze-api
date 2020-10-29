package database

import (
	"github.com/maxidelgado/maze-api/domain/maze"
	"math"
	"testing"
)

const (
	TopLeft     = "top left"
	TopRight    = "top right"
	BottomLeft  = "bottom left"
	BottomRight = "bottom right"

	Infinite = math.MaxInt64
)

func Test_datastore_Insert(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: maze created",
			args: args{
				val: maze.Maze{
					Id:        "1",
					Quadrants: setupQuadrants(0, 0),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			if err := d.Insert(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func setupQuadrants(x, y int64) [4]maze.Quadrant {
	var quadrants [4]maze.Quadrant
	quadrants[0] = maze.Quadrant{
		Id:     TopLeft,
		LimitX: maze.Coordinate{-Infinite, x},
		LimitY: maze.Coordinate{y, Infinite},
	}
	quadrants[1] = maze.Quadrant{
		Id:     TopRight,
		LimitX: maze.Coordinate{x, Infinite},
		LimitY: maze.Coordinate{y, Infinite},
	}
	quadrants[2] = maze.Quadrant{
		Id:     BottomLeft,
		LimitX: maze.Coordinate{-Infinite, x},
		LimitY: maze.Coordinate{-Infinite, y},
	}
	quadrants[3] = maze.Quadrant{
		Id:     BottomRight,
		LimitX: maze.Coordinate{x, Infinite},
		LimitY: maze.Coordinate{-Infinite, y},
	}

	return quadrants
}

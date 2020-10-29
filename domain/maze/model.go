package maze

import (
	"fmt"
	"math"
)

const (
	TopLeft     = "top left"
	TopRight    = "top right"
	BottomLeft  = "bottom left"
	BottomRight = "bottom right"

	TopLeftIndex     = 0
	TopRightIndex    = 1
	BottomLeftIndex  = 2
	BottomRightIndex = 3

	Infinite = math.MaxInt64
)

type Maze struct {
	Id        string          `json:"id" bson:"_id"`
	Quadrants [4]Quadrant     `json:"quadrants" bson:"quadrants"`
	Paths     map[string]Path `json:"paths" bson:"paths"`
}

func (m *Maze) setQuadrants(x, y int64) {
	m.Quadrants[TopLeftIndex] = Quadrant{
		Id:     TopLeft,
		LimitX: Coordinate{-Infinite, x},
		LimitY: Coordinate{y, Infinite},
	}
	m.Quadrants[TopRightIndex] = Quadrant{
		Id:     TopRight,
		LimitX: Coordinate{x, Infinite},
		LimitY: Coordinate{y, Infinite},
	}
	m.Quadrants[BottomLeftIndex] = Quadrant{
		Id:     BottomLeft,
		LimitX: Coordinate{-Infinite, x},
		LimitY: Coordinate{-Infinite, y},
	}
	m.Quadrants[BottomRightIndex] = Quadrant{
		Id:     BottomRight,
		LimitX: Coordinate{x, Infinite},
		LimitY: Coordinate{-Infinite, y},
	}
}

func (m *Maze) addSpot(spot Spot) {
	quadrant := m.Quadrants[TopLeftIndex]
	isLeft := spot.Coordinate.X() <= quadrant.LimitX.Y()
	isTop := spot.Coordinate.Y() >= quadrant.LimitY.X()

	switch {
	case isTop && isLeft:
		if m.Quadrants[TopLeftIndex].Spots == nil {
			m.Quadrants[TopLeftIndex].Spots = make(map[string]Spot)
		}
		m.Quadrants[TopLeftIndex].Spots[spot.Coordinate.Key()] = spot
	case isTop && !isLeft:
		if m.Quadrants[TopRightIndex].Spots == nil {
			m.Quadrants[TopRightIndex].Spots = make(map[string]Spot)
		}
		m.Quadrants[TopRightIndex].Spots[spot.Coordinate.Key()] = spot
	case !isTop && isLeft:
		if m.Quadrants[BottomLeftIndex].Spots == nil {
			m.Quadrants[BottomLeftIndex].Spots = make(map[string]Spot)
		}
		m.Quadrants[BottomLeftIndex].Spots[spot.Coordinate.Key()] = spot
	case !isTop && !isLeft:
		if m.Quadrants[BottomRightIndex].Spots == nil {
			m.Quadrants[BottomRightIndex].Spots = make(map[string]Spot)
		}
		m.Quadrants[BottomRightIndex].Spots[spot.Coordinate.Key()] = spot
	}
}

func (m *Maze) findSpot(coordinate Coordinate) bool {
	for _, quadrant := range m.Quadrants {
		_, ok := quadrant.Spots[coordinate.Key()]
		if ok {
			return true
		}
	}

	return false
}

func (m *Maze) addPath(path Path) {
	if m.Paths == nil {
		m.Paths = make(map[string]Path)
	}
	_, ok := m.Paths[path.Key()]
	if ok {
		return
	}
	path.calculateDistance()
	m.Paths[path.Key()] = path
}

func (m Maze) moveAxes(x, y int64) Maze {
	var maze Maze

	maze.Id = m.Id
	maze.Paths = m.Paths
	maze.setQuadrants(x, y)

	for _, quadrant := range m.Quadrants {
		for _, spot := range quadrant.Spots {
			maze.addSpot(spot)
		}
	}

	return maze
}

type Quadrant struct {
	Id     string          `json:"id" bson:"_id"`
	LimitX Coordinate      `json:"limit_x" bson:"limit_x"`
	LimitY Coordinate      `json:"limit_y" bson:"limit_y"`
	Spots  map[string]Spot `json:"spots" bson:"spots"`
}

type Spot struct {
	Id         string     `json:"id" bson:"_id"`
	Coordinate Coordinate `json:"coordinate" bson:"coordinate"`
	Name       string     `json:"name" bson:"name"`
	Gold       int        `json:"gold" bson:"gold"`
}

type Path struct {
	Source   Coordinate `json:"source" bson:"source"`
	Target   Coordinate `json:"target" bson:"target"`
	Distance float64    `json:"distance" bson:"distance"`
}

func (p Path) Key() string {
	return p.Source.Key() + "-" + p.Target.Key()
}

func (p *Path) calculateDistance() {
	a := math.Pow(float64(p.Source.X()-p.Target.X()), 2)
	b := math.Pow(float64(p.Source.Y()-p.Target.Y()), 2)
	p.Distance = math.Sqrt(a + b)
}

type Coordinate [2]int64

func (c Coordinate) X() int64 {
	return c[0]
}

func (c Coordinate) Y() int64 {
	return c[1]
}

func (c Coordinate) Key() string {
	return fmt.Sprintf("(%v,%v)", c.X(), c.Y())
}

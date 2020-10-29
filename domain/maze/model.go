package maze

import (
	"fmt"
	"math"
	"strings"
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

func (m *Maze) getCoordinateQuadrant(coordinate Coordinate) (id string, index int) {
	quadrant := m.Quadrants[TopLeftIndex]
	isLeft := coordinate.X() <= quadrant.LimitX.Y()
	isTop := coordinate.Y() >= quadrant.LimitY.X()

	switch {
	case isTop && isLeft:
		id = TopLeft
		index = TopLeftIndex
		return
	case isTop && !isLeft:
		id = TopRight
		index = TopRightIndex
		return
	case !isTop && isLeft:
		id = BottomLeft
		index = BottomLeftIndex
		return
	case !isTop && !isLeft:
		id = BottomRight
		index = BottomRightIndex
		return
	default:
		return "", 0
	}
}

func (m *Maze) addSpot(spot Spot) {
	_, index := m.getCoordinateQuadrant(spot.Coordinate)
	if m.Quadrants[index].Spots == nil {
		m.Quadrants[index].Spots = make(map[string]Spot)
	}
	m.Quadrants[index].Spots[spot.Coordinate.Key()] = spot
}

func (m *Maze) deleteSpot(coordinate Coordinate) {
	_, index := m.getCoordinateQuadrant(coordinate)
	if m.Quadrants[index].Spots == nil {
		m.Quadrants[index].Spots = make(map[string]Spot)
	}
	delete(m.Quadrants[index].Spots, coordinate.Key())
	for key := range m.Paths {
		if strings.Contains(key, coordinate.Key()) {
			delete(m.Paths, key)
		}
	}
}

func (m *Maze) deletePath(path Path) {
	delete(m.Paths, path.Key())
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

func (m *Maze) addPath(path Path) bool {
	if m.Paths == nil {
		m.Paths = make(map[string]Path)
	}

	// check if both spots already exist in the maze
	if !(m.findSpot(path.Edge[0]) && m.findSpot(path.Edge[1])) {
		return false
	}

	// check if path already exist
	_, ok := m.Paths[path.Key()]
	if ok {
		return true
	}

	// add a new path
	path.calculateDistance()
	m.Paths[path.Key()] = path
	return true
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

func (m Maze) getCenter() (int64, int64) {
	x := m.Quadrants[0].LimitX.Y()
	y := m.Quadrants[0].LimitY.X()

	return x, y
}

type Quadrant struct {
	Id     string          `json:"id" bson:"_id"`
	LimitX Coordinate      `json:"limit_x" bson:"limit_x"`
	LimitY Coordinate      `json:"limit_y" bson:"limit_y"`
	Spots  map[string]Spot `json:"spots" bson:"spots"`
}

type Spot struct {
	Name       string     `json:"name" bson:"name"`
	Coordinate Coordinate `json:"coordinate" bson:"coordinate"`
	GoldAmount int        `json:"gold_amount" bson:"gold_amount"`
}

type Path struct {
	Edge     [2]Coordinate `json:"edge" bson:"edge"`
	Distance float64       `json:"distance" bson:"distance"`
}

func (p Path) Key() string {
	return p.Edge[0].Key() + "-" + p.Edge[1].Key()
}

func (p *Path) calculateDistance() {
	a := math.Pow(float64(p.Edge[0].X()-p.Edge[1].X()), 2)
	b := math.Pow(float64(p.Edge[0].Y()-p.Edge[1].Y()), 2)
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
	return fmt.Sprintf("[%v,%v]", c.X(), c.Y())
}

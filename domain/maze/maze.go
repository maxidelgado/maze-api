package maze

import (
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
	Id        string      `json:"id" bson:"_id"`
	Quadrants [4]Quadrant `json:"quadrants"`
	Paths     PathsIndex  `json:"paths"`
}

// create the quadrants of the maze based on a central point in the cartesian plane
func (m *Maze) setQuadrants(x, y int64) {
	m.Quadrants = createQuadrants(x, y)
}

// returns the id and the index to get the plane where a coordinate is contained
func (m *Maze) getCoordinateQuadrant(coordinate Coordinates) (id string, index int) {
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

// add a spot to the maze in the corresponding quadrant
func (m *Maze) addSpot(spot Spot) {
	_, index := m.getCoordinateQuadrant(spot.Coordinate)
	m.Quadrants[index].Spots[spot.Coordinate.Key()] = spot
}

// delete a spot from the maze and produces a cascade deleting of all the related paths
func (m *Maze) deleteSpot(coordinate Coordinates) {
	_, index := m.getCoordinateQuadrant(coordinate)
	if m.Quadrants[index].Spots == nil {
		return
	}

	delete(m.Quadrants[index].Spots, coordinate.Key())

	paths := m.Paths[coordinate.Key()]

	for key := range paths {
		delete(m.Paths[key], coordinate.Key())
	}

	delete(m.Paths, coordinate.Key())
}

// find a spot by key
func (m *Maze) FindSpot(key string) (Spot, bool) {
	for _, quadrant := range m.Quadrants {
		if spot, ok := quadrant.Spots[key]; ok {
			return spot, true
		}
	}

	return Spot{}, false
}

// add an edge between two existing spots
func (m *Maze) addPath(origin, destiny Coordinates) bool {
	// check if both spots already exist in the maze
	_, originFound := m.FindSpot(origin.Key())
	_, destinyFound := m.FindSpot(destiny.Key())
	if !(originFound && destinyFound) {
		return false
	}

	if m.Paths[origin.Key()] == nil {
		m.Paths[origin.Key()] = map[string]float64{}
	}

	m.Paths.appendPath(origin, destiny)
	return true
}

// allows to change the central point of the entire maze, and moves all the spots to the corresponding quadrant
func (m *Maze) moveAxes(x, y int64) Maze {
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

// returns the central point of the maze
func (m *Maze) getCenter() (int64, int64) {
	x := m.Quadrants[0].LimitX.Y()
	y := m.Quadrants[0].LimitY.X()

	return x, y
}

// returns all the spots that are directly connected to the current spot
func (m *Maze) GetNeighbours(origin string) map[string]float64 {
	return m.Paths[origin]
}

// uses the Dijkstra algorithm to verify if two spots are already connected through any path,
// and returns the distance between them and the complete path
func (m *Maze) GetPath(origin, destiny string) (float64, []string) {
	h := newHeap()
	h.push(path{value: 0, nodes: []string{origin}})
	visited := make(map[string]bool)

	for len(*h.values) > 0 {
		// Find the nearest yet to visit node
		p := h.pop()
		node := p.nodes[len(p.nodes)-1]

		if visited[node] {
			continue
		}

		if node == destiny {
			return p.value, p.nodes
		}

		neighbours := m.GetNeighbours(node)
		for k, distance := range neighbours {
			if !visited[destiny] {
				// We calculate the total spent so far plus the cost and the path of getting here
				h.push(path{value: p.value + distance, nodes: append([]string{}, append(p.nodes, k)...)})
			}
		}

		visited[node] = true
	}

	return 0, nil
}

package maze

import (
	"errors"
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
	Entrance  string      `json:"-"`
	Exit      string      `json:"-"`
	Quadrants [4]Quadrant `json:"quadrants"`
	Paths     PathsIndex  `json:"paths"`
}

// Create the quadrants of the maze based on a central point in the cartesian plane - Default: [0,0]
func (m *Maze) SetQuadrants(x, y int64) {
	m.Quadrants = createQuadrants(x, y)
}

// Returns the name (top left, bottom right, etc) and the index of the quadrant which contains a given coordinate
func (m *Maze) getCoordinateQuadrant(coordinate Coordinates) (id string, index int) {
	quadrant := m.Quadrants[TopLeftIndex] // takes the top-left quadrant as a reference
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

// Add a spot to the corresponding quadrant in a maze
func (m *Maze) AddSpot(spot Spot) error {
	switch spot.Name {
	case EntranceSpot:
		if m.Entrance != "" {
			return errors.New("multiple entrance spots not allowed")
		}
		m.Entrance = spot.Coordinate.Key()
	case ExitSpot:
		if m.Exit != "" {
			return errors.New("multiple exit spots not allowed")
		}
		m.Exit = spot.Coordinate.Key()
	}

	_, index := m.getCoordinateQuadrant(spot.Coordinate)
	m.Quadrants[index].Spots[spot.Coordinate.Key()] = spot
	return nil
}

// Delete a spot from the maze and produces a cascade deleting of all the related paths to avoid orphan paths
func (m *Maze) DeleteSpot(coordinate Coordinates) {
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

// Check if a spot is present in the maze
func (m *Maze) FindSpot(key string) (Spot, bool) {
	for _, quadrant := range m.Quadrants {
		if spot, ok := quadrant.Spots[key]; ok {
			return spot, true
		}
	}

	return Spot{}, false
}

// Add an edge between two existing spots (and the corresponding reverse-path)
func (m *Maze) AddPath(origin, destiny Coordinates) bool {
	// check if both spots already exist in the maze
	_, originFound := m.FindSpot(origin.Key())
	_, destinyFound := m.FindSpot(destiny.Key())
	if !(originFound && destinyFound) {
		return false
	}

	if m.Paths[origin.Key()] == nil {
		m.Paths[origin.Key()] = map[string]float64{}
	}
	if m.Paths[destiny.Key()] == nil {
		m.Paths[destiny.Key()] = map[string]float64{}
	}

	m.Paths.appendPath(origin, destiny)
	m.Paths.appendPath(destiny, origin) // the reverse path
	return true
}

// Change the central point of the entire maze and moves all the spots to the corresponding quadrant
func (m *Maze) MoveAxes(x, y int64) Maze {
	var maze Maze

	maze.Id = m.Id
	maze.Paths = m.Paths
	maze.SetQuadrants(x, y)

	for _, quadrant := range m.Quadrants {
		for _, spot := range quadrant.Spots {
			_ = maze.AddSpot(spot)
		}
	}

	return maze
}

// Calculate and the central point of the maze
func (m *Maze) GetCenter() (int64, int64) {
	x := m.Quadrants[0].LimitX.Y()
	y := m.Quadrants[0].LimitY.X()

	return x, y
}

// Returns all the spots that are directly connected to the current spot
func (m *Maze) GetNeighbours(origin string) map[string]float64 {
	return m.Paths[origin]
}

/*
	Uses the Dijkstra algorithm to verify if two spots are already connected through any path
	and calculate the minimum distance between them.
	The algorithm is backed by a min-heap implementation.
*/
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

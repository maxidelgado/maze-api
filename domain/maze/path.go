package maze

import (
	"math"
)

// these nested maps provide the ability to bidirectionally index the paths between spots in the following format:
// [a,b]: {[c,d] : distance1, [e,f] : distance2}
type PathsIndex map[string]map[string]float64

// appends a path between two coordinates by doing: sqrt((x1-x2)²+(y1-y2)²)
func (p PathsIndex) appendPath(origin, destiny Coordinates) {
	a := math.Pow(float64(origin.X()-destiny.X()), 2)
	b := math.Pow(float64(origin.Y()-destiny.Y()), 2)
	p[origin.Key()][destiny.Key()] = math.Sqrt(a + b)
}

func (p PathsIndex) deletePath(origin, destiny Coordinates) {
	delete(p[origin.Key()], destiny.Key())
	delete(p[destiny.Key()], origin.Key())
}

type Path struct {
	Origin  Coordinates `json:"origin"`
	Destiny Coordinates `json:"destiny"`
}

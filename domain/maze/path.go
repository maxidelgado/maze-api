package maze

import (
	"math"
)

/*
	The usage of this two-ways index was inspired by the adjacency lists pattern commonly used in DynamoDB,
	described here: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-adjacency-graphs.html

	This in-memory index performs extremely well for querying paths in any direction, like origin -> destiny or
	destiny -> origin, so you can swap [a,b][c,d] to [c,d][a,b] at any time.

	Once we added a path to the in-memory index by uson, it will cal
*/
type PathsIndex map[string]map[string]float64

/*
	Once we add a new path to the index, it will be saved with the distance between both spots.
	I prefer to persist the distance in order to improve the performance when the client asks for the paths,
	and because the distance is immutable while the path exists.

	The distance is calculated as: sqrt((x1-x2)²+(y1-y2)²)
*/
func (p PathsIndex) appendPath(origin, destiny Coordinates) {
	a := math.Pow(float64(origin.X()-destiny.X()), 2)
	b := math.Pow(float64(origin.Y()-destiny.Y()), 2)
	p[origin.Key()][destiny.Key()] = math.Sqrt(a + b)
}

// Deletes both, the path and the reverse-path
func (p PathsIndex) DeletePath(origin, destiny Coordinates) {
	delete(p[origin.Key()], destiny.Key())
	delete(p[destiny.Key()], origin.Key())
}

type Path struct {
	Origin  Coordinates `json:"origin"`
	Destiny Coordinates `json:"destiny"`
}

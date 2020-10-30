package maze

import "fmt"

// Represents a point in a cartesian plane in the format [a,b]
type Coordinates [2]int64

func (c Coordinates) X() int64 {
	return c[0]
}

func (c Coordinates) Y() int64 {
	return c[1]
}

func (c Coordinates) Key() string {
	return fmt.Sprintf("[%v,%v]", c.X(), c.Y())
}

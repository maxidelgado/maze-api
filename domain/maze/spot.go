package maze

// represents a location inside the maze
type Spot struct {
	Name       string      `json:"name"`
	Coordinate Coordinates `json:"coordinate"`
	GoldAmount int         `json:"gold_amount"`
}

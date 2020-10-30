package maze

// create quadrants based on a given point (x,y)
func createQuadrants(x, y int64) [4]Quadrant {
	var quadrants [4]Quadrant
	quadrants[TopLeftIndex] = Quadrant{
		Id:     TopLeft,
		LimitX: Coordinates{-Infinite, x},
		LimitY: Coordinates{y, Infinite},
		Spots:  map[string]Spot{},
	}
	quadrants[TopRightIndex] = Quadrant{
		Id:     TopRight,
		LimitX: Coordinates{x, Infinite},
		LimitY: Coordinates{y, Infinite},
		Spots:  map[string]Spot{},
	}
	quadrants[BottomLeftIndex] = Quadrant{
		Id:     BottomLeft,
		LimitX: Coordinates{-Infinite, x},
		LimitY: Coordinates{-Infinite, y},
		Spots:  map[string]Spot{},
	}
	quadrants[BottomRightIndex] = Quadrant{
		Id:     BottomRight,
		LimitX: Coordinates{x, Infinite},
		LimitY: Coordinates{-Infinite, y},
		Spots:  map[string]Spot{},
	}

	return quadrants
}

type Quadrant struct {
	Id     string          `json:"id" bson:"_id"`
	LimitX Coordinates     `json:"limit_x"`
	LimitY Coordinates     `json:"limit_y"`
	Spots  map[string]Spot `json:"spots"`
}

package maze

const (
	EntranceSpot = "entrance"
	ExitSpot     = "exit"
)

// Represents a location inside the maze with the corresponding name and amount of gold.
type Spot struct {
	Name       string      `json:"name"`
	Coordinate Coordinates `json:"coordinate"`
	GoldAmount int         `json:"gold_amount"`
}

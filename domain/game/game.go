package game

import (
	"time"

	"github.com/maxidelgado/maze-api/domain/maze"
)

const (
	EntranceSpot = "entrance"
	ExitSpot     = "exit"
)

// Represents an in-progress game
type Game struct {
	Id              string      `json:"id" bson:"_id"`
	Entrance        string      `json:"entrance"`
	Exit            string      `json:"exit"`
	MinimumDistance float64     `json:"minimum_distance"`
	PlayerStats     PlayerStats `json:"player_stats"`
	Maze            maze.Maze   `json:"-"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date,omitempty"`
}

// Performs a movement to the spot selected by the player
func (g *Game) Move(selectedSpot string) {
	g.PlayerStats.Movements = append(g.PlayerStats.Movements, Movement{
		Date: time.Now(),
		From: g.PlayerStats.CurrentSpot,
		To:   selectedSpot,
	})
}

// Add the gold found in the selected spot to the player stats
func (g *Game) AddGold(selectedSpot string) {
	spot, _ := g.Maze.FindSpot(selectedSpot)
	g.PlayerStats.TotalGold += spot.GoldAmount
}

// Add the distance from the current spot to the one selected by the player to the stats
func (g *Game) AddDistance(selectedSpot string) {
	g.PlayerStats.DistanceCovered += g.Maze.Paths[g.PlayerStats.CurrentSpot][selectedSpot]
}

// Set the selected spot as the current one
func (g *Game) SetCurrentSpot(selectedSpot string) {
	g.PlayerStats.CurrentSpot = selectedSpot
}

// Set the allowed movements (player can't move to an out of radar spot)
func (g *Game) SetAllowedMovements(movements []string) {
	g.PlayerStats.AllowedMovements = movements
}

// Check if the user already passed by the selected spot
func (g *Game) HasVisited(selectedSpot string) bool {
	for _, movement := range g.PlayerStats.Movements {
		if movement.From == selectedSpot || movement.To == selectedSpot {
			return true
		}
	}

	return false
}

// Represents a moving from one spot to another
type Movement struct {
	Date time.Time `json:"date"`
	From string    `json:"from"`
	To   string    `json:"to"`
}

// Represents the player stats
type PlayerStats struct {
	TotalGold        int        `json:"total_gold"`
	DistanceCovered  float64    `json:"distance_covered"`
	CurrentSpot      string     `json:"current_spot"`
	Movements        []Movement `json:"movements"`
	AllowedMovements []string   `json:"allowed_movements"`
}

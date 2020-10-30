package game

import (
	"time"

	"github.com/maxidelgado/maze-api/domain/maze"
)

const (
	EntranceSpot = "entrance"
	ExitSpot     = "exit"
)

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

func (g *Game) move(next string) {
	g.PlayerStats.Movements = append(g.PlayerStats.Movements, Movement{
		Date: time.Now(),
		From: g.PlayerStats.CurrentSpot,
		To:   next,
	})
}

func (g *Game) addGold(spotId string) {
	spot, _ := g.Maze.FindSpot(spotId)
	g.PlayerStats.TotalGold += spot.GoldAmount
}

func (g *Game) addDistance(spotId string) {
	g.PlayerStats.DistanceCovered += g.Maze.Paths[g.PlayerStats.CurrentSpot][spotId]
}

func (g *Game) setCurrentSpot(spotId string) {
	g.PlayerStats.CurrentSpot = spotId
}

func (g *Game) setAllowedMovements(movements []string) {
	g.PlayerStats.AllowedMovements = movements
}

func (g *Game) hasVisited(spot string) bool {
	for _, movement := range g.PlayerStats.Movements {
		if movement.From == spot || movement.To == spot {
			return true
		}
	}

	return false
}

type Movement struct {
	Date time.Time `json:"date"`
	From string    `json:"from"`
	To   string    `json:"to"`
}

type PlayerStats struct {
	TotalGold        int        `json:"total_gold"`
	DistanceCovered  float64    `json:"distance_covered"`
	CurrentSpot      string     `json:"current_spot"`
	Movements        []Movement `json:"movements"`
	AllowedMovements []string   `json:"allowed_movements"`
}

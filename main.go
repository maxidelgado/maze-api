package main

import (
	"github.com/maxidelgado/maze-api/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/maxidelgado/maze-api/config"
	"github.com/maxidelgado/maze-api/database"
	"github.com/maxidelgado/maze-api/handlers"
)

func main() {
	// setup router
	app := fiber.New()
	app.Use(
		recover.New(),
	)
	api := app.Group(config.Router.BasePath)

	// setup repositories
	db := database.New()

	// setup services
	mazeSvc := services.NewMaze(db)
	gameSvc := services.NewGame(mazeSvc, db)

	// setup handlers
	handlers.NewMaze(api, mazeSvc)
	handlers.NewGames(api, gameSvc)

	log.Fatal(app.Listen(config.Router.Host))
}

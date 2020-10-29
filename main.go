package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/maxidelgado/maze-api/config"
	"github.com/maxidelgado/maze-api/database"
	"github.com/maxidelgado/maze-api/domain/maze"
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
	mazeSvc := maze.New(db)

	// setup handlers
	handlers.NewMaze(api, mazeSvc)

	log.Fatal(app.Listen(config.Router.Host))
}

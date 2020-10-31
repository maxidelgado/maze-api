package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxidelgado/maze-api/domain/game"
	"net/http"
)

func NewGames(router fiber.Router, svc game.Service) {
	h := gamesHandler{router: router, svc: svc}
	h.setupRoutes()
}

type gamesHandler struct {
	svc    game.Service
	router fiber.Router
}

func (h gamesHandler) setupRoutes() {
	m := h.router.Group("/games")
	{
		m.Post("", h.postGame)
		m.Get("", h.searchGames)
		m.Get("/:id", h.getGame)
		m.Delete("/:id", h.deleteGame)
		m.Put("/:id/move", h.putMove)
	}
}

/*
POST /api/v1/games :
	Starts a new game based on a given mazeId.
	Validates that the maze is able to be played.
	Returns all the required info to render the game on the client side.
*/
func (h gamesHandler) postGame(ctx *fiber.Ctx) error {
	var body struct {
		MazeId string `json:"maze_id"`
		Name   string `json:"name"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	newGame, err := h.svc.Start(ctx.Context(), body.MazeId, body.Name)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(newGame)
}

/*
GET /api/v1/games :
	Returns an existing game
*/
func (h gamesHandler) getGame(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := h.svc.Get(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

/*
DELETE /api/v1/games :
	Performs a hard delete to a given game
*/
func (h gamesHandler) deleteGame(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := h.svc.Delete(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

/*
PUT /api/v1/games/move :
	Performs a movement to a given spot (if valid).
*/
func (h gamesHandler) putMove(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var body struct {
		Spot string `json:"spot"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := h.svc.Move(ctx.Context(), id, body.Spot)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

/*
GET /api/v1/games?name=game_name
	Search matching games with a given name
*/
func (h gamesHandler) searchGames(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	if name == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "name param is required"})
	}

	response, err := h.svc.Query(ctx.Context(), name)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if len(response) == 0 {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

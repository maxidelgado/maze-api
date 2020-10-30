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
		m.Get("/:id", h.getGame)
		m.Delete("/:id", h.deleteGame)
		m.Put("/:id/move", h.putMove)
	}
}

func (h gamesHandler) postGame(ctx *fiber.Ctx) error {
	var body struct {
		MazeId string `json:"maze_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	newGame, err := h.svc.Start(ctx.Context(), body.MazeId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(newGame)
}

func (h gamesHandler) getGame(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := h.svc.Get(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

func (h gamesHandler) deleteGame(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := h.svc.Delete(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

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

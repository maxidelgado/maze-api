package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxidelgado/maze-api/domain/maze"
	"net/http"
)

func NewMaze(router fiber.Router, svc maze.Service) {
	h := mazeHandler{router: router, svc: svc}
	h.setupRoutes()
}

type mazeHandler struct {
	svc    maze.Service
	router fiber.Router
}

func (h mazeHandler) setupRoutes() {
	m := h.router.Group("/maze")
	{
		m.Get("/:id", h.getMaze)
		m.Post("", h.postMaze)
		m.Put("/spots", h.putSpots)
		m.Put("/paths", h.putPaths)
		m.Put("/quadrants", h.putQuadrants)
	}
}

func (h mazeHandler) postMaze(ctx *fiber.Ctx) error {
	var body struct {
		Coordinate maze.Coordinate `json:"coordinate"`
		Spots      []maze.Spot     `json:"spots"`
		Paths      []maze.Path     `json:"paths"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := h.svc.CreateMaze(ctx.Context(), body.Coordinate, body.Spots, body.Paths)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"maze_id": id})
}

func (h mazeHandler) getMaze(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id is required in path"})
	}

	m, err := h.svc.GetMaze(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(m)
}

func (h mazeHandler) putSpots(ctx *fiber.Ctx) error {
	var body struct {
		Id    string      `json:"id"`
		Spots []maze.Spot `json:"spots"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.svc.PutSpots(ctx.Context(), body.Id, body.Spots)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

func (h mazeHandler) putPaths(ctx *fiber.Ctx) error {
	var body struct {
		Id    string      `json:"id"`
		Paths []maze.Path `json:"paths"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.svc.PutPaths(ctx.Context(), body.Id, body.Paths)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

func (h mazeHandler) putQuadrants(ctx *fiber.Ctx) error {
	var body struct {
		Id         string          `json:"id"`
		Coordinate maze.Coordinate `json:"coordinate"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.svc.UpdateQuadrants(ctx.Context(), body.Id, body.Coordinate)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

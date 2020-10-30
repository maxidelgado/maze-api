package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/maxidelgado/maze-api/domain/maze"
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
	m := h.router.Group("/mazes")
	{
		m.Post("", h.postMaze)
		m.Get("/:id", h.getMaze)
		m.Put("/:id", h.putMaze)
		m.Delete("/:id", h.deleteMaze)

		m.Delete("/:id/spot", h.deleteSpot)
		m.Delete("/:id/path", h.deletePath)
	}
}

func (h mazeHandler) postMaze(ctx *fiber.Ctx) error {
	var body struct {
		Center maze.Coordinates `json:"center"`
		Spots  []maze.Spot      `json:"spots"`
		Paths  []maze.Path      `json:"paths"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := h.svc.Create(ctx.Context(), body.Center, body.Spots, body.Paths)
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

	m, err := h.svc.Get(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(m)
}

func (h mazeHandler) deleteSpot(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var coordinate maze.Coordinates
	if err := ctx.BodyParser(&coordinate); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.svc.DeleteSpot(ctx.Context(), id, coordinate)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

func (h mazeHandler) deletePath(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var path maze.Path
	if err := ctx.BodyParser(&path); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.svc.DeletePath(ctx.Context(), id, path)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

// UpdateMaze maze allows to:
//	- Add/replace existing spots
//	- Move quadrants by changing maze's center
//	- Add paths
func (h mazeHandler) putMaze(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var body struct {
		Center maze.Coordinates `json:"center"`
		Spots  []maze.Spot      `json:"spots"`
		Paths  []maze.Path      `json:"paths"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.svc.Update(ctx.Context(), id, body.Center, body.Spots, body.Paths)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

func (h mazeHandler) deleteMaze(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.svc.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

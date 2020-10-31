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

/*
POST /api/v1/mazes :
	Proceed to the creation of a new maze.
	Optionally the client can set some spots, paths, and can change the initial
	center of the maze (displaced quadrants).
*/
func (h mazeHandler) postMaze(ctx *fiber.Ctx) error {
	var body struct {
		Name   string           `json:"name"`
		Center maze.Coordinates `json:"center"`
		Spots  []maze.Spot      `json:"spots"`
		Paths  []maze.Path      `json:"paths"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := h.svc.Create(ctx.Context(), body.Name, body.Center, body.Spots, body.Paths)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"maze_id": id})
}

/*
GET /api/v1/mazes/{id} :
	Returns a given maze.
*/
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

/*
DELETE /api/v1/mazes/{id}/spot :
	Performs the deletion of a given spot from a maze.
	IMPORTANT: will produce a cascade deletion of all the related paths.
*/
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

/*
DELETE /api/v1/mazes/{id}/path :
	Performs the deletion of a given path from a maze.
	IMPORTANT: as all the paths have the corresponding edge/reverse-edge pair, both will be deleted.
*/
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

/*
PUT /api/v1/mazes/{id} :
	Perform a maze updating. Allows the following update:
		- Add/replace existing spots
		- Move quadrants by changing maze's center
		- Add paths
*/
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

/*
DELETE /api/v1/mazes/{id} :
	Performs a hard delete for a given maze.
*/
func (h mazeHandler) deleteMaze(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.svc.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

/*
GET /api/v1/mazes?name=game_name
	Search matching mazes with a given name
*/
func (h mazeHandler) searchMazes(ctx *fiber.Ctx) error {
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

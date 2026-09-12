package server

import (
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
)

// GameHandler exposes the game run lifecycle over HTTP. The editor asks for a
// game; the orchestrator decides what processes that means
type GameHandler struct {
	controller GameController
}

func NewGameHandler(controller GameController) *GameHandler {
	return &GameHandler{
		controller: controller,
	}
}

func (h *GameHandler) HandleRunGame(ctx *fiber.Ctx) error {
	var req entity.RunGameRequest

	if err := ctx.BodyParser(&req); err != nil {
		ls.ErrorWith(err).Msg("failed to parse request body")
		return ctx.Status(http.StatusBadRequest).JSON(entity.NewGameActionResponse(false, "failed to parse body"))
	}

	if req.Filepath == "" {
		ls.Error().Msg("no game file provided to run")
		return ctx.Status(http.StatusBadRequest).JSON(entity.NewGameActionResponse(false, "a game file is required to run a game"))
	}

	// Matches what the builder accepts, so an unusable path is rejected
	// before a process is spawned for it
	if filepath.Ext(req.Filepath) != ".json" {
		ls.Error().Msg("the filepath does not point to a json file")
		return ctx.Status(http.StatusBadRequest).JSON(entity.NewGameActionResponse(false, "games can only be run from json files"))
	}

	if err := h.controller.RunGame(req.Filepath); err != nil {
		ls.ErrorWith(err).Msg("failed to run the game")
		return ctx.Status(http.StatusInternalServerError).JSON(entity.NewGameActionResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(entity.NewGameActionResponse(true, "game started successfully"))
}

func (h *GameHandler) HandleStopGame(ctx *fiber.Ctx) error {
	h.controller.StopGame()

	ls.Info().Msg("game run stopped")

	return ctx.Status(http.StatusOK).JSON(entity.NewGameActionResponse(true, "game stopped successfully"))
}

func (h *GameHandler) HandleGameStatus(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(entity.GameStatusResponse{
		Success: true,
		Running: h.controller.IsGameRunning(),
	})
}

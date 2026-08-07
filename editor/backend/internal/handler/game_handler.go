package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/entity"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/logger"
)

var ls = logger.Get()

func GetGameState(ctx *fiber.Ctx) error {
	gsm := gamestatemanager.Get()

	resp := entity.GetGameStateResponse{
		Success:   true,
		GameState: gsm.GetGameState(),
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func HandleSaveGame(ctx *fiber.Ctx) error {
	var req entity.SaveGameRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		ls.ErrorWith(err).Msg("Failed to parse request body")
		return ctx.Status(http.StatusBadRequest).JSON(entity.NewErrorResponse("failed to parse body"))
	}
	// Enforce .json extension
	if filepath.Ext(req.Filepath) != ".json" {
		ls.Error().Msg("The filepath does not point to json file")
		return ctx.Status(http.StatusBadRequest).JSON(entity.NewErrorResponse("game state can only be stored in json files"))
	}

	err = saveGame(gamestatemanager.Get(), req.Filepath)
	if err != nil {
		ls.ErrorWith(err).Msg("Failed to save game")
		return ctx.Status(http.StatusInternalServerError).JSON(entity.NewErrorResponse("internal server error"))
	}

	resp := entity.SaveGameResponse{
		Success: true,
		Message: "game state saved successfully",
	}

	ls.Info().Msgf("game state saved successfully to file: %s", req.Filepath)

	return ctx.Status(http.StatusOK).JSON(resp)
}

func saveGame(gsm *gamestatemanager.GameStateManager, filepath string) error {
	data := gsm.GetGameState()

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

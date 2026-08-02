package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
)

// Request that triggers the save game functionality
type saveGameRequest struct {
	Filepath string `json:"filepath"` // destination where the game state should be saved. (JSON file)
}

func HandleSaveGame(ctx *fiber.Ctx) error {
	var req saveGameRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse body",
		})
	}
	// Enforce .json extension
	if filepath.Ext(req.Filepath) != ".json" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "game state can only be stored in json files",
		})
	}

	err = saveGame(gamestatemanager.Get(), req.Filepath)
	if err != nil {
		ls.ErrorWith(err).Msg("Failed to save game")
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "game state saved successfully",
	})
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

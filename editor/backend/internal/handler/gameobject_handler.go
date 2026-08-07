package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/entity"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
)

func AddGameobject(ctx *fiber.Ctx) error {
	gameobject := gameobject.NewGameobject()

	gsm := gamestatemanager.Get()
	gsm.AddGameobject(gameobject)

	resp := entity.CreateGameobjectResponse{
		Success:       true,
		ObjectDetails: gameobject.GetGameobjectDetails(),
	}

	ls.Info().Msg("gameobject added successfully")

	return ctx.Status(http.StatusOK).JSON(resp)
}

func DeleteGameobject(ctx *fiber.Ctx) error {
	gameobjectID := ctx.Params("objectID")

	gsm := gamestatemanager.Get()

	_, found := gsm.GetGameobject(gameobjectID)
	if !found {
		ls.Error().Msgf("gameobject with id: %s not found", gameobjectID)
		return ctx.Status(http.StatusNotFound).JSON(entity.NewErrorResponse("gameobject not found"))
	}

	gsm.DeleteGameobject(gameobjectID)

	resp := entity.DeleteGameobjectResponse{
		Success: true,
		Message: "deletion request successfully executed",
	}

	ls.Info().Msgf("gameobject with id: %s deleted successfully", gameobjectID)

	return ctx.Status(http.StatusOK).JSON(resp)
}

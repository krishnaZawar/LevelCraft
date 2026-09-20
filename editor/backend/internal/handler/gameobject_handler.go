package handler

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/entity"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
)

// Matches the "(2)" a previous duplicate appended, so copying a copy gives
// "Player (3)" rather than "Player (2) (1)".
var copyIndexSuffix = regexp.MustCompile(`\s*\(\d+\)$`)

// Names a copy of sourceName so it doesn't collide with any existing object.
func duplicateName(sourceName string, taken map[string]struct{}) string {
	if _, exists := taken[sourceName]; !exists {
		return sourceName
	}
	base := copyIndexSuffix.ReplaceAllString(sourceName, "")
	if base == "" {
		return sourceName
	}
	for index := 1; ; index++ {
		candidate := fmt.Sprintf("%s (%d)", base, index)
		if _, exists := taken[candidate]; !exists {
			return candidate
		}
	}
}

// collects the names currently in use across the scene
func takenNames(scene map[string]interface{}) map[string]struct{} {
	taken := map[string]struct{}{}
	for _, objData := range scene {
		data, ok := objData.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := data[gameobject.Gameobject_CurLabelName].(string); ok {
			taken[name] = struct{}{}
		}
	}
	return taken
}

func DuplicateGameobject(ctx *fiber.Ctx) error {
	gameobjectID := ctx.Params("objectID")

	gsm := gamestatemanager.Get()

	source, found := gsm.GetGameobject(gameobjectID)
	if !found {
		ls.Error().Msgf("gameobject with id: %s not found", gameobjectID)
		return ctx.Status(http.StatusNotFound).JSON(entity.NewErrorResponse("gameobject not found"))
	}

	// Rebuilt from the source's own details rather than sharing anything with
	// it, so the copy owns its components and editing one can't affect the other.
	clone := gameobject.NewGameobject()
	details := source.GetGameobjectDetails()
	details[gameobject.Gameobject_CurLabelID] = clone.GetID()
	details[gameobject.Gameobject_CurLabelName] = duplicateName(
		source.GetName(),
		takenNames(gsm.GetGameState()),
	)

	if err := clone.BuildFromDetails(details); err != nil {
		ls.ErrorWith(err).Msgf("failed to duplicate gameobject with id: %s", gameobjectID)
		return ctx.Status(http.StatusInternalServerError).JSON(entity.NewErrorResponse("failed to duplicate the gameobject"))
	}

	gsm.AddGameobject(clone)

	resp := entity.DuplicateGameobjectResponse{
		Success:       true,
		ObjectDetails: clone.GetGameobjectDetails(),
	}

	ls.Info().Msgf("gameobject with id: %s duplicated into %s", gameobjectID, clone.GetID())

	return ctx.Status(http.StatusOK).JSON(resp)
}

func AddGameobject(ctx *fiber.Ctx) error {
	gameobject := gameobject.NewGameobject()
	gsm := gamestatemanager.Get()

	gameobject.SetName(duplicateName(
		gameobject.GetName(),
		takenNames(gsm.GetGameState()),
	))

	gsm.AddGameobject(gameobject)

	resp := entity.CreateGameobjectResponse{
		Success:       true,
		ObjectDetails: gameobject.GetGameobjectDetails(),
	}

	ls.Info().Msg("gameobject added successfully")

	return ctx.Status(http.StatusOK).JSON(resp)
}

func UpdateGameobject(ctx *fiber.Ctx) error {
	gameobjectID := ctx.Params("objectID")

	var req entity.UpdateGameobjectRequest
	if err := ctx.BodyParser(&req); err != nil {
		ls.ErrorWith(err).Msg("failed to parse request body")
		return ctx.Status(http.StatusBadRequest).JSON(entity.NewErrorResponse("failed to parse body"))
	}

	gsm := gamestatemanager.Get()

	obj, found := gsm.GetGameobject(gameobjectID)
	if !found {
		ls.Error().Msgf("gameobject with id: %s not found", gameobjectID)
		return ctx.Status(http.StatusNotFound).JSON(entity.NewErrorResponse("gameobject not found"))
	}

	if req.Name != nil {
		obj.SetName(*req.Name)
	}
	if req.Group != nil {
		obj.SetGroup(*req.Group)
	}

	resp := entity.UpdateGameobjectResponse{
		Success:       true,
		ObjectDetails: obj.GetGameobjectDetails(),
	}

	ls.Info().Msgf("gameobject with id: %s updated successfully", gameobjectID)

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

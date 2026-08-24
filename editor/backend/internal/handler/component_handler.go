package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/entity"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
)

func GetComponents(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"components": component.ComponentList,
	})
}

func AddComponent(ctx *fiber.Ctx) error {
	gameobjectID := ctx.Params("objectID")
	componentName := ctx.Params("componentName")

	gsm := gamestatemanager.Get()
	obj, found := gsm.GetGameobject(gameobjectID)
	if !found {
		ls.Error().Msgf("gameobject with id: %s not found", gameobjectID)
		return ctx.Status(http.StatusNotFound).JSON(entity.NewErrorResponse("gameobject not found"))
	}

	comp, found := component.NewComponentRegistry().GetComponent(componentName)
	if !found {
		ls.Error().Msgf("component with name: %s does not exist", componentName)
		return ctx.Status(http.StatusNotFound).JSON(entity.NewErrorResponse("component not found"))
	}

	added := obj.AddComponent(comp)
	if !added {
		ls.Warn().Msgf("%s component already attached to gameobject with id: %s", componentName, gameobjectID)
		return ctx.Status(http.StatusBadRequest).JSON(entity.NewErrorResponse("component already attached to the gameobject"))
	}

	resp := entity.ComponentResponse{
		Success:       true,
		ObjectDetails: obj.GetGameobjectDetails(),
	}

	ls.Info().Msgf("%s component attached to gameobject with id: %s successfully", componentName, gameobjectID)

	return ctx.Status(http.StatusOK).JSON(resp)
}

func DeleteComponent(ctx *fiber.Ctx) error {
	gameobjectID := ctx.Params("objectID")
	componentName := ctx.Params("componentName")

	gsm := gamestatemanager.Get()

	obj, found := gsm.GetGameobject(gameobjectID)
	if !found {
		ls.Error().Msgf("gameobject with id: %s not found", gameobjectID)
		return ctx.Status(http.StatusNotFound).JSON(entity.NewErrorResponse("gameobject not found"))
	}

	_, found = obj.GetComponent(componentName)
	if !found {
		ls.Error().Msgf("%s component is not attached to gameobject with id: %s", componentName, gameobjectID)
		return ctx.Status(http.StatusNotFound).JSON(entity.NewErrorResponse("component not found"))
	}

	obj.RemoveComponent(componentName)

	resp := entity.ComponentResponse{
		Success:       true,
		ObjectDetails: obj.GetGameobjectDetails(),
	}

	ls.Info().Msgf("%s component removed from gameobject with id: %s successfully", componentName, gameobjectID)

	return ctx.Status(http.StatusOK).JSON(resp)
}

func UpdateComponent(ctx *fiber.Ctx) error {
	gameobjectID := ctx.Params("objectID")
	componentName := ctx.Params("componentName")

	var req entity.UpdateComponentRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		ls.ErrorWith(err).Msg("failed to parse request body")
		return ctx.Status(http.StatusBadRequest).BodyParser(entity.NewErrorResponse("failed to parse body"))
	}

	gsm := gamestatemanager.Get()

	obj, found := gsm.GetGameobject(gameobjectID)
	if !found {
		ls.Error().Msgf("gameobject with id: %s not found", gameobjectID)
		return ctx.Status(http.StatusNotFound).JSON(entity.NewErrorResponse("gameobject not found"))
	}

	comp, found := obj.GetComponent(componentName)
	if !found {
		ls.Error().Msgf("%s component is not attached to gameobject with id: %s", componentName, gameobjectID)
		return ctx.Status(http.StatusNotFound).JSON(entity.NewErrorResponse("component not found"))
	}

	err = comp.BuildFromDetails(req.Details)
	if err != nil {
		ls.ErrorWith(err).Msgf("failed to build %s component from details", componentName)
		return ctx.Status(http.StatusBadRequest).JSON(entity.NewErrorResponse(err.Error()))
	}

	resp := entity.ComponentResponse{
		Success:       true,
		ObjectDetails: obj.GetGameobjectDetails(),
	}

	ls.Info().Msgf("%s component updated for gameobject with id: %s successfully", componentName, gameobjectID)

	return ctx.Status(http.StatusOK).JSON(resp)
}

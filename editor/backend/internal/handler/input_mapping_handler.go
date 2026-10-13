package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/entity"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/inputmapping"
)

func GetInputMapping(ctx *fiber.Ctx) error {
	mapping := inputmapping.Get()
	resp := entity.GetInputMappingResponse{
		Success: true,
		Mapping: *mapping,
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

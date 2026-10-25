package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/utils/input"
)

type GetInputMappingResponse struct {
	Success bool               `json:"success"`
	Mapping input.InputMapping `json:"mapping"`
}

func GetInputMapping(ctx *fiber.Ctx) error {
	mapping := inputmanager.GetMapping()
	resp := GetInputMappingResponse{
		Success: true,
		Mapping: *mapping,
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CreateRoutes(app *fiber.App, controller GameController) {
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"data": "pong.",
		})
	})

	handler := NewGameHandler(controller)

	game := app.Group("/game")
	game.Post("/run", handler.HandleRunGame)
	game.Post("/stop", handler.HandleStopGame)
	game.Get("/status", handler.HandleGameStatus)
}

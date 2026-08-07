package router

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/handler"
)

func CreateRoutes(app *fiber.App) {
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"data": "pong.",
		})
	})

	game := app.Group("/game")
	game.Post("/save", handler.HandleSaveGame)
	game.Get("/state", handler.GetGameState)

	gameobjects := app.Group("/gameobjects")
	gameobjects.Post("/", handler.AddGameobject)
	gameobjects.Delete("/:objectID", handler.DeleteGameobject)

	components := gameobjects.Group("/:objectID/components")
	components.Post("/:componentName", handler.AddComponent)
	components.Delete("/:componentName", handler.DeleteComponent)
	components.Put("/:componentName", handler.UpdateComponent)
}

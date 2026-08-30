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
	game.Post("/load", handler.HandleLoadGame)
	game.Get("/state", handler.GetGameState)

	components := app.Group("/components")
	components.Get("/", handler.GetComponents)

	gameobjects := app.Group("/gameobjects")
	gameobjects.Post("/", handler.AddGameobject)
	gameobjects.Delete("/:objectID", handler.DeleteGameobject)

	objComponents := gameobjects.Group("/:objectID/components")
	objComponents.Post("/:componentName", handler.AddComponent)
	objComponents.Delete("/:componentName", handler.DeleteComponent)
	objComponents.Put("/:componentName", handler.UpdateComponent)
}

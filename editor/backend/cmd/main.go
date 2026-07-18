package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/handler"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/router"
)

func main() {
	app := fiber.New()

	// run 2 threads
	// - 1. game loop: process command requests and update game state
	// - 2. main thread: accept comman requests and send out responses back to the frontend
	go handler.EditorLoop()

	router.CreateRoutes(app)

	app.Listen(":3000")
}

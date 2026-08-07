package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/router"
)

func main() {
	app := fiber.New()

	router.CreateRoutes(app)

	app.Listen(":3000")
}

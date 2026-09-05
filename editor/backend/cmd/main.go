package main

import (
	"flag"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/router"
)

func main() {
	app := fiber.New()

	port := flag.Int("port", 3000, "the value of the port on which the service starts")

	flag.Parse()

	// The only consumer of this API is our own Electron client running
	// locally (dev server origin or a file:// packaged build), never a
	// public caller, so a permissive CORS policy is appropriate here.
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	router.CreateRoutes(app)

	err := app.Listen(":" + strconv.Itoa(*port))
	if err != nil {
		panic(err)
	}
}

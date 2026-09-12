package server

import (
	"net"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/utils/logger"
)

var ls = logger.New(base.ServiceName)

// GameController is the slice of the orchestrator the control API needs.
// Declared here so the dependency only points one way
type GameController interface {
	// starts the builder processes for the scene at the given path
	RunGame(scenePath string) error

	// stops every process belonging to the current game run
	StopGame()

	// reports whether a game run is currently live
	IsGameRunning() bool
}

// Start brings up the control API on an OS assigned port and returns its base
// URL. Assigned rather than guessed: this address is handed to the editor
func Start(controller GameController) (string, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// The only caller is the local editor client, never a public one, so this
	// matches the permissive policy the editor backend already uses
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Set("Access-Control-Allow-Origin", "*")
		ctx.Set("Access-Control-Allow-Headers", "Content-Type")
		ctx.Set("Access-Control-Allow-Methods", "GET,POST")
		if ctx.Method() == fiber.MethodOptions {
			return ctx.SendStatus(fiber.StatusNoContent)
		}
		return ctx.Next()
	})

	CreateRoutes(app, controller)

	port := listener.Addr().(*net.TCPAddr).Port

	go func() {
		if err := app.Listener(listener); err != nil {
			ls.ErrorWith(err).Msg("orchestrator control API stopped")
		}
	}()

	return "http://localhost:" + strconv.Itoa(port), nil
}

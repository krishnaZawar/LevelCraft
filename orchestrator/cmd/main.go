package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/orchestrator"
)

func main() {
	orch := orchestrator.NewOrchestrator()
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	err := orch.Run(ctx)
	if err != nil {
		panic(err)
	}
}

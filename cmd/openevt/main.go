package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/brandon1024/cmder"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	err := cmder.Execute(ctx, cmd, cmder.WithEnvironmentBinding())
	cancel()

	if err != nil {
		slog.Error("error caught - shutting down", "err", err)
		os.Exit(1)
	}
}

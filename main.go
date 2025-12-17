package main

import (
	"os"
	"scope-guardian/logger"

	"golang.org/x/exp/slog"
)

func main() {
	logger.SetGlobalLogger(
		logger.NewSlogLogger(
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))))

	logger.Info("Tool started")
}

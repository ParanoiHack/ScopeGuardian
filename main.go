package main

import (
	"os"
	"scope-guardian/display"
	"scope-guardian/engine"
	"scope-guardian/loader"
	"scope-guardian/logger"
	"scope-guardian/parser"

	"golang.org/x/exp/slog"
)

const (
	logInfoLoadConfigFile = "Loading configuration file"
)

func main() {
	logger.SetGlobalLogger(
		logger.NewSlogLogger(
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))))

	display.DisplayBanner()
	display.DisplayCredit()

	logger.Info(logInfoLoadConfigFile)

	args, err := parser.Parse(os.Args[1:])
	if err != nil {
		logger.Panic(err.Error())
	}

	config, err := loader.Load(args.Config)
	if err != nil {
		logger.Panic(err.Error())
	}

	engine := engine.NewEngine()

	engine.Initialize(config)
	engine.Start()

	findings := engine.LoadFindings()

	engine.SyncResults(13, "final-test")

	display.DisplayFindings(findings)
}

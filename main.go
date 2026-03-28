package main

import (
	"os"
	"scope-guardian/display"
	"scope-guardian/engine"
	"scope-guardian/loader"
	"scope-guardian/logger"
	"scope-guardian/parser"
	securitygate "scope-guardian/features/security-gate"

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
		logger.Error(err.Error())
		parser.PrintUsage(os.Stdout)
		os.Exit(1)
	}

	config, err := loader.Load(args.Config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	engine := engine.NewEngine()

	engine.Initialize(config)
	engine.Start()

	findings := engine.LoadFindings()

	if args.Sync {
		engine.SyncResults(args.ProjectName, args.Branch, config.ProtectedBranches)
	}

	display.DisplayFindings(findings)

	if len(args.Thresholds) > 0 {
		if !securitygate.Evaluate(findings, args.Thresholds) {
			os.Exit(-1)
		}
	}
}

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
	logInfoLoadConfigFile  = "Loading configuration file"
	logErrOutputFile       = "Failed to create output log file"
	logErrCloseOutputFile  = "Failed to close output log file"
)

func main() {
	logger.SetGlobalLogger(
		logger.NewSlogLogger(
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))))

	display.DisplayBanner()
	display.DisplayCredit()

	args, err := parser.Parse(os.Args[1:])
	if err != nil {
		logger.Error(err.Error())
		parser.PrintUsage(os.Stdout)
		os.Exit(1)
	}

	if args.Quiet {
		logger.SetGlobalLogger(logger.NewNullLogger())
	} else if args.Output != "" {
		f, err := os.Create(args.Output)
		if err != nil {
			logger.Error(logErrOutputFile, logger.Err(err))
			os.Exit(1)
		}
		defer func() {
			if cerr := f.Close(); cerr != nil {
				logger.Error(logErrCloseOutputFile, logger.Err(cerr))
			}
		}()
		logger.SetGlobalLogger(
			logger.NewSlogLogger(
				slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{}))))
	}

	logger.Info(logInfoLoadConfigFile)

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
		findingsToEvaluate := findings
		if args.Sync {
			remoteFindings, err := engine.GetDefectDojoFindings(args.ProjectName, args.Branch, config.ProtectedBranches)
			if err != nil {
				logger.Error(err.Error())
			} else {
				findingsToEvaluate = remoteFindings
			}
		}
		if !securitygate.Evaluate(findingsToEvaluate, args.Thresholds) {
			os.Exit(-1)
		}
	}
}

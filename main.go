package main

import (
	"io"
	"os"
	"ScopeGuardian/display"
	"ScopeGuardian/engine"
	"ScopeGuardian/loader"
	"ScopeGuardian/logger"
	"ScopeGuardian/parser"
	securitygate "ScopeGuardian/features/security-gate"

	"golang.org/x/exp/slog"
)

const (
	logInfoLoadConfigFile = "Loading configuration file"
	logErrOutputFile      = "Failed to create output log file"
	logErrCloseOutputFile = "Failed to close output log file"
)

func main() {
	logger.SetGlobalLogger(
		logger.NewSlogLogger(
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))))

	args, err := parser.Parse(os.Args[1:])
	if err != nil {
		logger.Error(err.Error())
		parser.PrintUsage(os.Stdout)
		os.Exit(1)
	}

	// displayOut is the writer used for banner, credit, and findings output.
	// logOut is the writer used for structured log messages.
	// When -o is set both writers tee to stdout AND the file.
	displayOut := io.Writer(os.Stdout)

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
		tee := io.MultiWriter(os.Stdout, f)
		displayOut = tee
		logger.SetGlobalLogger(
			logger.NewSlogLogger(
				slog.New(slog.NewTextHandler(tee, &slog.HandlerOptions{}))))
	}

	display.DisplayBanner(displayOut)
	display.DisplayCredit(displayOut)

	logger.Info(logInfoLoadConfigFile)

	config, err := loader.Load(args.Config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	eng := engine.NewEngine()

	eng.Initialize(config)
	eng.Start()

	findings := eng.LoadFindings()

	if args.Sync {
		eng.SyncResults(args.ProjectName, args.Branch, config.ProtectedBranches)
	}

	display.DisplayFindings(displayOut, findings)

	if len(args.Thresholds) > 0 {
		findingsToEvaluate := findings
		if args.Sync {
			remoteFindings, err := eng.GetDefectDojoFindings(args.ProjectName, args.Branch, config.ProtectedBranches)
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

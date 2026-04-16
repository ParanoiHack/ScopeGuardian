package main

import (
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
	logInfoDumpFindings   = "Findings successfully written to output file"
	logErrOutputFile      = "Failed to create output file"
	logErrCloseOutputFile = "Failed to close output file"
	logErrDumpFindings    = "Failed to write findings to output file"
	logErrFilterByDD      = "Failed to filter findings against DefectDojo; displaying all local findings instead"
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

	if args.Quiet {
		logger.SetGlobalLogger(logger.NewNullLogger())
	}

	display.DisplayBanner(os.Stdout)
	display.DisplayCredit(os.Stdout)

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
		filtered, err := eng.FilterFindingsByDD(findings, args.ProjectName, args.Branch, config.ProtectedBranches)
		if err != nil {
			logger.Error(logErrFilterByDD, logger.Err(err))
		} else {
			findings = filtered
		}
	}

	display.DisplayFindings(os.Stdout, findings)

	if args.Output != "" {
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
		if err := display.DumpFindings(f, findings, args.Format); err != nil {
			logger.Error(logErrDumpFindings, logger.Err(err))
			os.Exit(1)
		}
		logger.Info(logInfoDumpFindings, logger.Any("file", args.Output))
	}

	if len(args.Thresholds) > 0 {
		if !securitygate.Evaluate(findings, args.Thresholds) {
			os.Exit(-1)
		}
	}
}

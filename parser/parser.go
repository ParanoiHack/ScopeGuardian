package parser

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var validSeverities = []string{severityCritical, severityHigh, severityMedium, severityLow, severityInfo}

func Parse(args []string) (Args, error) {
	fs := flag.NewFlagSet("scope-guardian", flag.ContinueOnError)

	sync        := fs.Bool("sync", false, "Enable sync result with DefectDojo")
	threshold   := fs.String("threshold", "", "Enable security gate (e.g., critical=1)")
	projectName := fs.String("projectName", "", "Name of the project to scan")
	branch      := fs.String("branch", "", "Project branch to scan")

	if err := fs.Parse(args); err != nil {
		return Args{}, err
	}

	remaining := fs.Args()
	if len(remaining) < 1 {
		return Args{}, errors.New(errConfigRequired)
	}

	config := remaining[0]

	if *projectName == "" {
		return Args{}, errors.New(errProjectNameRequired)
	}

	if *branch == "" {
		return Args{}, errors.New(errBranchRequired)
	}

	var parsedThreshold *Threshold
	if *threshold != "" {
		t, err := parseThreshold(*threshold)
		if err != nil {
			return Args{}, err
		}
		parsedThreshold = t
	}

	return Args{
		Config:      config,
		ProjectName: *projectName,
		Branch:      *branch,
		Sync:        *sync,
		Threshold:   parsedThreshold,
	}, nil
}

func parseThreshold(s string) (*Threshold, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return nil, errors.New(errInvalidThreshold)
	}

	severity := strings.ToLower(strings.TrimSpace(parts[0]))
	valueStr := strings.TrimSpace(parts[1])

	if !isValidSeverity(severity) {
		return nil, fmt.Errorf(errInvalidSeverity, severity)
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil || value < 0 {
		return nil, fmt.Errorf(errInvalidThresholdValue, valueStr)
	}

	return &Threshold{Severity: severity, Value: value}, nil
}

func isValidSeverity(severity string) bool {
	for _, s := range validSeverities {
		if s == severity {
			return true
		}
	}
	return false
}

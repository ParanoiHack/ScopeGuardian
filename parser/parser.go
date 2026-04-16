package parser

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var validSeverities = []string{severityCritical, severityHigh, severityMedium, severityLow, severityInfo}

// PrintUsage writes the CLI usage help to w.
func PrintUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: ScopeGuardian [flags] <config-file>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Arguments:")
	fmt.Fprintln(w, "  <config-file>  Path to the TOML configuration file (required)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --projectName string   Name of the project to scan (required)")
	fmt.Fprintln(w, "  --branch string        Project branch to scan (required)")
	fmt.Fprintln(w, "  --sync                 Enable sync result with DefectDojo (default: false)")
	fmt.Fprintln(w, "  --threshold string     Enable security gate, e.g. critical=1 or critical=1,high=2 (optional)")
	fmt.Fprintln(w, "                         Supported severities: critical, high, medium, low, info")
	fmt.Fprintln(w, "                         Multiple thresholds can be comma-separated (e.g. critical=1,high=2)")
	fmt.Fprintln(w, "  --filter string        Comma-separated finding statuses to include in display and output (optional)")
	fmt.Fprintln(w, "                         Supported values: ACTIVE, INACTIVE, DUPLICATE (default: ACTIVE)")
	fmt.Fprintln(w, "                         Example: --filter ACTIVE,DUPLICATE")
	fmt.Fprintln(w, "  -q                     Quiet mode: suppress all log output (default: false)")
	fmt.Fprintln(w, "  -o string              Write findings to the specified file (optional)")
	fmt.Fprintln(w, "  --format string        Output format for -o: json, csv, or raw (default: json)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Example:")
	fmt.Fprintln(w, "  ScopeGuardian --projectName my-service --branch main ./config.toml")
	fmt.Fprintln(w, "  ScopeGuardian --projectName my-service --branch main --threshold critical=1,high=2 --sync ./config.toml")
	fmt.Fprintln(w, "  ScopeGuardian --projectName my-service --branch main -q ./config.toml")
	fmt.Fprintln(w, "  ScopeGuardian --projectName my-service --branch main -o /tmp/scan.json --format json ./config.toml")
}

// Parse parses the CLI arguments in args and returns a validated Args struct.
// It returns an error if required flags (--projectName, --branch) are missing,
// the config-file positional argument is absent, or the threshold flag is malformed.
func Parse(args []string) (Args, error) {
	fs := flag.NewFlagSet("ScopeGuardian", flag.ContinueOnError)

	sync        := fs.Bool("sync", false, "Enable sync result with DefectDojo")
	threshold   := fs.String("threshold", "", "Enable security gate (e.g., critical=1,high=2)")
	projectName := fs.String("projectName", "", "Name of the project to scan")
	branch      := fs.String("branch", "", "Project branch to scan")
	quiet       := fs.Bool("q", false, "Quiet mode: suppress all log output")
	output      := fs.String("o", "", "Write findings to the specified file")
	format      := fs.String("format", FormatJSON, "Output format for -o: json, csv, or raw")
	filter      := fs.String("filter", "", "Comma-separated finding statuses to display: ACTIVE, INACTIVE, DUPLICATE (default: ACTIVE)")

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

	var parsedThresholds []Threshold
	if *threshold != "" {
		ts, err := parseThresholds(*threshold)
		if err != nil {
			return Args{}, err
		}
		parsedThresholds = ts
	}

	if !isValidFormat(*format) {
		return Args{}, fmt.Errorf(errInvalidFormat, *format)
	}

	statusFilters, err := parseStatusFilters(*filter)
	if err != nil {
		return Args{}, err
	}

	return Args{
		Config:        config,
		ProjectName:   *projectName,
		Branch:        *branch,
		Sync:          *sync,
		Quiet:         *quiet,
		Output:        *output,
		Format:        *format,
		Thresholds:    parsedThresholds,
		StatusFilters: statusFilters,
	}, nil
}

// parseThresholds parses a comma-separated list of threshold strings, each of the
// form "severity=value" (e.g. "critical=1,high=2"), and returns a slice of
// Threshold values. Returns an error if any token is malformed.
func parseThresholds(s string) ([]Threshold, error) {
	tokens := strings.Split(s, ",")
	thresholds := make([]Threshold, 0, len(tokens))
	for _, token := range tokens {
		t, err := parseThreshold(strings.TrimSpace(token))
		if err != nil {
			return nil, err
		}
		thresholds = append(thresholds, *t)
	}
	return thresholds, nil
}

// parseThreshold parses a threshold string of the form "severity=value"
// (e.g. "critical=1") and returns a Threshold. Returns an error if the
// format is invalid, the severity is unrecognised, or the value is negative.
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

// isValidSeverity reports whether severity is one of the recognised severity levels.
func isValidSeverity(severity string) bool {
	for _, s := range validSeverities {
		if s == severity {
			return true
		}
	}
	return false
}

// isValidFilterStatus reports whether status (case-insensitive) is one of the
// recognised finding status values: ACTIVE, INACTIVE, DUPLICATE.
func isValidFilterStatus(status string) bool {
	upper := strings.ToUpper(status)
	for _, s := range validFilterStatuses {
		if s == upper {
			return true
		}
	}
	return false
}

// parseStatusFilters parses a comma-separated list of finding status values
// (e.g. "ACTIVE,DUPLICATE") and returns a normalised (upper-case) slice.
// When s is empty the default filter (ACTIVE only) is returned.
// Returns an error if any token is not a recognised status.
func parseStatusFilters(s string) ([]string, error) {
	if s == "" {
		return defaultFilterStatuses, nil
	}
	tokens := strings.Split(s, ",")
	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		trimmed := strings.TrimSpace(token)
		if !isValidFilterStatus(trimmed) {
			return nil, fmt.Errorf(errInvalidFilter, trimmed)
		}
		result = append(result, strings.ToUpper(trimmed))
	}
	return result, nil
}

// isValidFormat reports whether format is one of the recognised output formats.
func isValidFormat(format string) bool {
	for _, f := range validFormats {
		if f == format {
			return true
		}
	}
	return false
}

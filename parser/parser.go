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
	fmt.Fprintln(w, "Usage: scope-guardian [flags] <config-file>")
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
	fmt.Fprintln(w, "  -q                     Quiet mode: suppress all log output (default: false)")
	fmt.Fprintln(w, "  -o string              Write output to the specified file in addition to stdout (optional)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Example:")
	fmt.Fprintln(w, "  scope-guardian --projectName my-service --branch main ./config.toml")
	fmt.Fprintln(w, "  scope-guardian --projectName my-service --branch main --threshold critical=1,high=2 --sync ./config.toml")
	fmt.Fprintln(w, "  scope-guardian --projectName my-service --branch main -q ./config.toml")
	fmt.Fprintln(w, "  scope-guardian --projectName my-service --branch main -o /tmp/scan.log ./config.toml")
}

// Parse parses the CLI arguments in args and returns a validated Args struct.
// It returns an error if required flags (--projectName, --branch) are missing,
// the config-file positional argument is absent, or the threshold flag is malformed.
func Parse(args []string) (Args, error) {
	fs := flag.NewFlagSet("scope-guardian", flag.ContinueOnError)

	sync        := fs.Bool("sync", false, "Enable sync result with DefectDojo")
	threshold   := fs.String("threshold", "", "Enable security gate (e.g., critical=1,high=2)")
	projectName := fs.String("projectName", "", "Name of the project to scan")
	branch      := fs.String("branch", "", "Project branch to scan")
	quiet       := fs.Bool("q", false, "Quiet mode: suppress all log output")
	output      := fs.String("o", "", "Write output to the specified file in addition to stdout")

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

	return Args{
		Config:      config,
		ProjectName: *projectName,
		Branch:      *branch,
		Sync:        *sync,
		Quiet:       *quiet,
		Output:      *output,
		Thresholds:  parsedThresholds,
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

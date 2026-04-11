package parser

const (
	errConfigRequired        = "config filepath is required"
	errProjectNameRequired   = "projectName is required"
	errBranchRequired        = "branch is required"
	errInvalidThreshold      = "invalid threshold format, expected severity=value (e.g., critical=1)"
	errInvalidSeverity       = "invalid severity: %s, must be one of: critical, high, medium, low, info"
	errInvalidThresholdValue = "invalid threshold value: %s, must be a non-negative integer"
	errInvalidFormat         = "invalid format: %s, must be one of: json, csv, raw"
)

const (
	severityCritical = "critical"
	severityHigh     = "high"
	severityMedium   = "medium"
	severityLow      = "low"
	severityInfo     = "info"
)

const (
	FormatJSON = "json"
	FormatCSV  = "csv"
	FormatRaw  = "raw"
)

var validFormats = []string{FormatJSON, FormatCSV, FormatRaw}

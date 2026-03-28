package parser

const (
	errConfigRequired        = "config filepath is required"
	errProjectNameRequired   = "projectName is required"
	errBranchRequired        = "branch is required"
	errInvalidThreshold      = "invalid threshold format, expected severity=value (e.g., critical=1)"
	errInvalidSeverity       = "invalid severity: %s, must be one of: critical, high, medium, low, info"
	errInvalidThresholdValue = "invalid threshold value: %s, must be a non-negative integer"
)

const (
	severityCritical = "critical"
	severityHigh     = "high"
	severityMedium   = "medium"
	severityLow      = "low"
	severityInfo     = "info"
)

package securitygate

const (
	severityCritical = "CRITICAL"
	severityHigh     = "HIGH"
	severityMedium   = "MEDIUM"
	severityLow      = "LOW"
	severityInfo     = "INFO"
)

// severityRank maps a normalised (upper-case) severity label to a numeric rank
// used for "at-or-above" comparisons. Higher rank means more severe.
var severityRank = map[string]int{
	severityCritical: 5,
	severityHigh:     4,
	severityMedium:   3,
	severityLow:      2,
	severityInfo:     1,
}

const (
	logInfoGatePass   = "Security gate passed: %d finding(s) at or above %s (threshold: %d)"
	logErrorGateFail  = "Security gate failed: %d finding(s) at or above %s (threshold: %d)"
)

// Package securitygate implements the security-gate feature for scope-guardian.
// It evaluates a set of scan findings against a configurable severity threshold
// and signals whether the build should be blocked (gate failed) or allowed to
// continue (gate passed).
package securitygate

import (
	"fmt"
	"scope-guardian/domains/models"
	"scope-guardian/logger"
	"scope-guardian/parser"
	"strings"
)

// Evaluate counts findings whose severity is at or above the threshold severity
// and compares the count against the threshold value.
// It returns true when the gate passes (count < threshold.Value) and false when
// the gate fails (count >= threshold.Value).
func Evaluate(findings []models.Finding, threshold parser.Threshold) bool {
	thresholdRank, ok := severityRank[strings.ToUpper(threshold.Severity)]
	if !ok {
		return true
	}

	count := 0
	for _, f := range findings {
		rank, ok := severityRank[strings.ToUpper(f.Severity)]
		if ok && rank >= thresholdRank {
			count++
		}
	}

	if count >= threshold.Value {
		logger.Error(fmt.Sprintf(logErrorGateFail, count, strings.ToUpper(threshold.Severity), threshold.Value))
		return false
	}

	logger.Info(fmt.Sprintf(logInfoGatePass, count, strings.ToUpper(threshold.Severity), threshold.Value))
	return true
}

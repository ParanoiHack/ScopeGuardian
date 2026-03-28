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

// Evaluate checks every threshold against the findings.
// For each threshold it counts findings whose severity is at or above the threshold
// severity and fails the gate when that count reaches the threshold value.
// It returns false (gate failed) as soon as any threshold is breached, or true
// (gate passed) when all thresholds are satisfied.
func Evaluate(findings []models.Finding, thresholds []parser.Threshold) bool {
	for _, threshold := range thresholds {
		if !evaluateSingle(findings, threshold) {
			return false
		}
	}
	return true
}

// evaluateSingle evaluates a single threshold against the findings.
// It returns true when the gate passes (count < threshold.Value) and false when
// the gate fails (count >= threshold.Value).
func evaluateSingle(findings []models.Finding, threshold parser.Threshold) bool {
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

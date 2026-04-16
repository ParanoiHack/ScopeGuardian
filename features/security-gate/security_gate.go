// Package securitygate implements the security-gate feature for ScopeGuardian.
// It evaluates a set of scan findings against a configurable severity threshold
// and signals whether the build should be blocked (gate failed) or allowed to
// continue (gate passed).
package securitygate

import (
	"fmt"
	"ScopeGuardian/domains/models"
	"ScopeGuardian/logger"
	"ScopeGuardian/parser"
	"strings"
)

// Evaluate checks every threshold against the findings.
// For each threshold it counts findings whose severity is at or above the threshold
// severity and fails the gate when that count reaches the threshold value.
// It returns false (gate failed) as soon as any threshold is breached, logging only
// the failure. When all thresholds pass a single summary pass log is emitted.
func Evaluate(findings []models.Finding, thresholds []parser.Threshold) bool {
	for _, threshold := range thresholds {
		if !evaluateSingle(findings, threshold) {
			return false
		}
	}
	logger.Info(logInfoGatePass)
	return true
}

// evaluateSingle evaluates a single threshold against the findings.
// It returns true when the gate passes (count < threshold.Value) and false when
// the gate fails (count >= threshold.Value). It only logs on failure.
// Findings with status INACTIVE or DUPLICATE are excluded from the count: INACTIVE
// findings have been suppressed in DefectDojo (false positive, accepted risk, etc.)
// and DUPLICATE findings are already tracked elsewhere in the product — neither
// should block a build.
func evaluateSingle(findings []models.Finding, threshold parser.Threshold) bool {
	thresholdRank, ok := severityRank[strings.ToUpper(threshold.Severity)]
	if !ok {
		return true
	}

	count := 0
	for _, f := range findings {
		if f.Status == models.FindingStatusInactive || f.Status == models.FindingStatusDuplicate {
			continue
		}
		rank, ok := severityRank[strings.ToUpper(f.Severity)]
		if ok && rank >= thresholdRank {
			count++
		}
	}

	if count >= threshold.Value {
		logger.Error(fmt.Sprintf(logErrorGateFail, count, strings.ToUpper(threshold.Severity), threshold.Value))
		return false
	}

	return true
}

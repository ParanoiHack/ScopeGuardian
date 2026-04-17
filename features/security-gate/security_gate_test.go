package securitygate

import (
	"ScopeGuardian/domains/models"
	"ScopeGuardian/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate_NoThresholdExceeded(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH", Status: models.FindingStatusActive},
		{Severity: "MEDIUM", Status: models.FindingStatusActive},
	}
	// critical=1: no criticals found → count=0 < 1 → pass
	threshold := parser.Threshold{Severity: "critical", Value: 1}

	assert.True(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_ThresholdExceeded(t *testing.T) {
	findings := []models.Finding{
		{Severity: "CRITICAL", Status: models.FindingStatusActive},
		{Severity: "HIGH", Status: models.FindingStatusActive},
	}
	threshold := parser.Threshold{Severity: "critical", Value: 0}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_ExactlyAtThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH", Status: models.FindingStatusActive},
		{Severity: "HIGH", Status: models.FindingStatusActive},
	}
	// high=2: exactly 2 HIGH findings → count >= threshold → gate fails
	threshold := parser.Threshold{Severity: "high", Value: 2}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_OneHighWithThresholdHighOne_ShouldFail(t *testing.T) {
	// Reproducer: 1 HIGH finding with threshold high=1 must fail.
	findings := []models.Finding{
		{Severity: "HIGH", Status: models.FindingStatusActive},
	}
	threshold := parser.Threshold{Severity: "high", Value: 1}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_AboveThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH", Status: models.FindingStatusActive},
		{Severity: "CRITICAL", Status: models.FindingStatusActive},
		{Severity: "HIGH", Status: models.FindingStatusActive},
	}
	threshold := parser.Threshold{Severity: "high", Value: 2}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_CountsAllSeveritiesAtOrAbove(t *testing.T) {
	findings := []models.Finding{
		{Severity: "LOW", Status: models.FindingStatusActive},
		{Severity: "MEDIUM", Status: models.FindingStatusActive},
		{Severity: "HIGH", Status: models.FindingStatusActive},
		{Severity: "CRITICAL", Status: models.FindingStatusActive},
	}
	// medium=1 means threshold at MEDIUM; count of MEDIUM+HIGH+CRITICAL = 3 >= 1 → fail
	threshold := parser.Threshold{Severity: "medium", Value: 1}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_LowSeverityBelowThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "LOW", Status: models.FindingStatusActive},
		{Severity: "LOW", Status: models.FindingStatusActive},
	}
	// threshold high=1: LOW is below HIGH → not counted → count=0 < 1 → pass
	threshold := parser.Threshold{Severity: "high", Value: 1}

	assert.True(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_EmptyFindings(t *testing.T) {
	threshold := parser.Threshold{Severity: "critical", Value: 1}

	assert.True(t, Evaluate([]models.Finding{}, []parser.Threshold{threshold}))
}

func TestEvaluate_UnknownThresholdSeverity(t *testing.T) {
	findings := []models.Finding{
		{Severity: "CRITICAL", Status: models.FindingStatusActive},
	}
	threshold := parser.Threshold{Severity: "unknown", Value: 0}

	assert.True(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_CaseInsensitiveFindingSeverity(t *testing.T) {
	findings := []models.Finding{
		{Severity: "critical", Status: models.FindingStatusActive},
		{Severity: "High", Status: models.FindingStatusActive},
	}
	threshold := parser.Threshold{Severity: "high", Value: 0}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_InfoThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "INFO", Status: models.FindingStatusActive},
		{Severity: "LOW", Status: models.FindingStatusActive},
		{Severity: "MEDIUM", Status: models.FindingStatusActive},
	}
	// info=2: count all findings at INFO or above = 3 >= 2 → fail
	threshold := parser.Threshold{Severity: "info", Value: 2}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_MultipleThresholds_BothPass(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH", Status: models.FindingStatusActive},
	}
	// critical=1: 0 criticals < 1 → pass; high=2: 1 high < 2 → pass
	thresholds := []parser.Threshold{
		{Severity: "critical", Value: 1},
		{Severity: "high", Value: 2},
	}

	assert.True(t, Evaluate(findings, thresholds))
}

func TestEvaluate_MultipleThresholds_SecondBreached(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH", Status: models.FindingStatusActive},
	}
	// critical=1: 0 criticals < 1 → pass; high=1: 1 high >= 1 → fail
	thresholds := []parser.Threshold{
		{Severity: "critical", Value: 1},
		{Severity: "high", Value: 1},
	}

	assert.False(t, Evaluate(findings, thresholds))
}

func TestEvaluate_MultipleThresholds_FirstBreached(t *testing.T) {
	findings := []models.Finding{
		{Severity: "CRITICAL", Status: models.FindingStatusActive},
	}
	// critical=1: 1 critical >= 1 → fail (stops here)
	thresholds := []parser.Threshold{
		{Severity: "critical", Value: 1},
		{Severity: "high", Value: 5},
	}

	assert.False(t, Evaluate(findings, thresholds))
}

func TestEvaluate_DuplicateFindingExcludedFromGate(t *testing.T) {
	// A DUPLICATE finding that would otherwise breach the threshold must not
	// be counted — it is already tracked elsewhere in the product.
	findings := []models.Finding{
		{Severity: "CRITICAL", Status: models.FindingStatusDuplicate},
	}
	threshold := parser.Threshold{Severity: "critical", Value: 1}

	assert.True(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_InactiveFindingExcludedFromGate(t *testing.T) {
	// An INACTIVE finding (suppressed / false positive / accepted risk) must not
	// be counted — DefectDojo has explicitly dismissed it.
	findings := []models.Finding{
		{Severity: "CRITICAL", Status: models.FindingStatusInactive},
	}
	threshold := parser.Threshold{Severity: "critical", Value: 1}

	assert.True(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_MixedStatuses_OnlyActiveCountedTowardGate(t *testing.T) {
	// Mix of ACTIVE, INACTIVE, DUPLICATE: only the ACTIVE HIGH should be counted.
	findings := []models.Finding{
		{Severity: "HIGH", Status: models.FindingStatusActive},
		{Severity: "CRITICAL", Status: models.FindingStatusDuplicate},
		{Severity: "HIGH", Status: models.FindingStatusInactive},
	}
	// high=2: only 1 active HIGH → count=1 < 2 → pass
	threshold := parser.Threshold{Severity: "high", Value: 2}

	assert.True(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_AllDuplicateOrInactive_GatePasses(t *testing.T) {
	// When every finding is DUPLICATE or INACTIVE the gate must always pass,
	// regardless of severity, because the security team has already addressed them.
	findings := []models.Finding{
		{Severity: "CRITICAL", Status: models.FindingStatusDuplicate},
		{Severity: "HIGH", Status: models.FindingStatusInactive},
		{Severity: "MEDIUM", Status: models.FindingStatusDuplicate},
	}
	// threshold info=1: 0 active findings (all skipped) < 1 → pass
	threshold := parser.Threshold{Severity: "info", Value: 1}

	assert.True(t, Evaluate(findings, []parser.Threshold{threshold}))
}

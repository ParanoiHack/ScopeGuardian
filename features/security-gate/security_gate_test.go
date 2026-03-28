package securitygate

import (
	"scope-guardian/domains/models"
	"scope-guardian/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate_NoThresholdExceeded(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH"},
		{Severity: "MEDIUM"},
	}
	// critical=1: no criticals found → count=0 < 1 → pass
	threshold := parser.Threshold{Severity: "critical", Value: 1}

	assert.True(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_ThresholdExceeded(t *testing.T) {
	findings := []models.Finding{
		{Severity: "CRITICAL"},
		{Severity: "HIGH"},
	}
	threshold := parser.Threshold{Severity: "critical", Value: 0}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_ExactlyAtThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH"},
		{Severity: "HIGH"},
	}
	// high=2: exactly 2 HIGH findings → count >= threshold → gate fails
	threshold := parser.Threshold{Severity: "high", Value: 2}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_OneHighWithThresholdHighOne_ShouldFail(t *testing.T) {
	// Reproducer: 1 HIGH finding with threshold high=1 must fail.
	findings := []models.Finding{
		{Severity: "HIGH"},
	}
	threshold := parser.Threshold{Severity: "high", Value: 1}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}


func TestEvaluate_AboveThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH"},
		{Severity: "CRITICAL"},
		{Severity: "HIGH"},
	}
	threshold := parser.Threshold{Severity: "high", Value: 2}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_CountsAllSeveritiesAtOrAbove(t *testing.T) {
	findings := []models.Finding{
		{Severity: "LOW"},
		{Severity: "MEDIUM"},
		{Severity: "HIGH"},
		{Severity: "CRITICAL"},
	}
	// medium=1 means threshold at MEDIUM; count of MEDIUM+HIGH+CRITICAL = 3 >= 1 → fail
	threshold := parser.Threshold{Severity: "medium", Value: 1}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_LowSeverityBelowThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "LOW"},
		{Severity: "LOW"},
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
		{Severity: "CRITICAL"},
	}
	threshold := parser.Threshold{Severity: "unknown", Value: 0}

	assert.True(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_CaseInsensitiveFindingSeverity(t *testing.T) {
	findings := []models.Finding{
		{Severity: "critical"},
		{Severity: "High"},
	}
	threshold := parser.Threshold{Severity: "high", Value: 0}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_InfoThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "INFO"},
		{Severity: "LOW"},
		{Severity: "MEDIUM"},
	}
	// info=2: count all findings at INFO or above = 3 >= 2 → fail
	threshold := parser.Threshold{Severity: "info", Value: 2}

	assert.False(t, Evaluate(findings, []parser.Threshold{threshold}))
}

func TestEvaluate_MultipleThresholds_BothPass(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH"},
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
		{Severity: "HIGH"},
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
		{Severity: "CRITICAL"},
	}
	// critical=1: 1 critical >= 1 → fail (stops here)
	thresholds := []parser.Threshold{
		{Severity: "critical", Value: 1},
		{Severity: "high", Value: 5},
	}

	assert.False(t, Evaluate(findings, thresholds))
}

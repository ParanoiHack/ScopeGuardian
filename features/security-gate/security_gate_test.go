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
	threshold := parser.Threshold{Severity: "critical", Value: 0}

	assert.True(t, Evaluate(findings, threshold))
}

func TestEvaluate_ThresholdExceeded(t *testing.T) {
	findings := []models.Finding{
		{Severity: "CRITICAL"},
		{Severity: "HIGH"},
	}
	threshold := parser.Threshold{Severity: "critical", Value: 0}

	assert.False(t, Evaluate(findings, threshold))
}

func TestEvaluate_ExactlyAtThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH"},
		{Severity: "HIGH"},
	}
	threshold := parser.Threshold{Severity: "high", Value: 2}

	assert.True(t, Evaluate(findings, threshold))
}

func TestEvaluate_AboveThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "HIGH"},
		{Severity: "CRITICAL"},
		{Severity: "HIGH"},
	}
	threshold := parser.Threshold{Severity: "high", Value: 2}

	assert.False(t, Evaluate(findings, threshold))
}

func TestEvaluate_CountsAllSeveritiesAtOrAbove(t *testing.T) {
	findings := []models.Finding{
		{Severity: "LOW"},
		{Severity: "MEDIUM"},
		{Severity: "HIGH"},
		{Severity: "CRITICAL"},
	}
	// medium=1 means threshold at MEDIUM; count of MEDIUM+HIGH+CRITICAL = 3 > 1 → fail
	threshold := parser.Threshold{Severity: "medium", Value: 1}

	assert.False(t, Evaluate(findings, threshold))
}

func TestEvaluate_LowSeverityBelowThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "LOW"},
		{Severity: "LOW"},
	}
	// threshold high=0: LOW is below HIGH → not counted → count=0 <= 0 → pass
	threshold := parser.Threshold{Severity: "high", Value: 0}

	assert.True(t, Evaluate(findings, threshold))
}

func TestEvaluate_EmptyFindings(t *testing.T) {
	threshold := parser.Threshold{Severity: "critical", Value: 0}

	assert.True(t, Evaluate([]models.Finding{}, threshold))
}

func TestEvaluate_UnknownThresholdSeverity(t *testing.T) {
	findings := []models.Finding{
		{Severity: "CRITICAL"},
	}
	threshold := parser.Threshold{Severity: "unknown", Value: 0}

	assert.True(t, Evaluate(findings, threshold))
}

func TestEvaluate_CaseInsensitiveFindingSeverity(t *testing.T) {
	findings := []models.Finding{
		{Severity: "critical"},
		{Severity: "High"},
	}
	threshold := parser.Threshold{Severity: "high", Value: 0}

	assert.False(t, Evaluate(findings, threshold))
}

func TestEvaluate_InfoThreshold(t *testing.T) {
	findings := []models.Finding{
		{Severity: "INFO"},
		{Severity: "LOW"},
		{Severity: "MEDIUM"},
	}
	// info=2: count all findings at INFO or above = 3 > 2 → fail
	threshold := parser.Threshold{Severity: "info", Value: 2}

	assert.False(t, Evaluate(findings, threshold))
}

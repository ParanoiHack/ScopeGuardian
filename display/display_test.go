package display

import (
	"bytes"
	"scope-guardian/domains/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplayBanner(t *testing.T) {
	assert.NotPanics(t, func() {
		DisplayBanner(&bytes.Buffer{})
	})
}

func TestDisplayCredit(t *testing.T) {
	assert.NotPanics(t, func() {
		DisplayCredit(&bytes.Buffer{})
	})
}

func TestDisplayFindings_Empty(t *testing.T) {
	assert.NotPanics(t, func() {
		DisplayFindings(&bytes.Buffer{}, []models.Finding{})
	})
}

func TestDisplayFindings_WithData(t *testing.T) {
	findings := []models.Finding{
		{
			Engine:         "IACST",
			Severity:       "HIGH",
			Name:           "Privileged Container",
			Cwe:            "CWE-284",
			Description:    "Container is running as privileged",
			SinkFile:       "Dockerfile",
			SinkLine:       10,
			Recommendation: "Do not run containers as privileged",
		},
		{
			Engine:         "IACST",
			Severity:       "MEDIUM",
			Name:           "Exposed Port",
			Cwe:            "CWE-200",
			Description:    "Sensitive port exposed",
			SinkFile:       "docker-compose.yml",
			SinkLine:       25,
			Recommendation: "Avoid exposing sensitive ports",
		},
	}

	assert.NotPanics(t, func() {
		DisplayFindings(&bytes.Buffer{}, findings)
	})
}

func TestDisplayFindings_NilFindings(t *testing.T) {
	assert.NotPanics(t, func() {
		DisplayFindings(&bytes.Buffer{}, nil)
	})
}

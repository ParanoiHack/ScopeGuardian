package display

import (
	"bytes"
	"encoding/json"
	"strings"
	"ScopeGuardian/domains/models"
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

var testFindings = []models.Finding{
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
}

func TestDumpFindings_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := DumpFindings(&buf, testFindings, "json")

	assert.Nil(t, err)

	var result []models.Finding
	assert.Nil(t, json.Unmarshal(buf.Bytes(), &result))
	assert.Len(t, result, 1)
	assert.EqualValues(t, testFindings[0].Engine, result[0].Engine)
	assert.EqualValues(t, testFindings[0].Severity, result[0].Severity)
	assert.EqualValues(t, testFindings[0].Name, result[0].Name)
}

func TestDumpFindings_CSV(t *testing.T) {
	var buf bytes.Buffer
	err := DumpFindings(&buf, testFindings, "csv")

	assert.Nil(t, err)

	output := buf.String()
	assert.Contains(t, output, "Engine")
	assert.Contains(t, output, "Severity")
	assert.Contains(t, output, "IACST")
	assert.Contains(t, output, "HIGH")
	assert.Contains(t, output, "Privileged Container")
}

func TestDumpFindings_Raw(t *testing.T) {
	var buf bytes.Buffer
	err := DumpFindings(&buf, testFindings, "raw")

	assert.Nil(t, err)
	assert.NotEmpty(t, buf.String())
}

func TestDumpFindings_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := DumpFindings(&buf, []models.Finding{}, "json")

	assert.Nil(t, err)

	var result []models.Finding
	assert.Nil(t, json.Unmarshal(buf.Bytes(), &result))
	assert.Empty(t, result)
}

func TestDumpFindings_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	err := DumpFindings(&buf, testFindings, "xml")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "xml")
}

func TestDumpFindings_CSV_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	err := DumpFindings(&buf, []models.Finding{}, "csv")

	assert.Nil(t, err)
	output := buf.String()
	assert.True(t, strings.HasPrefix(output, "Engine,Severity,Name,CWE,Description,SinkFile,SinkLine,Recommendation"))
}

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterInactiveFindings(t *testing.T) {
	active := Finding{Status: FindingStatusActive, Name: "SQL Injection"}
	duplicate := Finding{Status: FindingStatusDuplicate, Name: "XSS"}
	inactive := Finding{Status: FindingStatusInactive, Name: "Old Issue"}

	t.Run("removes INACTIVE findings and keeps ACTIVE and DUPLICATE", func(t *testing.T) {
		input := []Finding{active, inactive, duplicate}
		got := FilterInactiveFindings(input)
		assert.Equal(t, []Finding{active, duplicate}, got)
	})

	t.Run("returns empty slice when all findings are INACTIVE", func(t *testing.T) {
		input := []Finding{inactive, inactive}
		got := FilterInactiveFindings(input)
		assert.Empty(t, got)
	})

	t.Run("returns all findings when none are INACTIVE", func(t *testing.T) {
		input := []Finding{active, duplicate}
		got := FilterInactiveFindings(input)
		assert.Equal(t, []Finding{active, duplicate}, got)
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		got := FilterInactiveFindings([]Finding{})
		assert.Empty(t, got)
	})

	t.Run("returns nil-safe result for nil input", func(t *testing.T) {
		got := FilterInactiveFindings(nil)
		assert.Empty(t, got)
	})
}

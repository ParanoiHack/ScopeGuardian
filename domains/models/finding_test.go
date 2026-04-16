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

func TestFilterFindingsByStatus(t *testing.T) {
	active := Finding{Status: FindingStatusActive, Name: "SQL Injection"}
	duplicate := Finding{Status: FindingStatusDuplicate, Name: "XSS"}
	inactive := Finding{Status: FindingStatusInactive, Name: "Old Issue"}

	t.Run("returns only ACTIVE findings when filter is ACTIVE", func(t *testing.T) {
		input := []Finding{active, inactive, duplicate}
		got := FilterFindingsByStatus(input, []string{FindingStatusActive})
		assert.Equal(t, []Finding{active}, got)
	})

	t.Run("returns ACTIVE and DUPLICATE when filter is ACTIVE,DUPLICATE", func(t *testing.T) {
		input := []Finding{active, inactive, duplicate}
		got := FilterFindingsByStatus(input, []string{FindingStatusActive, FindingStatusDuplicate})
		assert.Equal(t, []Finding{active, duplicate}, got)
	})

	t.Run("returns only INACTIVE findings when filter is INACTIVE", func(t *testing.T) {
		input := []Finding{active, inactive, duplicate}
		got := FilterFindingsByStatus(input, []string{FindingStatusInactive})
		assert.Equal(t, []Finding{inactive}, got)
	})

	t.Run("returns all findings when filter contains all statuses", func(t *testing.T) {
		input := []Finding{active, inactive, duplicate}
		got := FilterFindingsByStatus(input, []string{FindingStatusActive, FindingStatusInactive, FindingStatusDuplicate})
		assert.Equal(t, []Finding{active, inactive, duplicate}, got)
	})

	t.Run("returns original slice when statuses is empty", func(t *testing.T) {
		input := []Finding{active, inactive, duplicate}
		got := FilterFindingsByStatus(input, []string{})
		assert.Equal(t, input, got)
	})

	t.Run("returns original slice when statuses is nil", func(t *testing.T) {
		input := []Finding{active, inactive, duplicate}
		got := FilterFindingsByStatus(input, nil)
		assert.Equal(t, input, got)
	})

	t.Run("is case-insensitive", func(t *testing.T) {
		input := []Finding{active, inactive}
		got := FilterFindingsByStatus(input, []string{"active"})
		assert.Equal(t, []Finding{active}, got)
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		got := FilterFindingsByStatus([]Finding{}, []string{FindingStatusActive})
		assert.Empty(t, got)
	})
}

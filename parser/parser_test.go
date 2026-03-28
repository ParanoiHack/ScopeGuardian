package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("Should parse config filepath as positional argument", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "./config.toml"})

		assert.Nil(t, err)
		assert.EqualValues(t, "./config.toml", args.Config)
		assert.EqualValues(t, "my-project", args.ProjectName)
		assert.EqualValues(t, "main", args.Branch)
		assert.EqualValues(t, false, args.Sync)
		assert.Nil(t, args.Thresholds)
	})

	t.Run("Should parse sync flag as true", func(t *testing.T) {
		args, err := Parse([]string{"--sync", "--projectName", "my-project", "--branch", "main", "./config.toml"})

		assert.Nil(t, err)
		assert.EqualValues(t, "./config.toml", args.Config)
		assert.EqualValues(t, "my-project", args.ProjectName)
		assert.EqualValues(t, "main", args.Branch)
		assert.EqualValues(t, true, args.Sync)
		assert.Empty(t, args.Thresholds)
	})

	t.Run("Should parse threshold with critical severity", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "critical=1", "./config.toml"})

		assert.Nil(t, err)
		assert.EqualValues(t, "./config.toml", args.Config)
		assert.Len(t, args.Thresholds, 1)
		assert.EqualValues(t, severityCritical, args.Thresholds[0].Severity)
		assert.EqualValues(t, 1, args.Thresholds[0].Value)
	})

	t.Run("Should parse threshold with high severity", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "high=0", "./config.toml"})

		assert.Nil(t, err)
		assert.Len(t, args.Thresholds, 1)
		assert.EqualValues(t, severityHigh, args.Thresholds[0].Severity)
		assert.EqualValues(t, 0, args.Thresholds[0].Value)
	})

	t.Run("Should parse threshold with medium severity", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "medium=3", "./config.toml"})

		assert.Nil(t, err)
		assert.Len(t, args.Thresholds, 1)
		assert.EqualValues(t, severityMedium, args.Thresholds[0].Severity)
		assert.EqualValues(t, 3, args.Thresholds[0].Value)
	})

	t.Run("Should parse threshold with low severity", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "low=5", "./config.toml"})

		assert.Nil(t, err)
		assert.Len(t, args.Thresholds, 1)
		assert.EqualValues(t, severityLow, args.Thresholds[0].Severity)
		assert.EqualValues(t, 5, args.Thresholds[0].Value)
	})

	t.Run("Should parse threshold with info severity", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "info=2", "./config.toml"})

		assert.Nil(t, err)
		assert.Len(t, args.Thresholds, 1)
		assert.EqualValues(t, severityInfo, args.Thresholds[0].Severity)
		assert.EqualValues(t, 2, args.Thresholds[0].Value)
	})

	t.Run("Should parse all flags together", func(t *testing.T) {
		args, err := Parse([]string{"--sync", "--projectName", "my-project", "--branch", "main", "--threshold", "critical=1", "./config.toml"})

		assert.Nil(t, err)
		assert.EqualValues(t, "./config.toml", args.Config)
		assert.EqualValues(t, "my-project", args.ProjectName)
		assert.EqualValues(t, "main", args.Branch)
		assert.EqualValues(t, true, args.Sync)
		assert.Len(t, args.Thresholds, 1)
		assert.EqualValues(t, severityCritical, args.Thresholds[0].Severity)
		assert.EqualValues(t, 1, args.Thresholds[0].Value)
	})

	t.Run("Should parse multiple comma-separated thresholds", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "critical=1,high=2", "./config.toml"})

		assert.Nil(t, err)
		assert.Len(t, args.Thresholds, 2)
		assert.EqualValues(t, severityCritical, args.Thresholds[0].Severity)
		assert.EqualValues(t, 1, args.Thresholds[0].Value)
		assert.EqualValues(t, severityHigh, args.Thresholds[1].Severity)
		assert.EqualValues(t, 2, args.Thresholds[1].Value)
	})

	t.Run("Should not parse when config filepath is missing", func(t *testing.T) {
		args, err := Parse([]string{})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
		assert.EqualValues(t, errConfigRequired, err.Error())
	})

	t.Run("Should not parse when projectName is missing", func(t *testing.T) {
		args, err := Parse([]string{"--branch", "main", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
		assert.EqualValues(t, errProjectNameRequired, err.Error())
	})

	t.Run("Should not parse when branch is missing", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
		assert.EqualValues(t, errBranchRequired, err.Error())
	})

	t.Run("Should not parse when threshold format is invalid", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "critical", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
		assert.EqualValues(t, errInvalidThreshold, err.Error())
	})

	t.Run("Should not parse when threshold severity is invalid", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "unknown=1", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
	})

	t.Run("Should not parse when threshold value is not an integer", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "critical=abc", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
	})

	t.Run("Should not parse when threshold value is negative", func(t *testing.T) {
		args, err := Parse([]string{"--projectName", "my-project", "--branch", "main", "--threshold", "critical=-1", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
	})
}

func TestPrintUsage(t *testing.T) {
	var buf bytes.Buffer
	PrintUsage(&buf)

	output := buf.String()
	assert.Contains(t, output, "scope-guardian")
	assert.Contains(t, output, "--projectName")
	assert.Contains(t, output, "--branch")
	assert.Contains(t, output, "--sync")
	assert.Contains(t, output, "--threshold")
	assert.Contains(t, output, "<config-file>")
}

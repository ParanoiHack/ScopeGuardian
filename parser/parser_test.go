package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("Should parse config filepath as positional argument", func(t *testing.T) {
		args, err := Parse([]string{"./config.toml"})

		assert.Nil(t, err)
		assert.EqualValues(t, "./config.toml", args.Config)
		assert.EqualValues(t, false, args.Sync)
		assert.Nil(t, args.Threshold)
	})

	t.Run("Should parse sync flag as true", func(t *testing.T) {
		args, err := Parse([]string{"--sync", "./config.toml"})

		assert.Nil(t, err)
		assert.EqualValues(t, "./config.toml", args.Config)
		assert.EqualValues(t, true, args.Sync)
		assert.Nil(t, args.Threshold)
	})

	t.Run("Should parse threshold with critical severity", func(t *testing.T) {
		args, err := Parse([]string{"--threshold", "critical=1", "./config.toml"})

		assert.Nil(t, err)
		assert.EqualValues(t, "./config.toml", args.Config)
		assert.NotNil(t, args.Threshold)
		assert.EqualValues(t, severityCritical, args.Threshold.Severity)
		assert.EqualValues(t, 1, args.Threshold.Value)
	})

	t.Run("Should parse threshold with high severity", func(t *testing.T) {
		args, err := Parse([]string{"--threshold", "high=0", "./config.toml"})

		assert.Nil(t, err)
		assert.NotNil(t, args.Threshold)
		assert.EqualValues(t, severityHigh, args.Threshold.Severity)
		assert.EqualValues(t, 0, args.Threshold.Value)
	})

	t.Run("Should parse threshold with medium severity", func(t *testing.T) {
		args, err := Parse([]string{"--threshold", "medium=3", "./config.toml"})

		assert.Nil(t, err)
		assert.NotNil(t, args.Threshold)
		assert.EqualValues(t, severityMedium, args.Threshold.Severity)
		assert.EqualValues(t, 3, args.Threshold.Value)
	})

	t.Run("Should parse threshold with low severity", func(t *testing.T) {
		args, err := Parse([]string{"--threshold", "low=5", "./config.toml"})

		assert.Nil(t, err)
		assert.NotNil(t, args.Threshold)
		assert.EqualValues(t, severityLow, args.Threshold.Severity)
		assert.EqualValues(t, 5, args.Threshold.Value)
	})

	t.Run("Should parse threshold with info severity", func(t *testing.T) {
		args, err := Parse([]string{"--threshold", "info=2", "./config.toml"})

		assert.Nil(t, err)
		assert.NotNil(t, args.Threshold)
		assert.EqualValues(t, severityInfo, args.Threshold.Severity)
		assert.EqualValues(t, 2, args.Threshold.Value)
	})

	t.Run("Should parse all flags together", func(t *testing.T) {
		args, err := Parse([]string{"--sync", "--threshold", "critical=1", "./config.toml"})

		assert.Nil(t, err)
		assert.EqualValues(t, "./config.toml", args.Config)
		assert.EqualValues(t, true, args.Sync)
		assert.NotNil(t, args.Threshold)
		assert.EqualValues(t, severityCritical, args.Threshold.Severity)
		assert.EqualValues(t, 1, args.Threshold.Value)
	})

	t.Run("Should not parse when config filepath is missing", func(t *testing.T) {
		args, err := Parse([]string{})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
		assert.EqualValues(t, errConfigRequired, err.Error())
	})

	t.Run("Should not parse when threshold format is invalid", func(t *testing.T) {
		args, err := Parse([]string{"--threshold", "critical", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
		assert.EqualValues(t, errInvalidThreshold, err.Error())
	})

	t.Run("Should not parse when threshold severity is invalid", func(t *testing.T) {
		args, err := Parse([]string{"--threshold", "unknown=1", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
	})

	t.Run("Should not parse when threshold value is not an integer", func(t *testing.T) {
		args, err := Parse([]string{"--threshold", "critical=abc", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
	})

	t.Run("Should not parse when threshold value is negative", func(t *testing.T) {
		args, err := Parse([]string{"--threshold", "critical=-1", "./config.toml"})

		assert.NotNil(t, err)
		assert.EqualValues(t, Args{}, args)
	})
}

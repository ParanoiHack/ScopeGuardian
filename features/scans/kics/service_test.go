package kics

import (
	"fmt"
	"os"
	"scope-guardian/domains/interfaces"
	environment_variable "scope-guardian/environnement_variable"
	"scope-guardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKicsService(t *testing.T) {
	service := newKicsService(loader.Kics{"./test", ""})

	_, ok := service.(interfaces.ScanServiceImpl)
	assert.NotNil(t, service)
	assert.True(t, ok)
}

func TestVerifyConfig(t *testing.T) {
	t.Run("Config should be OK", func(t *testing.T) {
		ok, err := verifyConfig(fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks"))
		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Config should not be OK", func(t *testing.T) {
		ok, err := verifyConfig(fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/doesnotexist"))
		assert.NotNil(t, err)
		assert.False(t, ok)
	})
}

func TestLoadFinding(t *testing.T) {
	t.Run("Should load findings", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newKicsService(loader.Kics{"./test", ""})

		findings, err := service.LoadFindings()

		assert.Nil(t, err)
		assert.EqualValues(t, 8, len(findings))
	})

	t.Run("Should not load findings due to lack of results", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newKicsService(loader.Kics{"./test", ""})

		findings, err := service.LoadFindings()

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})

	t.Run("Should not load findings due to unmarshalling issue", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/bad_format_results"))
		environment_variable.ReloadEnv()

		service := newKicsService(loader.Kics{"./test", ""})

		findings, err := service.LoadFindings()

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})
}

// func TestSync(t *testing.T) {
// 	service := newKicsService(loader.Kics{"./test", ""})

// 	err := service.Sync()

// 	assert.Nil(t, err)
// }

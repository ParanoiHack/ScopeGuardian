package syft

import (
	"fmt"
	"os"
	"scope-guardian/domains/interfaces"
	environment_variable "scope-guardian/environnement_variable"
	"scope-guardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSyftService(t *testing.T) {
	service := newSyftService(loader.Grype{Path: "./test"})

	_, ok := service.(interfaces.ScanServiceImpl)
	assert.NotNil(t, service)
	assert.True(t, ok)
}

func TestSyftStart(t *testing.T) {
	t.Run("Should return error when directory not found", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newSyftService(loader.Grype{Path: "./doesnotexist"})

		ok, err := service.Start()

		assert.NotNil(t, err)
		assert.False(t, ok)
		assert.EqualValues(t, errDirectoryNotFound, err.Error())
	})
}

func TestSyftLoadFindings(t *testing.T) {
	t.Run("Should return nil findings and nil error", func(t *testing.T) {
		service := newSyftService(loader.Grype{Path: "./test"})

		findings, err := service.LoadFindings()

		assert.Nil(t, err)
		assert.Nil(t, findings)
	})
}

func TestSyftSync(t *testing.T) {
	t.Run("Should return nil error", func(t *testing.T) {
		service := newSyftService(loader.Grype{Path: "./test"})

		err := service.Sync(1, "main", nil)

		assert.Nil(t, err)
	})
}

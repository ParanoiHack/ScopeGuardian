package syft

import (
	"fmt"
	"os"
	"scope-guardian/domains/interfaces"
	environment_variable "scope-guardian/environnement_variable"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSyftService(t *testing.T) {
	service := newSyftService("./test", false)

	_, ok := service.(interfaces.ScanServiceImpl)
	assert.NotNil(t, service)
	assert.True(t, ok)
}

func TestSyftStart(t *testing.T) {
	t.Run("Should return error when directory not found", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newSyftService("./doesnotexist", false)

		ok, err := service.Start()

		assert.NotNil(t, err)
		assert.False(t, ok)
		assert.EqualValues(t, errDirectoryNotFound, err.Error())
	})

	t.Run("Should return error when binary not found but directory exists", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", os.TempDir())
		environment_variable.ReloadEnv()

		service := newSyftService(".", false)

		ok, err := service.Start()

		assert.NotNil(t, err)
		assert.False(t, ok)
	})

	t.Run("Should log transitive libraries message and return error when binary not found", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", os.TempDir())
		environment_variable.ReloadEnv()

		service := newSyftService(".", true)

		ok, err := service.Start()

		assert.NotNil(t, err)
		assert.False(t, ok)
	})
}

func TestSyftLoadFindings(t *testing.T) {
	t.Run("Should return nil findings and nil error", func(t *testing.T) {
		service := newSyftService("./test", false)

		findings, err := service.LoadFindings()

		assert.Nil(t, err)
		assert.Nil(t, findings)
	})
}

func TestSyftSync(t *testing.T) {
	t.Run("Should return nil error", func(t *testing.T) {
		service := newSyftService("./test", false)

		err := service.Sync(1, "main", nil)

		assert.Nil(t, err)
	})
}

package syft

import (
	"ScopeGuardian/domains/interfaces"
	environment_variable "ScopeGuardian/environnement_variable"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSyftService(t *testing.T) {
	service := newSyftService("./test", false, nil, 0, nil)

	_, ok := service.(interfaces.ScanServiceImpl)
	assert.NotNil(t, service)
	assert.True(t, ok)
}

func TestSyftStart(t *testing.T) {
	t.Run("Should return error when directory not found", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newSyftService("./doesnotexist", false, nil, 0, nil)

		ok, err := service.Start()

		assert.NotNil(t, err)
		assert.False(t, ok)
		assert.EqualValues(t, errDirectoryNotFound, err.Error())
	})

	t.Run("Should return error when runner fails and directory exists", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", os.TempDir())
		environment_variable.ReloadEnv()

		svc := newSyftService(".", false, nil, 0, nil).(*SyftServiceImpl)
		svc.runner = func(_ string, _ string, _ []string, _ io.Writer, _ io.Writer, _ ...string) (bool, error) {
			return false, fmt.Errorf("runner error")
		}

		ok, err := svc.Start()

		assert.NotNil(t, err)
		assert.False(t, ok)
	})

	t.Run("Should log transitive libraries message and succeed when runner succeeds", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", os.TempDir())
		environment_variable.ReloadEnv()

		svc := newSyftService(".", true, nil, 0, nil).(*SyftServiceImpl)
		svc.runner = func(_ string, _ string, _ []string, _ io.Writer, _ io.Writer, _ ...string) (bool, error) {
			return true, nil
		}

		ok, err := svc.Start()

		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Should pass each exclude pattern as a separate --exclude arg (quoted)", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", os.TempDir())
		environment_variable.ReloadEnv()

		patterns := []string{"**/src/test/**", "**/testdata/**"}
		var capturedArgs []string
		svc := newSyftService(".", false, patterns, 1, nil).(*SyftServiceImpl)
		svc.runner = func(_ string, _ string, args []string, _ io.Writer, _ io.Writer, _ ...string) (bool, error) {
			capturedArgs = args
			return true, nil
		}

		ok, err := svc.Start()

		assert.Nil(t, err)
		assert.True(t, ok)
		assert.Contains(t, capturedArgs, excludeArgument)
		for _, p := range patterns {
			assert.Contains(t, capturedArgs, fmt.Sprintf(`"%s"`, p))
		}
	})

	t.Run("Should not pass --exclude arg when exclude list is empty", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", os.TempDir())
		environment_variable.ReloadEnv()

		var capturedArgs []string
		svc := newSyftService(".", false, nil, 0, nil).(*SyftServiceImpl)
		svc.runner = func(_ string, _ string, args []string, _ io.Writer, _ io.Writer, _ ...string) (bool, error) {
			capturedArgs = args
			return true, nil
		}

		ok, err := svc.Start()

		assert.Nil(t, err)
		assert.True(t, ok)
		assert.NotContains(t, capturedArgs, excludeArgument)
	})

	t.Run("Should pass SYFT_JAVA_MAX_PARENT_RECURSIVE_DEPTH env var with configured depth", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", os.TempDir())
		environment_variable.ReloadEnv()

		var capturedEnv []string
		svc := newSyftService(".", false, nil, 3, nil).(*SyftServiceImpl)
		svc.runner = func(_ string, _ string, _ []string, _ io.Writer, _ io.Writer, extraEnv ...string) (bool, error) {
			capturedEnv = extraEnv
			return true, nil
		}

		ok, err := svc.Start()

		assert.Nil(t, err)
		assert.True(t, ok)
		assert.Contains(t, capturedEnv, "SYFT_JAVA_MAX_PARENT_RECURSIVE_DEPTH=3")
	})

	t.Run("Should pass SYFT_JAVA_MAX_PARENT_RECURSIVE_DEPTH=1 when default depth is used", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", os.TempDir())
		environment_variable.ReloadEnv()

		var capturedEnv []string
		svc := newSyftService(".", false, nil, 1, nil).(*SyftServiceImpl)
		svc.runner = func(_ string, _ string, _ []string, _ io.Writer, _ io.Writer, extraEnv ...string) (bool, error) {
			capturedEnv = extraEnv
			return true, nil
		}

		ok, err := svc.Start()

		assert.Nil(t, err)
		assert.True(t, ok)
		assert.Contains(t, capturedEnv, "SYFT_JAVA_MAX_PARENT_RECURSIVE_DEPTH=1")
	})
}

func TestSyftLoadFindings(t *testing.T) {
	t.Run("Should return nil findings and nil error", func(t *testing.T) {
		service := newSyftService("./test", false, nil, 0, nil)

		findings, err := service.LoadFindings()

		assert.Nil(t, err)
		assert.Nil(t, findings)
	})
}

func TestSyftSync(t *testing.T) {
	t.Run("Should return nil error", func(t *testing.T) {
		service := newSyftService("./test", false, nil, 0, nil)

		err := service.Sync(1, "main", nil)

		assert.Nil(t, err)
	})
}

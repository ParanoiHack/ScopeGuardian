package environment_variable

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoadSuccess(t *testing.T) {
	_ = os.Setenv("test", "test")
	loadEnv()

	assert.Equal(t, "test", EnvironmentVariable["test"])
}

func TestReloadEnv(t *testing.T) {
	_ = os.Setenv("test", "test")
	loadEnv()

	assert.Equal(t, "test", EnvironmentVariable["test"])

	_ = os.Setenv("test", "test_2")
	ReloadEnv()

	assert.Equal(t, "test_2", EnvironmentVariable["test"])
}

func TestConfigLoadFail(t *testing.T) {
	assert.Equal(t, "", EnvironmentVariable["fail"])
}

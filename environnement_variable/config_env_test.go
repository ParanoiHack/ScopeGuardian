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

func TestConfigLoadFail(t *testing.T) {
	assert.Equal(t, "", EnvironmentVariable["fail"])
}

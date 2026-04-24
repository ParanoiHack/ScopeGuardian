package exec

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	t.Run("Should run a binary", func(t *testing.T) {
		ok, err := Wrap("/bin/ls", "/", []string{"-l", "/opt"}, os.Stdout, os.Stderr)

		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Should return error for non-existent binary", func(t *testing.T) {
		ok, err := Wrap("/bin/doesnotexist", "/", []string{"-l", "/opt/kics"}, io.Discard, io.Discard)

		assert.NotNil(t, err)
		assert.False(t, ok)
	})
}

func TestWrapAllowExitCodes(t *testing.T) {
	t.Run("Should treat allowed non-zero exit code as success", func(t *testing.T) {
		// 'false' exits with code 1
		ok, err := WrapAllowExitCodes("/bin/false", "/", []string{}, io.Discard, io.Discard, []int{1})

		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Should return error for disallowed non-zero exit code", func(t *testing.T) {
		// 'false' exits with code 1; only code 2 is allowed here
		ok, err := WrapAllowExitCodes("/bin/false", "/", []string{}, io.Discard, io.Discard, []int{2})

		assert.NotNil(t, err)
		assert.False(t, ok)
	})

	t.Run("Should succeed normally on exit 0", func(t *testing.T) {
		ok, err := WrapAllowExitCodes("/bin/true", "/", []string{}, io.Discard, io.Discard, []int{1})

		assert.Nil(t, err)
		assert.True(t, ok)
	})
}

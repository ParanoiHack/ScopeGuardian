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

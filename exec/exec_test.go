package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	t.Run("Should run a binary", func(t *testing.T) {
		ok, err := Wrap("/bin/ls", "/", []string{"-l", "/opt"})

		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Should run a binary", func(t *testing.T) {
		ok, err := Wrap("/bin/doesnotexist", "/", []string{"-l", "/opt/kics"})

		assert.NotNil(t, err)
		assert.False(t, ok)
	})
}

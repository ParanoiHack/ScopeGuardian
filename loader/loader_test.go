package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("Should load configuration file", func(t *testing.T) {
		config, err := Load("./mocks/config.toml")

		assert.Nil(t, err)
		assert.EqualValues(t, "Scope-guardian configuration file", config.Title)
		assert.EqualValues(t, "./", config.Kics.Path)
		assert.EqualValues(t, "terraform", config.Kics.Platform)
	})

	t.Run("Should load configuration file without engine", func(t *testing.T) {
		config, err := Load("./mocks/config_no_engine.toml")

		assert.Nil(t, err)
		assert.EqualValues(t, "Scope-guardian configuration file", config.Title)
		assert.EqualValues(t, Kics{}, config.Kics)
	})

	t.Run("Should not load configuration file cause wrong pathname", func(t *testing.T) {
		config, err := Load("./mocks/does_not_exist.toml")

		assert.NotNil(t, err)
		assert.EqualValues(t, Config{}, config)
		assert.EqualValues(t, errFileNotFound, err.Error())
	})

	t.Run("Should not load configuration file cause corruption", func(t *testing.T) {
		config, err := Load("./mocks/config_corrupted.toml")

		assert.NotNil(t, err)
		assert.EqualValues(t, Config{}, config)
		assert.EqualValues(t, errDecodingToml, err.Error())
	})
}

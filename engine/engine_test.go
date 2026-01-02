package engine

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testEngine = "test_engine"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()

	assert.NotNil(t, engine)
	assert.EqualValues(t, 0, len(engine.scanners))
}

func TestRegisterScanner(t *testing.T) {
	t.Run("Should register scanner", func(t *testing.T) {
		engine := NewEngine()

		assert.EqualValues(t, 0, len(engine.scanners))

		var mockScanner interfaces.ScanServiceImpl
		ok := engine.registerScanner(testEngine, mockScanner)

		assert.True(t, ok)
		assert.EqualValues(t, 1, len(engine.scanners))
	})

	t.Run("Should not register two scanner under the same name", func(t *testing.T) {
		engine := NewEngine()

		assert.EqualValues(t, 0, len(engine.scanners))

		var mockScanner interfaces.ScanServiceImpl
		engine.registerScanner(testEngine, mockScanner)

		ok := engine.registerScanner(testEngine, mockScanner)
		assert.False(t, ok)
		assert.EqualValues(t, 1, len(engine.scanners))
	})

	t.Run("Should not register a scanner with an empty name", func(t *testing.T) {
		engine := NewEngine()

		var mockScanner interfaces.ScanServiceImpl
		ok := engine.registerScanner("", mockScanner)

		assert.False(t, ok)
		assert.EqualValues(t, 0, len(engine.scanners))
	})
}

func TestInitialize(t *testing.T) {
	t.Run("Should initialize engine with kics runner", func(t *testing.T) {
		engine := NewEngine()

		config, err := loader.Load("../loader/mocks/config.toml")
		assert.Nil(t, err)
		assert.NotNil(t, config)

		engine.Initialize(config)

		assert.EqualValues(t, 1, len(engine.scanners))
	})

	t.Run("Should not initialize engine with kics runner", func(t *testing.T) {
		engine := NewEngine()

		config, err := loader.Load("../loader/mocks/config_no_engine.toml")
		assert.Nil(t, err)
		assert.NotNil(t, config)

		engine.Initialize(config)

		assert.EqualValues(t, 0, len(engine.scanners))
	})
}

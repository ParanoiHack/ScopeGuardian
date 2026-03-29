package engine

import (
	"errors"
	"scope-guardian/connectors/defectdojo"
	"scope-guardian/domains/interfaces"
	"scope-guardian/domains/models"
	"scope-guardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testEngine = "test_engine"
)

type mockScanService struct {
	startOk     bool
	startErr    error
	findings    []models.Finding
	findingsErr error
	syncErr     error
}

func (m *mockScanService) Start() (bool, error) {
	return m.startOk, m.startErr
}

func (m *mockScanService) LoadFindings() ([]models.Finding, error) {
	return m.findings, m.findingsErr
}

func (m *mockScanService) Sync(_ int, _ string, _ defectdojo.DefectDojoService) error {
	return m.syncErr
}

var _ interfaces.ScanServiceImpl = &mockScanService{}

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

	t.Run("Should initialize engine with syft runner when grype is configured", func(t *testing.T) {
		engine := NewEngine()

		config, err := loader.Load("../loader/mocks/config_with_grype.toml")
		assert.Nil(t, err)
		assert.NotNil(t, config)

		engine.Initialize(config)

		assert.EqualValues(t, 1, len(engine.scanners))
		_, ok := engine.scanners[syftScannerName]
		assert.True(t, ok)
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

func TestStart(t *testing.T) {
	t.Run("Should start scanner successfully", func(t *testing.T) {
		engine := NewEngine()
		mock := &mockScanService{startOk: true, startErr: nil}
		engine.registerScanner(testEngine, mock)

		assert.NotPanics(t, func() {
			engine.Start()
		})
	})

	t.Run("Should handle scanner start failure", func(t *testing.T) {
		engine := NewEngine()
		mock := &mockScanService{startOk: false, startErr: errors.New("scan failed")}
		engine.registerScanner(testEngine, mock)

		assert.NotPanics(t, func() {
			engine.Start()
		})
	})

	t.Run("Should complete with no scanners registered", func(t *testing.T) {
		engine := NewEngine()

		assert.NotPanics(t, func() {
			engine.Start()
		})
	})
}

func TestLoadFindings(t *testing.T) {
	t.Run("Should load findings from registered scanner", func(t *testing.T) {
		engine := NewEngine()
		expectedFindings := []models.Finding{
			{Engine: "IACST", Severity: "HIGH", Name: "Test Finding"},
		}
		mock := &mockScanService{findings: expectedFindings, findingsErr: nil}
		engine.registerScanner(testEngine, mock)

		results := engine.LoadFindings()

		assert.EqualValues(t, 1, len(results))
		assert.Equal(t, "IACST", results[0].Engine)
	})

	t.Run("Should return empty slice when scanner returns error", func(t *testing.T) {
		engine := NewEngine()
		mock := &mockScanService{findings: nil, findingsErr: errors.New("load failed")}
		engine.registerScanner(testEngine, mock)

		results := engine.LoadFindings()

		assert.Empty(t, results)
	})

	t.Run("Should return empty slice when no scanners registered", func(t *testing.T) {
		engine := NewEngine()

		results := engine.LoadFindings()

		assert.Empty(t, results)
	})

	t.Run("Should aggregate findings from multiple scanners", func(t *testing.T) {
		engine := NewEngine()
		mock1 := &mockScanService{findings: []models.Finding{{Engine: "IACST", Name: "Finding 1"}}}
		mock2 := &mockScanService{findings: []models.Finding{{Engine: "IACST", Name: "Finding 2"}}}
		engine.registerScanner("scanner1", mock1)
		engine.registerScanner("scanner2", mock2)

		results := engine.LoadFindings()

		assert.EqualValues(t, 2, len(results))
	})
}

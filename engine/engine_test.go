package engine

import (
	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/domains/models"
	"ScopeGuardian/loader"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testEngine = "test_engine"
)

type mockScanService struct {
	startOk     bool
	startErr    error
	started     bool
	findings    []models.Finding
	findingsErr error
	syncErr     error
}

func (m *mockScanService) Start() (bool, error) {
	m.started = true
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
	assert.EqualValues(t, 0, len(engine.prerequisites))
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

func TestRegisterPrerequisite(t *testing.T) {
	t.Run("Should register prerequisite", func(t *testing.T) {
		engine := NewEngine()
		mock := &mockScanService{startOk: true}

		ok := engine.registerPrerequisite(testEngine, mock)

		assert.True(t, ok)
		assert.EqualValues(t, 1, len(engine.prerequisites))
	})

	t.Run("Should not register two prerequisites under the same name", func(t *testing.T) {
		engine := NewEngine()
		mock := &mockScanService{startOk: true}
		engine.registerPrerequisite(testEngine, mock)

		ok := engine.registerPrerequisite(testEngine, mock)
		assert.False(t, ok)
		assert.EqualValues(t, 1, len(engine.prerequisites))
	})

	t.Run("Should not register a prerequisite with an empty name", func(t *testing.T) {
		engine := NewEngine()
		mock := &mockScanService{startOk: true}

		ok := engine.registerPrerequisite("", mock)
		assert.False(t, ok)
		assert.EqualValues(t, 0, len(engine.prerequisites))
	})
}

func TestRegisterDependentScanner(t *testing.T) {
	t.Run("Should register dependent scanner with DependsOn set", func(t *testing.T) {
		engine := NewEngine()
		mock := &mockScanService{startOk: true}

		ok := engine.registerDependentScanner(testEngine, mock, "some-prerequisite")

		assert.True(t, ok)
		assert.EqualValues(t, 1, len(engine.scanners))
		assert.EqualValues(t, "some-prerequisite", engine.scanners[testEngine].DependsOn)
	})

	t.Run("Should not register two dependent scanners under the same name", func(t *testing.T) {
		engine := NewEngine()
		mock := &mockScanService{startOk: true}
		engine.registerDependentScanner(testEngine, mock, "prereq")

		ok := engine.registerDependentScanner(testEngine, mock, "prereq")
		assert.False(t, ok)
		assert.EqualValues(t, 1, len(engine.scanners))
	})
}

func TestInitialize(t *testing.T) {
	t.Run("Should initialize engine with kics runner", func(t *testing.T) {
		engine := NewEngine()

		config, err := loader.Load("../loader/mocks/config.toml")
		assert.Nil(t, err)
		assert.NotNil(t, config)

		err = engine.Initialize(config)
		assert.Nil(t, err)

		assert.EqualValues(t, 1, len(engine.scanners))
	})

	t.Run("Should initialize engine with syft runner when grype is configured", func(t *testing.T) {
		engine := NewEngine()

		config, err := loader.Load("../loader/mocks/config_with_grype.toml")
		assert.Nil(t, err)
		assert.NotNil(t, config)

		err = engine.Initialize(config)
		assert.Nil(t, err)

		assert.EqualValues(t, 1, len(engine.prerequisites))
		assert.EqualValues(t, 1, len(engine.scanners))
		_, ok := engine.prerequisites[syftScannerName]
		assert.True(t, ok)
		_, ok = engine.scanners[grypeScannerName]
		assert.True(t, ok)
		assert.EqualValues(t, syftScannerName, engine.scanners[grypeScannerName].DependsOn)
	})

	t.Run("Should not initialize engine with kics runner", func(t *testing.T) {
		engine := NewEngine()

		config, err := loader.Load("../loader/mocks/config_no_engine.toml")
		assert.Nil(t, err)
		assert.NotNil(t, config)

		err = engine.Initialize(config)
		assert.Nil(t, err)

		assert.EqualValues(t, 0, len(engine.scanners))
	})

	t.Run("Should initialize engine with opengrep runner", func(t *testing.T) {
		engine := NewEngine()

		config, err := loader.Load("../loader/mocks/config_with_opengrep.toml")
		assert.Nil(t, err)
		assert.NotNil(t, config)

		err = engine.Initialize(config)
		assert.Nil(t, err)

		assert.EqualValues(t, 1, len(engine.scanners))
		_, ok := engine.scanners[opengrepScannerName]
		assert.True(t, ok)
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

	t.Run("Should run prerequisite before dependent scanner", func(t *testing.T) {
		engine := NewEngine()
		prereq := &mockScanService{startOk: true, startErr: nil}
		dependent := &mockScanService{startOk: true, startErr: nil}

		engine.registerPrerequisite("prereq", prereq)
		engine.registerDependentScanner("dependent", dependent, "prereq")

		assert.NotPanics(t, func() {
			engine.Start()
		})

		assert.True(t, prereq.started)
		assert.True(t, dependent.started)
	})

	t.Run("Should skip dependent scanner when prerequisite fails", func(t *testing.T) {
		engine := NewEngine()
		prereq := &mockScanService{startOk: false, startErr: errors.New("prereq failed")}
		dependent := &mockScanService{startOk: true, startErr: nil}

		engine.registerPrerequisite("prereq", prereq)
		engine.registerDependentScanner("dependent", dependent, "prereq")

		assert.NotPanics(t, func() {
			engine.Start()
		})

		assert.True(t, prereq.started)
		assert.False(t, dependent.started)
	})

	t.Run("Should still run independent scanner when prerequisite fails", func(t *testing.T) {
		engine := NewEngine()
		prereq := &mockScanService{startOk: false, startErr: errors.New("prereq failed")}
		dependent := &mockScanService{startOk: true, startErr: nil}
		independent := &mockScanService{startOk: true, startErr: nil}

		engine.registerPrerequisite("prereq", prereq)
		engine.registerDependentScanner("dependent", dependent, "prereq")
		engine.registerScanner("independent", independent)

		assert.NotPanics(t, func() {
			engine.Start()
		})

		assert.True(t, prereq.started)
		assert.False(t, dependent.started)
		assert.True(t, independent.started)
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

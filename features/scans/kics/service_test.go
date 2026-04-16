package kics

import (
	"fmt"
	"os"
	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/interfaces"
	environment_variable "ScopeGuardian/environnement_variable"
	"ScopeGuardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDefectDojoService struct {
	importScanOk  bool
	importScanErr error
}

func (m *mockDefectDojoService) GetProductByName(_ string) (defectdojo.Product, error) {
	return defectdojo.Product{}, nil
}

func (m *mockDefectDojoService) CreateEngagement(_ string, _ string, _ int, _ bool) (int, error) {
	return 0, nil
}

func (m *mockDefectDojoService) GetEngagements(_ uint, _ int, _ int, _ []defectdojo.Engagement) ([]defectdojo.Engagement, error) {
	return []defectdojo.Engagement{}, nil
}

func (m *mockDefectDojoService) UpdateEngagementEndDate(_, _ int, _ bool) (bool, error) {
	return true, nil
}

func (m *mockDefectDojoService) ImportScan(_ defectdojo.ScanPayload, _ string) (bool, error) {
	return m.importScanOk, m.importScanErr
}

func (m *mockDefectDojoService) SetAccessToken(_ string) {}

func (m *mockDefectDojoService) SetURL(_ string) {}

func (m *mockDefectDojoService) GetFindings(_ int, _ int, _ int, _ []defectdojo.Finding) ([]defectdojo.Finding, error) {
	return []defectdojo.Finding{}, nil
}

func (m *mockDefectDojoService) GetAllEngagementFindings(_ int, _ int, _ int, _ []defectdojo.Finding) ([]defectdojo.Finding, error) {
	return []defectdojo.Finding{}, nil
}

var _ defectdojo.DefectDojoService = &mockDefectDojoService{}

// TestNewKicsServiceImplementsInterface verifies that newKicsService returns a
// value that satisfies the interfaces.ScanServiceImpl contract, enforcing that
// all required scanner methods are present at compile and test time.
func TestNewKicsServiceImplementsInterface(t *testing.T) {
	service := newKicsService("./test", loader.Kics{})

	_, ok := service.(interfaces.ScanServiceImpl)
	assert.True(t, ok)
}

func TestNewKicsService(t *testing.T) {
	service := newKicsService("./test", loader.Kics{ExcludeQueries: []string{"a227ec01-f97a-4084-91a4-47b350c1db54"}})

	impl, ok := service.(*KicsServiceImpl)
	assert.True(t, ok)
	assert.EqualValues(t, []string{"a227ec01-f97a-4084-91a4-47b350c1db54"}, impl.excludeQueries)
}

func TestVerifyConfig(t *testing.T) {
	t.Run("Config should be OK", func(t *testing.T) {
		ok, err := verifyConfig(fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks"))
		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Config should not be OK", func(t *testing.T) {
		ok, err := verifyConfig(fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/doesnotexist"))
		assert.NotNil(t, err)
		assert.False(t, ok)
	})
}

func TestLoadFinding(t *testing.T) {
	t.Run("Should load findings", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newKicsService("./test", loader.Kics{})

		findings, err := service.LoadFindings()

		assert.Nil(t, err)
		assert.EqualValues(t, 8, len(findings))
	})

	t.Run("Should not load findings due to lack of results", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newKicsService("./test", loader.Kics{})

		findings, err := service.LoadFindings()

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})

	t.Run("Should not load findings due to unmarshalling issue", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/bad_format_results"))
		environment_variable.ReloadEnv()

		service := newKicsService("./test", loader.Kics{})

		findings, err := service.LoadFindings()

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})
}

// func TestSync(t *testing.T) {
// 	service := newKicsService(loader.Kics{"./test", ""})

// 	err := service.Sync()

// 	assert.Nil(t, err)
// }

func TestSync(t *testing.T) {
	t.Run("Should sync successfully", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newKicsService("./test", loader.Kics{})
		ddMock := &mockDefectDojoService{importScanOk: true, importScanErr: nil}

		err := service.Sync(1, "main", ddMock)

		assert.Nil(t, err)
	})

	t.Run("Should return error when import scan fails", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newKicsService("./test", loader.Kics{})
		ddMock := &mockDefectDojoService{importScanOk: false, importScanErr: fmt.Errorf("import failed")}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})

	t.Run("Should return error when output file not found", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newKicsService("./test", loader.Kics{})
		ddMock := &mockDefectDojoService{}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})
}

package grype

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
	importScanOk      bool
	importScanErr     error
	reimportScanOk    bool
	reimportScanErr   error
	reimportedPayload defectdojo.ScanPayload
	testsToReturn     []defectdojo.Test
	getTestsErr       error
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

func (m *mockDefectDojoService) ReimportScan(payload defectdojo.ScanPayload, _ string) (bool, error) {
	m.reimportedPayload = payload
	return m.reimportScanOk, m.reimportScanErr
}

func (m *mockDefectDojoService) GetTests(_ int, _ string) ([]defectdojo.Test, error) {
	return m.testsToReturn, m.getTestsErr
}

func (m *mockDefectDojoService) SetAccessToken(_ string) {}

func (m *mockDefectDojoService) SetURL(_ string) {}

func (m *mockDefectDojoService) GetFindings(_ int, _ int, _ int, _ []defectdojo.Finding) ([]defectdojo.Finding, error) {
	return []defectdojo.Finding{}, nil
}

var _ defectdojo.DefectDojoService = &mockDefectDojoService{}

func TestNewGrypeService(t *testing.T) {
	service := newGrypeService(loader.Grype{})

	_, ok := service.(interfaces.ScanServiceImpl)
	assert.NotNil(t, service)
	assert.True(t, ok)
}

func TestGrypeStart(t *testing.T) {
	t.Run("Should return error when sbom file not found", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})

		ok, err := service.Start()

		assert.NotNil(t, err)
		assert.False(t, ok)
		assert.EqualValues(t, errSbomNotFound, err.Error())
	})
}

func TestLoadFindings(t *testing.T) {
	t.Run("Should load findings", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})

		findings, err := service.LoadFindings()

		assert.Nil(t, err)
		assert.EqualValues(t, 2, len(findings))

		assert.EqualValues(t, "test-package 1.0.0", findings[0].Name)
		assert.EqualValues(t, "CVE-2021-1234", findings[0].VulnId)
		assert.EqualValues(t, "HIGH", findings[0].Severity)
		assert.EqualValues(t, "Upgrade to 1.2.0", findings[0].Recommendation)

		assert.EqualValues(t, "another-package 2.0.0", findings[1].Name)
		assert.EqualValues(t, "CVE-2021-5678", findings[1].VulnId)
		assert.EqualValues(t, "MEDIUM", findings[1].Severity)
		assert.EqualValues(t, "", findings[1].Recommendation)
	})

	t.Run("Should not load findings due to lack of results", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})

		findings, err := service.LoadFindings()

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})

	t.Run("Should not load findings due to unmarshalling issue", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/bad_format_results"))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})

		findings, err := service.LoadFindings()

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})
}

func TestSync(t *testing.T) {
	t.Run("Should sync successfully using import when no tests exist", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})
		ddMock := &mockDefectDojoService{importScanOk: true, importScanErr: nil}

		err := service.Sync(1, "main", ddMock)

		assert.Nil(t, err)
	})

	t.Run("Should sync successfully using reimport when tests exist", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})
		ddMock := &mockDefectDojoService{
			testsToReturn:  []defectdojo.Test{{Id: 3, ScanType: "Anchore Grype"}},
			reimportScanOk: true,
		}

		err := service.Sync(1, "main", ddMock)

		assert.Nil(t, err)
		assert.EqualValues(t, 3, ddMock.reimportedPayload.TestId)
	})

	t.Run("Should return error when import scan fails", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})
		ddMock := &mockDefectDojoService{importScanOk: false, importScanErr: fmt.Errorf("import failed")}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})

	t.Run("Should return error when reimport scan fails", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})
		ddMock := &mockDefectDojoService{
			testsToReturn:   []defectdojo.Test{{Id: 3, ScanType: "Anchore Grype"}},
			reimportScanOk:  false,
			reimportScanErr: fmt.Errorf("reimport failed"),
		}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})

	t.Run("Should return error when GetTests fails", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})
		ddMock := &mockDefectDojoService{getTestsErr: fmt.Errorf("cannot retrieve tests")}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})

	t.Run("Should return error when output file not found", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newGrypeService(loader.Grype{})
		ddMock := &mockDefectDojoService{}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})
}

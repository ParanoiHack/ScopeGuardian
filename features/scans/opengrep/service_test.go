package opengrep

import (
	"fmt"
	"os"
	"scope-guardian/connectors/defectdojo"
	"scope-guardian/domains/interfaces"
	environment_variable "scope-guardian/environnement_variable"
	"scope-guardian/loader"
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

var _ defectdojo.DefectDojoService = &mockDefectDojoService{}

func TestNewOpenGrepService(t *testing.T) {
	service := newOpenGrepService("./test", loader.Opengrep{})

	_, ok := service.(interfaces.ScanServiceImpl)
	assert.NotNil(t, service)
	assert.True(t, ok)
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

func TestOpenGrepStart(t *testing.T) {
	t.Run("Should return error when directory not found", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./nonexistent", loader.Opengrep{})

		ok, err := service.Start()

		assert.NotNil(t, err)
		assert.False(t, ok)
		assert.EqualValues(t, errDirectoryNotFound, err.Error())
	})
}

func TestOpenGrepLoadFindings(t *testing.T) {
	t.Run("Should load findings", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})

		findings, err := service.LoadFindings()

		assert.Nil(t, err)
		assert.EqualValues(t, 2, len(findings))

		assert.EqualValues(t, "python.lang.security.deserialization.avoided-pickle-usage", findings[0].Name)
		assert.EqualValues(t, "ERROR", findings[0].Severity)
		assert.EqualValues(t, "/app/src/utils.py", findings[0].SinkFile)
		assert.EqualValues(t, 42, findings[0].SinkLine)

		assert.EqualValues(t, "python.lang.security.audit.formatted-sql-query", findings[1].Name)
		assert.EqualValues(t, "WARNING", findings[1].Severity)
	})

	t.Run("Should not load findings due to lack of results", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})

		findings, err := service.LoadFindings()

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})

	t.Run("Should not load findings due to unmarshalling issue", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/bad_format_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})

		findings, err := service.LoadFindings()

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})
}

func TestOpenGrepSync(t *testing.T) {
	t.Run("Should sync successfully", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})
		ddMock := &mockDefectDojoService{importScanOk: true, importScanErr: nil}

		err := service.Sync(1, "main", ddMock)

		assert.Nil(t, err)
	})

	t.Run("Should return error when import scan fails", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})
		ddMock := &mockDefectDojoService{importScanOk: false, importScanErr: fmt.Errorf("import failed")}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})
}

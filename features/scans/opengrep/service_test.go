package opengrep

import (
	"encoding/json"
	"fmt"
	"os"
	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/domains/models"
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

func (m *mockDefectDojoService) GetAllEngagementFindings(_ int, _ int, _ int, _ []defectdojo.Finding) ([]defectdojo.Finding, error) {
	return []defectdojo.Finding{}, nil
}

var _ defectdojo.DefectDojoService = &mockDefectDojoService{}

// TestNewOpenGrepServiceImplementsInterface verifies that newOpenGrepService returns
// a value that satisfies the interfaces.ScanServiceImpl contract, enforcing that
// all required scanner methods are present at compile and test time.
func TestNewOpenGrepServiceImplementsInterface(t *testing.T) {
	service := newOpenGrepService("./test", loader.Opengrep{})

	_, ok := service.(interfaces.ScanServiceImpl)
	assert.True(t, ok)
}

func TestNewOpenGrepService(t *testing.T) {
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
		assert.EqualValues(t, "HIGH", findings[0].Severity)
		assert.EqualValues(t, "CWE-502: Deserialization of Untrusted Data", findings[0].Cwe)
		assert.EqualValues(t, "A8:2017-Insecure Deserialization", findings[0].Description)
		assert.EqualValues(t, "Deserialization with `pickle` is insecure; it can lead to arbitrary code execution.", findings[0].Recommendation)
		assert.EqualValues(t, "/app/src/utils.py", findings[0].SinkFile)
		assert.EqualValues(t, 42, findings[0].SinkLine)

		assert.EqualValues(t, "python.lang.security.audit.formatted-sql-query", findings[1].Name)
		assert.EqualValues(t, "MEDIUM", findings[1].Severity)
		assert.EqualValues(t, "CWE-89: Improper Neutralization of Special Elements used in an SQL Command", findings[1].Cwe)
		assert.EqualValues(t, "A1:2017-Injection", findings[1].Description)
		assert.EqualValues(t, "Detected possible formatted SQL query. Use parameterized queries instead.", findings[1].Recommendation)
	})

	t.Run("Should load findings when cwe and owasp are plain strings", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/string_metadata_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})

		findings, err := service.LoadFindings()

		assert.Nil(t, err)
		assert.EqualValues(t, 1, len(findings))
		assert.EqualValues(t, "CWE-502: Deserialization of Untrusted Data", findings[0].Cwe)
		assert.EqualValues(t, "A8:2017-Insecure Deserialization", findings[0].Description)
		assert.EqualValues(t, "HIGH", findings[0].Severity)
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

func TestEnrichOpenGrepResults(t *testing.T) {
	t.Run("Should add severity from impact when missing", func(t *testing.T) {
		input := []byte(`{"results":[{"extra":{"message":"msg","metadata":{"impact":"HIGH","confidence":"HIGH"}}}]}`)
		output := enrichOpenGrepResults(input)

		var wrapper map[string]interface{}
		err := json.Unmarshal(output, &wrapper)
		assert.Nil(t, err)

		results := wrapper["results"].([]interface{})
		extra := results[0].(map[string]interface{})["extra"].(map[string]interface{})
		assert.EqualValues(t, "HIGH", extra["severity"])
	})

	t.Run("Should not overwrite existing severity", func(t *testing.T) {
		input := []byte(`{"results":[{"extra":{"severity":"LOW","message":"msg","metadata":{"impact":"HIGH"}}}]}`)
		output := enrichOpenGrepResults(input)

		var wrapper map[string]interface{}
		err := json.Unmarshal(output, &wrapper)
		assert.Nil(t, err)

		results := wrapper["results"].([]interface{})
		extra := results[0].(map[string]interface{})["extra"].(map[string]interface{})
		assert.EqualValues(t, "LOW", extra["severity"])
	})

	t.Run("Should return original bytes on invalid JSON", func(t *testing.T) {
		input := []byte(`not json`)
		output := enrichOpenGrepResults(input)
		assert.Equal(t, input, output)
	})

	t.Run("Should preserve all existing fields", func(t *testing.T) {
		input := []byte(`{"results":[{"check_id":"rule","path":"/a.py","extra":{"message":"msg","metadata":{"impact":"MEDIUM"}}}],"errors":[]}`)
		output := enrichOpenGrepResults(input)

		var wrapper map[string]interface{}
		err := json.Unmarshal(output, &wrapper)
		assert.Nil(t, err)
		assert.NotNil(t, wrapper["errors"])

		results := wrapper["results"].([]interface{})
		result := results[0].(map[string]interface{})
		assert.EqualValues(t, "rule", result["check_id"])
		assert.EqualValues(t, "/a.py", result["path"])

		extra := result["extra"].(map[string]interface{})
		assert.EqualValues(t, "MEDIUM", extra["severity"])
	})

	t.Run("Should inject fingerprint hash into extra.fingerprint", func(t *testing.T) {
		// enrichOpenGrepResults must inject the content hash into extra.fingerprint so
		// that DefectDojo's Semgrep parser stores it as unique_id_from_tool. The hash
		// must equal models.ComputeFindingHash(severity, path, line, "") — the same
		// formula used by LoadFindings — so both sides produce identical values.
		input := []byte(`{"results":[{"check_id":"go.lang.sql","path":"/app/db.go","start":{"line":10,"col":1},"extra":{"message":"msg","metadata":{"impact":"HIGH"}}}]}`)
		output := enrichOpenGrepResults(input)

		var wrapper map[string]interface{}
		err := json.Unmarshal(output, &wrapper)
		assert.Nil(t, err)

		results := wrapper["results"].([]interface{})
		extra := results[0].(map[string]interface{})["extra"].(map[string]interface{})

		wantHash := models.ComputeFindingHash("HIGH", "/app/db.go", 10, "")
		assert.EqualValues(t, wantHash, extra["fingerprint"])
	})

	t.Run("Should produce matching fingerprint for two findings with same rule but different locations", func(t *testing.T) {
		// Two findings with the same check_id but different file/line must get different
		// fingerprints — this is the core correctness property of the hash approach.
		input := []byte(`{"results":[` +
			`{"check_id":"rule.x","path":"/a.py","start":{"line":1},"extra":{"metadata":{"impact":"HIGH"}}},` +
			`{"check_id":"rule.x","path":"/b.py","start":{"line":2},"extra":{"metadata":{"impact":"HIGH"}}}` +
			`]}`)
		output := enrichOpenGrepResults(input)

		var wrapper map[string]interface{}
		err := json.Unmarshal(output, &wrapper)
		assert.Nil(t, err)

		results := wrapper["results"].([]interface{})
		fp1 := results[0].(map[string]interface{})["extra"].(map[string]interface{})["fingerprint"].(string)
		fp2 := results[1].(map[string]interface{})["extra"].(map[string]interface{})["fingerprint"].(string)
		assert.NotEqual(t, fp1, fp2)
	})
}

func TestOpenGrepSync(t *testing.T) {
	t.Run("Should sync successfully using import when no tests exist", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})
		ddMock := &mockDefectDojoService{importScanOk: true, importScanErr: nil}

		err := service.Sync(1, "main", ddMock)

		assert.Nil(t, err)
	})

	t.Run("Should sync successfully using reimport when tests exist", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})
		ddMock := &mockDefectDojoService{
			testsToReturn:  []defectdojo.Test{{Id: 7, ScanType: "Semgrep JSON Report"}},
			reimportScanOk: true,
		}

		err := service.Sync(1, "main", ddMock)

		assert.Nil(t, err)
		assert.EqualValues(t, 7, ddMock.reimportedPayload.TestId)
	})

	t.Run("Should return error when import scan fails", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})
		ddMock := &mockDefectDojoService{importScanOk: false, importScanErr: fmt.Errorf("import failed")}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})

	t.Run("Should return error when reimport scan fails", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})
		ddMock := &mockDefectDojoService{
			testsToReturn:   []defectdojo.Test{{Id: 7, ScanType: "Semgrep JSON Report"}},
			reimportScanOk:  false,
			reimportScanErr: fmt.Errorf("reimport failed"),
		}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})

	t.Run("Should return error when GetTests fails", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], "./mocks/working_results"))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})
		ddMock := &mockDefectDojoService{getTestsErr: fmt.Errorf("cannot retrieve tests")}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})

	t.Run("Should return error when output file not found", func(t *testing.T) {
		_ = os.Setenv("SCAN_DIR", fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["PWD"], ""))
		environment_variable.ReloadEnv()

		service := newOpenGrepService("./test", loader.Opengrep{})
		ddMock := &mockDefectDojoService{}

		err := service.Sync(1, "main", ddMock)

		assert.NotNil(t, err)
	})
}

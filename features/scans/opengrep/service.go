package opengrep

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/domains/models"
	environment_variable "ScopeGuardian/environnement_variable"
	"ScopeGuardian/exec"
	"ScopeGuardian/loader"
	"ScopeGuardian/logger"
	"strings"
	"time"
)

// OpenGrepServiceImpl implements ScanServiceImpl for the OpenGrep SAST scanner.
type OpenGrepServiceImpl struct {
	path        string
	output      string
	exclude     []string
	excludeRule []string
	proxyEnv    []string
}

// newOpenGrepService builds an OpenGrepServiceImpl from the scan path and loader configuration,
// resolving the scan path and output file path relative to the SCAN_DIR environment variable.
// proxyEnv is an optional list of "KEY=VALUE" proxy environment variable entries
// (see loader.Proxy.ToEnv) forwarded to the OpenGrep process.
func newOpenGrepService(path string, config loader.Opengrep, proxyEnv []string) interfaces.ScanServiceImpl {
	return &OpenGrepServiceImpl{
		path:        fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], path),
		output:      fmt.Sprintf("%s/%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], outputFolder, outputNameParameter),
		exclude:     config.Exclude,
		excludeRule: config.ExcludeRule,
		proxyEnv:    proxyEnv,
	}
}

// verifyConfig checks that the directory at path exists and is accessible.
// Returns false and an error if the path cannot be stat'd.
func verifyConfig(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		logger.Error(fmt.Sprintf(logErrorDirectoryNotFound, path))
		return false, errors.New(errDirectoryNotFound)
	}

	return true, nil
}

// Start validates the scan directory and then invokes the OpenGrep binary with the
// appropriate arguments. It returns true on success or false and an error if the
// directory is missing or the OpenGrep process exits with a non-zero status.
func (s *OpenGrepServiceImpl) Start() (bool, error) {
	if ok, err := verifyConfig(s.path); !ok && err != nil {
		return ok, err
	}

	if err := os.MkdirAll(filepath.Dir(s.output), 0755); err != nil {
		return false, err
	}

	args := []string{
		fmt.Sprintf("%s%s", jsonOutputArgument, s.output),
		ossOnlyArgument,
		quietArgument,
		skipUnknownExtArgument,
	}

	for _, ex := range s.exclude {
		args = append(args, excludeArgument, ex)
	}

	for _, rule := range s.excludeRule {
		args = append(args, excludeRuleArgument, rule)
	}

	args = append(args, s.path)

	logger.Info(fmt.Sprintf(logInfoCommandLine, strings.Join(args, " ")))

	return exec.Wrap(binaryPath, dirPath, args, io.Discard, io.Discard, s.proxyEnv...)
}

// LoadFindings reads the OpenGrep JSON output file and converts each result
// into a slice of domain Finding objects. Returns an error if the file cannot
// be read or parsed.
func (s *OpenGrepServiceImpl) LoadFindings() ([]models.Finding, error) {
	buffer, err := os.ReadFile(s.output)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorFileNotFound, s.output))
		return nil, err
	}

	var results OpenGrepResults
	if err := json.Unmarshal(buffer, &results); err != nil {
		logger.Error(logErrorParseResults)
		return nil, err
	}

	var findings []models.Finding
	for _, item := range results.Results {
		severity := strings.ToUpper(item.Extra.Metadata.Impact)

		cwe := ""
		if len(item.Extra.Metadata.Cwe) > 0 {
			cwe = item.Extra.Metadata.Cwe[0]
		}

		description := strings.Join(item.Extra.Metadata.Owasp, ", ")

		f := models.Finding{
			Engine:         scannerType,
			Severity:       severity,
			Name:           item.CheckId,
			VulnId:         item.CheckId,
			Cwe:            cwe,
			Description:    description,
			// Message is kept for display but intentionally excluded from the hash:
			// DefectDojo's Semgrep parser stores extra.message in description, not
			// mitigation, so f.Mitigation from the DD API will be empty for these findings.
			Recommendation: item.Extra.Message,
			SinkFile:       item.Path,
			SinkLine:       item.Start.Line,
		}
		f.Hash = models.ComputeFindingHash(f.Severity, f.SinkFile, f.SinkLine, "")
		findings = append(findings, f)
	}

	return findings, nil
}

// enrichOpenGrepResults post-processes the OpenGrep JSON output before it is uploaded
// to DefectDojo. It performs two enrichments per result entry:
//
//  1. extra.severity — copied from extra.metadata.impact when absent. DefectDojo's
//     Semgrep JSON Report parser requires this field to determine the finding severity.
//
//  2. extra.fingerprint — set to the content hash computed by models.ComputeFindingHash.
//     DefectDojo's Semgrep parser maps extra.fingerprint to unique_id_from_tool, which
//     is returned by the findings API. This lets FilterByActiveFindings match local
//     findings directly via UniqueIdFromTool without fragile multi-field heuristics and
//     without collision risk from multiple findings sharing the same rule ID (check_id).
//
// The original bytes are returned unchanged if the JSON cannot be parsed.
func enrichOpenGrepResults(data []byte) []byte {
	var wrapper map[string]interface{}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return data
	}

	results, ok := wrapper["results"].([]interface{})
	if !ok {
		return data
	}

	for _, r := range results {
		result, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		extra, ok := result["extra"].(map[string]interface{})
		if !ok {
			continue
		}
		metadata, ok := extra["metadata"].(map[string]interface{})
		if !ok {
			continue
		}

		// 1. Inject extra.severity from extra.metadata.impact when absent.
		if _, hasSeverity := extra["severity"]; !hasSeverity {
			if impact, ok := metadata["impact"].(string); ok && impact != "" {
				extra["severity"] = strings.ToUpper(impact)
			}
		}

		// 2. Inject a pre-computed content hash into extra.fingerprint. DefectDojo's
		// Semgrep parser maps extra.fingerprint to unique_id_from_tool, which is
		// returned by the findings API. FilterByActiveFindings can then match local
		// findings directly via that field using the same hash formula, giving a
		// per-finding identifier that is stable even when multiple findings share
		// the same rule ID (check_id).
		fingerSeverity := ""
		if impact, ok := metadata["impact"].(string); ok {
			fingerSeverity = strings.ToUpper(impact)
		}
		fingerPath, _ := result["path"].(string)
		var fingerLine int
		if start, ok := result["start"].(map[string]interface{}); ok {
			if lineFloat, ok := start["line"].(float64); ok {
				fingerLine = int(lineFloat)
			}
		}
		extra["fingerprint"] = models.ComputeFindingHash(fingerSeverity, fingerPath, fingerLine, "")
	}

	enriched, err := json.Marshal(wrapper)
	if err != nil {
		return data
	}
	return enriched
}

// Sync uploads the OpenGrep scan output to DefectDojo via the given service.
// It constructs a ScanPayload from the stored output file and the provided
// engagement ID and branch, then calls ReimportScan if a test with this scan
// type already exists for the engagement, or ImportScan otherwise.
func (s *OpenGrepServiceImpl) Sync(engagementId int, branch string, service defectdojo.DefectDojoService) error {
	var payload defectdojo.ScanPayload

	payload.Timestamp = time.Now().Format("2006-01-02")
	payload.SeverityThreshold = severityThreshold
	payload.Branch = branch
	payload.Tags = []string{SASTEngineTag}
	payload.GroupBy = groupByProperty
	payload.FindingGroup = findingGroupProperty
	payload.FindingTag = findingTagProperty
	payload.ScanType = scanType
	payload.EngagementId = engagementId
	payload.CloseOldFinding = closeOldFinding
	payload.DoNotReactivate = doNotReactivate

	fileContent, err := os.ReadFile(s.output)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorFileNotFound, s.output))
		return err
	}
	payload.File = enrichOpenGrepResults(fileContent)

	tests, err := service.GetTests(engagementId, scanType)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorGetTests, engagementId))
		return err
	}

	if len(tests) > 0 {
		payload.TestId = tests[0].Id
		if ok, err := service.ReimportScan(payload, s.output); !ok || err != nil {
			logger.Error(fmt.Sprintf(logErrorReimportScan, engagementId))
			return err
		}
	} else {
		if ok, err := service.ImportScan(payload, s.output); !ok || err != nil {
			logger.Error(fmt.Sprintf(logErrorImportScan, engagementId))
			return err
		}
	}

	return nil
}

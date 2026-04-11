package opengrep

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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
}

// newOpenGrepService builds an OpenGrepServiceImpl from the scan path and loader configuration,
// resolving the scan path and output file path relative to the SCAN_DIR environment variable.
func newOpenGrepService(path string, config loader.Opengrep) interfaces.ScanServiceImpl {
	return &OpenGrepServiceImpl{
		path:        fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], path),
		output:      fmt.Sprintf("%s/%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], outputFolder, outputNameParameter),
		exclude:     config.Exclude,
		excludeRule: config.ExcludeRule,
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

	return exec.Wrap(binaryPath, dirPath, args, io.Discard, io.Discard)
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

		findings = append(findings, models.Finding{
			Engine:         scannerType,
			Severity:       severity,
			Name:           item.CheckId,
			Cwe:            cwe,
			Description:    description,
			Recommendation: item.Extra.Message,
			SinkFile:       item.Path,
			SinkLine:       item.Start.Line,
		})
	}

	return findings, nil
}

// enrichOpenGrepResults ensures each result in the OpenGrep JSON output carries an
// extra.severity field, which DefectDojo's Semgrep JSON Report parser requires.
// OpenGrep stores severity as extra.metadata.impact; this function copies that value
// into extra.severity when the field is absent. The original bytes are returned
// unchanged if the JSON cannot be parsed.
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
		if _, hasSeverity := extra["severity"]; hasSeverity {
			continue
		}
		metadata, ok := extra["metadata"].(map[string]interface{})
		if !ok {
			continue
		}
		impact, ok := metadata["impact"].(string)
		if !ok || impact == "" {
			continue
		}
		extra["severity"] = strings.ToUpper(impact)
	}

	enriched, err := json.Marshal(wrapper)
	if err != nil {
		return data
	}
	return enriched
}

// Sync uploads the OpenGrep scan output to DefectDojo via the given service.
// It constructs a ScanPayload from the stored output file and the provided
// engagement ID and branch, then calls ImportScan.
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

	fileContent, err := os.ReadFile(s.output)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorFileNotFound, s.output))
		return err
	}
	payload.File = enrichOpenGrepResults(fileContent)

	if ok, err := service.ImportScan(payload, s.output); !ok || err != nil {
		logger.Error(fmt.Sprintf(logErrorImportScan, engagementId))
		return err
	}

	return nil
}

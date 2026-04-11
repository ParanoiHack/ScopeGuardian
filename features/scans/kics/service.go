package kics

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

// KicsServiceImpl implements ScanServiceImpl for the KICS infrastructure-as-code scanner.
type KicsServiceImpl struct {
	path           string
	platform       string
	output         string
	excludeQueries []string
}

// newKicsService builds a KicsServiceImpl from the scan path and loader configuration,
// resolving the scan path and output file path relative to the SCAN_DIR environment variable.
func newKicsService(path string, config loader.Kics) interfaces.ScanServiceImpl {
	return &KicsServiceImpl{
		path:           fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], path),
		output:         fmt.Sprintf("%s/%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], outputFolder, outputNameParameter),
		platform:       config.Platform,
		excludeQueries: config.ExcludeQueries,
	}
}

// verifyConfig checks that the directory at path exists and is accessible.
// Returns false and an error if the path cannot be stat'd.
func verifyConfig(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		logger.Error(fmt.Sprintf(logErrorDirectorNotFound, path))
		return false, errors.New(errDirectoryNotFound)
	}

	return true, nil
}

// Start validates the scan directory and then invokes the KICS binary with the
// appropriate arguments. It returns true on success or false and an error if the
// directory is missing or the KICS process exits with a non-zero status.
func (s *KicsServiceImpl) Start() (bool, error) {
	if ok, err := verifyConfig(s.path); !ok && err != nil {
		return ok, err
	}

	args := []string{scanArgument}

	args = append(args, []string{
		pathArgument,
		s.path,
		ciArgument,
		librariesPathArgument,
		librariesPathParameter,
		queriesPathArgument,
		queriesPathParameter,
		outputPathArgument,
		fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], outputFolder),
		outputNameArgument,
		outputNameParameter,
		ignoreOnExitArgument,
		ignoreOnExitParameter,
	}...)

	if s.platform != "" {
		args = append(args, []string{typeArgument, s.platform}...)
	}

	for _, q := range s.excludeQueries {
		args = append(args, excludeQueriesArgument, q)
	}

	logger.Info(fmt.Sprintf(logInfoCommandLine, strings.Join(args, " ")))

	return exec.Wrap(binaryPath, dirPath, args, io.Discard, io.Discard)
}

// LoadFindings reads the KICS JSON output file and converts each query result
// into a slice of domain Finding objects. Returns an error if the file cannot
// be read or parsed.
func (s *KicsServiceImpl) LoadFindings() ([]models.Finding, error) {
	buffer, err := os.ReadFile(s.output)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorFileNotFound, s.output))
		return nil, err
	}

	var results KicsResults
	if err := json.Unmarshal(buffer, &results); err != nil {
		logger.Error(logErrorParseResults)
		return nil, err
	}

	var findings []models.Finding
	for _, item := range results.Queries {
		for _, sink := range item.Files {
			findings = append(findings, models.Finding{
				Engine:         scannerType,
				Severity:       item.Severity,
				Name:           item.QueryName,
				Cwe:            item.Cwe,
				Description:    item.Description,
				SinkFile:       sink.FileName,
				SinkLine:       sink.Line,
				Recommendation: sink.Recommendation,
			})
		}

	}

	return findings, nil
}

// Sync uploads the KICS scan output to DefectDojo via the given service.
// It constructs a ScanPayload from the stored output file and the provided
// engagement ID and branch, then calls ImportScan.
func (s *KicsServiceImpl) Sync(engagementId int, branch string, service defectdojo.DefectDojoService) error {
	var payload defectdojo.ScanPayload

	payload.Timestamp = time.Now().Format("2006-01-02")
	payload.SeverityThreshold = severityThreshold
	payload.Branch = branch
	payload.Tags = []string{IACSTEngineTag}
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
	payload.File = fileContent

	if ok, err := service.ImportScan(payload, s.output); !ok || err != nil {
		logger.Error(fmt.Sprintf(logErrorImportScan, engagementId))
		return err
	}

	return nil
}

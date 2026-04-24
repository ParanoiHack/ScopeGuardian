package grype

import (
	"encoding/json"
	"errors"
	"fmt"
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

// GrypeServiceImpl implements ScanServiceImpl for the Grype vulnerability scanner.
// It scans the Syft-generated SBOM and reports software composition analysis findings.
type GrypeServiceImpl struct {
	sbom         string
	output       string
	ignoreStates string
	exclude      []string
}

// newGrypeService builds a GrypeServiceImpl from the Grype loader configuration,
// resolving the SBOM input and result output paths from the SCAN_DIR environment variable.
func newGrypeService(config loader.Grype) interfaces.ScanServiceImpl {
	return &GrypeServiceImpl{
		sbom:         fmt.Sprintf("%s/%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], outputFolder, sbomInputNameParameter),
		output:       fmt.Sprintf("%s/%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], outputFolder, outputNameParameter),
		ignoreStates: config.IgnoreStates,
		exclude:      config.Exclude,
	}
}

// Start verifies that the Syft SBOM file exists and then invokes the Grype binary
// to scan it for known vulnerabilities. It returns true on success or false and an
// error if the SBOM is missing or the Grype process exits with a non-zero status.
func (s *GrypeServiceImpl) Start() (bool, error) {
	if err := os.MkdirAll(filepath.Dir(s.output), 0755); err != nil {
		return false, err
	}

	if _, err := os.Stat(s.sbom); err != nil {
		logger.Error(fmt.Sprintf(logErrorSbomNotFound, s.sbom))
		return false, errors.New(errSbomNotFound)
	}

	args := []string{
		fmt.Sprintf("%s%s", sbomScheme, s.sbom),
	}

	if s.ignoreStates != "" {
		args = append(args, ignoreStatesArgument, s.ignoreStates)
	}

	args = append(args, outputArgument, outputFormatParameter)
	args = append(args, quietArgument)
	args = append(args, fileArgument, s.output)
	args = append(args, configArgument, configPath)

	for _, ex := range s.exclude {
		args = append(args, excludeArgument, ex)
	}

	logger.Info(fmt.Sprintf(logInfoCommandLine, strings.Join(args, " ")))

	return exec.WrapAllowExitCodes(binaryPath, dirPath, args, os.Stdout, os.Stderr, []int{exitCodeFindings})
}

// LoadFindings reads the Grype JSON output file and converts each vulnerability
// match into a slice of domain Finding objects. Returns an error if the file
// cannot be read or parsed.
func (s *GrypeServiceImpl) LoadFindings() ([]models.Finding, error) {
	buffer, err := os.ReadFile(s.output)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorFileNotFound, s.output))
		return nil, err
	}

	var results GrypeResults
	if err := json.Unmarshal(buffer, &results); err != nil {
		logger.Error(logErrorParseResults)
		return nil, err
	}

	var findings []models.Finding
	for _, match := range results.Matches {
		sinkFile := ""
		if len(match.Artifact.Locations) > 0 {
			sinkFile = match.Artifact.Locations[0].Path
		}

		recommendation := ""
		if len(match.Vulnerability.Fix.Versions) == 1 {
			recommendation = fmt.Sprintf(recommendationUpgrade, match.Vulnerability.Fix.Versions[0])
		} else if len(match.Vulnerability.Fix.Versions) > 1 {
			recommendation = fmt.Sprintf(recommendationUpgradeMultiple, strings.Join(match.Vulnerability.Fix.Versions, ", "))
		}

		severity := strings.ToUpper(match.Vulnerability.Severity)
		f := models.Finding{
			Engine:         scannerType,
			Severity:       severity,
			Name:           fmt.Sprintf("%s %s", match.Artifact.Name, match.Artifact.Version),
			VulnId:         match.Vulnerability.ID,
			Description:    match.Vulnerability.Description,
			SinkFile:       sinkFile,
			Recommendation: recommendation,
		}
		f.Hash = models.ComputeFindingHash(f.Severity, f.SinkFile, f.SinkLine, f.Recommendation)
		findings = append(findings, f)
	}

	return findings, nil
}

// Sync uploads the Grype scan output to DefectDojo via the given service.
// It constructs a ScanPayload from the stored output file and the provided
// engagement ID and branch, then calls ReimportScan if a test with this scan
// type already exists for the engagement, or ImportScan otherwise.
func (s *GrypeServiceImpl) Sync(engagementId int, branch string, service defectdojo.DefectDojoService) error {
	var payload defectdojo.ScanPayload

	payload.Timestamp = time.Now().Format("2006-01-02")
	payload.SeverityThreshold = severityThreshold
	payload.Branch = branch
	payload.Tags = []string{SCAEngineTag}
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
	payload.File = fileContent

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

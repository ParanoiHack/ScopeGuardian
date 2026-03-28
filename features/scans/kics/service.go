package kics

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"scope-guardian/connectors/defectdojo"
	"scope-guardian/domains/interfaces"
	"scope-guardian/domains/models"
	environment_variable "scope-guardian/environnement_variable"
	"scope-guardian/exec"
	"scope-guardian/loader"
	"scope-guardian/logger"
	"strings"
	"time"
)

type KicsServiceImpl struct {
	path     string
	platform string
	output   string
}

func newKicsService(config loader.Kics) interfaces.ScanServiceImpl {
	return &KicsServiceImpl{
		path:     fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], config.Path),
		output:   fmt.Sprintf("%s/%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], outputFolder, outputNameParameter),
		platform: config.Platform,
	}
}

func verifyConfig(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		logger.Error(fmt.Sprintf(logErrorDirectorNotFound, path))
		return false, errors.New(errDirectoryNotFound)
	}

	return true, nil
}

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

	logger.Info(fmt.Sprintf(logInfoCommandLine, strings.Join(args, " ")))

	return exec.Wrap(binaryPath, dirPath, args)
}

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
	payload.File, _ = os.ReadFile(s.output)

	if ok, err := service.ImportScan(payload, s.output); !ok || err != nil {
		logger.Error(fmt.Sprintf(logErrorImportScan, engagementId))
		return err
	}

	return nil
}

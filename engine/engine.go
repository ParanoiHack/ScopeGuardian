package engine

import (
	"fmt"
	"net/http"
	"reflect"
	"scope-guardian/connectors/defectdojo"
	"scope-guardian/connectors/defectdojo/client"
	"scope-guardian/domains/interfaces"
	"scope-guardian/domains/models"
	environment_variable "scope-guardian/environnement_variable"
	featuresync "scope-guardian/features/sync"
	"scope-guardian/features/scans/kics"
	"scope-guardian/loader"
	"scope-guardian/logger"
	"sync"
)

type Engine struct {
	scanners map[string]Scanner
}

type Scanner struct {
	Service interfaces.ScanServiceImpl
}

func NewEngine() *Engine {
	return &Engine{
		scanners: make(map[string]Scanner),
	}
}

func (e *Engine) Initialize(config loader.Config) {
	if !reflect.DeepEqual(config.Kics, loader.Kics{}) {
		logger.Info(logInfoKicsRegister)
		e.registerScanner(kicsScannerName, kics.GetKicsService(config.Kics))
	}
}

func (e *Engine) Start() {
	var wg sync.WaitGroup

	for k, scanner := range e.scanners {
		wg.Add(1)
		go func(scannerName string, service interfaces.ScanServiceImpl) {
			defer wg.Done()

			logger.Info(fmt.Sprintf(logInfoScannerStarting, scannerName))
			if ok, err := service.Start(); !ok || err != nil {
				logger.Error(err.Error())
				logger.Error(fmt.Sprintf(logErrorScannerFailed, scannerName))
			} else {
				logger.Info(fmt.Sprintf(logInfoScannerSuccess, scannerName))
			}
		}(k, scanner.Service)
	}

	wg.Wait()
}

func (e *Engine) LoadFindings() []models.Finding {
	var results []models.Finding

	for k, scanner := range e.scanners {
		findings, err := scanner.Service.LoadFindings()
		if err != nil {
			logger.Error(err.Error())
			logger.Error(fmt.Sprintf(logErrorLoadFinding, k))
		} else {
			logger.Info(fmt.Sprintf(logInfoFindingsLoaded, k))
			results = append(results, findings...)
		}
	}

	return results
}

func (e *Engine) SyncResults(projectName string, branch string) {
	ddService := defectdojo.GetDefectDojoService(
		client.NewClient(&http.Client{}),
		environment_variable.EnvironmentVariable["DD_URL"],
		environment_variable.EnvironmentVariable["DD_ACCESS_TOKEN"])

	engagementId, err := featuresync.GetEngagementId(ddService, projectName, branch)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorRetrieveEngagementId, projectName, branch))
		return
	}

	for k, scanner := range e.scanners {
		scanner.Service.Sync(engagementId, branch, ddService) // check
		_, _ = k, scanner
		logger.Info(fmt.Sprintf(logInfoSyncResult, k))
	}
}

func (e *Engine) registerScanner(name string, service interfaces.ScanServiceImpl) bool {
	if _, ok := e.scanners[name]; ok || name == "" {
		logger.Error(fmt.Sprintf(logErrorRegisterScanner, name))
		return false
	}

	e.scanners[name] = Scanner{
		Service: service,
	}

	return true
}

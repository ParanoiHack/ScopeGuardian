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

// Engine orchestrates one or more security scanners: it registers them,
// runs them in parallel, aggregates their findings, and optionally syncs
// results to DefectDojo.
type Engine struct {
	scanners map[string]Scanner
}

// Scanner wraps a ScanServiceImpl so it can be stored in the engine's registry.
type Scanner struct {
	Service interfaces.ScanServiceImpl
}

// NewEngine allocates and returns an empty Engine ready to accept scanners.
func NewEngine() *Engine {
	return &Engine{
		scanners: make(map[string]Scanner),
	}
}

// Initialize reads the provided configuration and registers any scanner whose
// section is present and non-empty. Currently supports KICS.
func (e *Engine) Initialize(config loader.Config) {
	if !reflect.DeepEqual(config.Kics, loader.Kics{}) {
		logger.Info(logInfoKicsRegister)
		e.registerScanner(kicsScannerName, kics.GetKicsService(config.Kics))
	}
}

// Start launches all registered scanners concurrently and waits for them to finish.
// Errors from individual scanners are logged but do not stop the other scanners.
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

// LoadFindings collects and merges the findings from all registered scanners.
// Errors from individual scanners are logged; successfully loaded findings are
// still included in the returned slice.
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

// SyncResults uploads each scanner's findings to DefectDojo under the engagement
// matching the given projectName and branch. If no matching engagement exists one
// is created automatically. Errors are logged and the sync is skipped for that scanner.
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

// registerScanner adds a scanner under name in the engine's registry.
// It returns false (and logs an error) if the name is empty or already registered.
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

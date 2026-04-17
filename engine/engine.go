package engine

import (
	"fmt"
	"net/http"
	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/connectors/defectdojo/client"
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/domains/models"
	environment_variable "ScopeGuardian/environnement_variable"
	featuresync "ScopeGuardian/features/sync"
	"ScopeGuardian/features/scans/grype"
	"ScopeGuardian/features/scans/kics"
	"ScopeGuardian/features/scans/opengrep"
	"ScopeGuardian/features/scans/syft"
	"ScopeGuardian/loader"
	"ScopeGuardian/logger"
	"sync"
)

// Engine orchestrates one or more security scanners: it registers them,
// runs them in parallel, aggregates their findings, and optionally syncs
// results to DefectDojo.
type Engine struct {
	prerequisites map[string]Scanner
	scanners      map[string]Scanner
}

// Scanner wraps a ScanServiceImpl so it can be stored in the engine's registry.
// DependsOn is the name of a prerequisite scanner that must succeed before this
// scanner is allowed to start (empty means no dependency).
type Scanner struct {
	Service   interfaces.ScanServiceImpl
	DependsOn string
}

// NewEngine allocates and returns an empty Engine ready to accept scanners.
func NewEngine() *Engine {
	return &Engine{
		prerequisites: make(map[string]Scanner),
		scanners:      make(map[string]Scanner),
	}
}

// Initialize reads the provided configuration and registers any scanner whose
// section is present and non-empty. Syft is registered as a prerequisite for
// Grype so that it always runs and completes before Grype starts.
func (e *Engine) Initialize(config loader.Config) {
	if config.Kics != nil {
		logger.Info(logInfoKicsRegister)
		e.registerScanner(kicsScannerName, kics.GetKicsService(config))
	}

	if config.Grype != nil {
		logger.Info(logInfoSyftRegister)
		e.registerPrerequisite(syftScannerName, syft.GetSyftService(config))
		logger.Info(logInfoGrypeRegister)
		e.registerDependentScanner(grypeScannerName, grype.GetGrypeService(config), syftScannerName)
	}

	if config.Opengrep != nil {
		logger.Info(logInfoOpenGrepRegister)
		e.registerScanner(opengrepScannerName, opengrep.GetOpenGrepService(config))
	}
}

// Start runs all registered scanners in two phases:
//  1. Prerequisites (e.g. Syft) are executed concurrently and the engine waits
//     for all of them to finish. The name of each failed prerequisite is recorded.
//  2. Regular scanners (e.g. Grype, KICS) are then executed concurrently. Any
//     scanner whose DependsOn prerequisite failed is skipped with a log message.
//
// Errors from individual scanners are logged but do not stop the other scanners.
func (e *Engine) Start() {
	var wg sync.WaitGroup
	failedPrereqs := make(map[string]bool)
	var mu sync.Mutex

	// Phase 1: run prerequisites concurrently and wait.
	for k, scanner := range e.prerequisites {
		wg.Add(1)
		go func(scannerName string, service interfaces.ScanServiceImpl) {
			defer wg.Done()

			logger.Info(fmt.Sprintf(logInfoScannerStarting, scannerName))
			if ok, err := service.Start(); !ok || err != nil {
				logger.Error(err.Error())
				logger.Error(fmt.Sprintf(logErrorScannerFailed, scannerName))
				mu.Lock()
				failedPrereqs[scannerName] = true
				mu.Unlock()
			} else {
				logger.Info(fmt.Sprintf(logInfoScannerSuccess, scannerName))
			}
		}(k, scanner.Service)
	}

	wg.Wait()

	// Phase 2: run scanners concurrently, skipping those whose prerequisite failed.
	for k, scanner := range e.scanners {
		if scanner.DependsOn != "" && failedPrereqs[scanner.DependsOn] {
			logger.Error(fmt.Sprintf(logErrorSkippingScanner, k, scanner.DependsOn))
			continue
		}

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
// Prerequisites (e.g. Syft) do not contribute findings and are not iterated here.
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
// is created automatically. protectedBranches determines the engagement end date duration.
// Prerequisites (e.g. Syft) do not contribute findings and are not synced here.
func (e *Engine) SyncResults(projectName string, branch string, protectedBranches []string) {
	ddService := defectdojo.GetDefectDojoService(
		client.NewClient(&http.Client{}),
		environment_variable.EnvironmentVariable["DD_URL"],
		environment_variable.EnvironmentVariable["DD_ACCESS_TOKEN"])

	engagementId, err := featuresync.GetEngagementId(ddService, projectName, branch, protectedBranches)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorRetrieveEngagementId, projectName, branch))
		return
	}

	for k, scanner := range e.scanners {
		logger.Info(fmt.Sprintf(logInfoSyncResult, k))
		if err := scanner.Service.Sync(engagementId, branch, ddService); err != nil {
			logger.Error(fmt.Sprintf(logErrorSyncResult, k))
		} else {
			logger.Info(fmt.Sprintf(logInfoSyncResultSuccess, k))
		}
	}
}

// MarkFindingsByDD fetches all findings from DefectDojo for the given project and
// branch after the scan has been synced and marks each local finding with the
// appropriate status based on DefectDojo's "active" and "duplicate" fields:
//   - DUPLICATE — DD deduplication identified a prior occurrence in the product
//   - INACTIVE  — DD has suppressed the finding (false positive, accepted risk, …)
//   - ACTIVE    — finding is open and confirmed; newly found findings also get ACTIVE
//
// All local findings are returned — nothing is filtered out. If the DD fetch fails
// (e.g. the engagement does not exist yet) the original findings slice is returned
// unchanged together with the error so the caller can keep all findings at ACTIVE.
func (e *Engine) MarkFindingsByDD(findings []models.Finding, projectName string, branch string, protectedBranches []string) ([]models.Finding, error) {
	ddService := defectdojo.GetDefectDojoService(
		client.NewClient(&http.Client{}),
		environment_variable.EnvironmentVariable["DD_URL"],
		environment_variable.EnvironmentVariable["DD_ACCESS_TOKEN"])

	ddFindings, err := featuresync.GetEngagementFindings(ddService, projectName, branch, protectedBranches)
	if err != nil {
		return findings, err
	}

	return featuresync.MarkFindingsByDDFindings(findings, ddFindings), nil
}

// registerPrerequisite adds a scanner that must run and finish before any dependent
// scanner in the regular registry is allowed to start. It returns false (and logs
// an error) if the name is empty or already registered.
func (e *Engine) registerPrerequisite(name string, service interfaces.ScanServiceImpl) bool {
	if _, ok := e.prerequisites[name]; ok || name == "" {
		logger.Error(fmt.Sprintf(logErrorRegisterScanner, name))
		return false
	}

	e.prerequisites[name] = Scanner{
		Service: service,
	}

	return true
}

// registerDependentScanner adds a scanner under name in the engine's registry with
// an explicit dependency on a named prerequisite. The scanner will be skipped during
// Start if that prerequisite failed. It returns false (and logs an error) if the name
// is empty or already registered.
func (e *Engine) registerDependentScanner(name string, service interfaces.ScanServiceImpl, dependsOn string) bool {
	if _, ok := e.scanners[name]; ok || name == "" {
		logger.Error(fmt.Sprintf(logErrorRegisterScanner, name))
		return false
	}

	e.scanners[name] = Scanner{
		Service:   service,
		DependsOn: dependsOn,
	}

	return true
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

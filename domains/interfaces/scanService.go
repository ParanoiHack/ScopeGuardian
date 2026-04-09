package interfaces

import (
	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/models"
)

// ScanServiceImpl is the interface that every scanner must implement.
// It covers launching a scan, loading its findings, and syncing them to DefectDojo.
type ScanServiceImpl interface {
	Start() (bool, error)
	LoadFindings() ([]models.Finding, error)
	Sync(int, string, defectdojo.DefectDojoService) error
}

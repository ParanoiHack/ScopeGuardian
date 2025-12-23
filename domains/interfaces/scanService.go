package interfaces

import "scope-guardian/domains/models"

type ScanServiceImpl interface {
	Start() (bool, error)
	LoadFindings() ([]models.Finding, error)
	Sync() error
}

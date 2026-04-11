package kics

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
)

// GetKicsService constructs and returns a ScanServiceImpl for the KICS scanner
// using the provided loader configuration.
func GetKicsService(config loader.Config) interfaces.ScanServiceImpl {
	return newKicsService(config.Path, *config.Kics)
}
